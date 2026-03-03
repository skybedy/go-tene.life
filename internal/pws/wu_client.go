package pws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var errWURateLimited = errors.New("wu rate limit exceeded")

type wuFetchMeta struct {
	CacheHit bool
	Expired  bool
}

type wuCacheEntry struct {
	body    []byte
	stored  time.Time
	status  int
	expired bool
}

type wuClientConfig struct {
	cacheTTL         time.Duration
	staleFallbackMax time.Duration
	rateLimitPerMin  int
	rateLimitBurst   int
	usageLogInterval time.Duration
}

type wuClient struct {
	httpClient  *http.Client
	apiKey      string
	baseURL     string
	cacheTTL    time.Duration
	fallbackTTL time.Duration
	limiter     *rate.Limiter
	usage       *wuUsageTracker

	mu    sync.RWMutex
	cache map[string]wuCacheEntry
}

type wuRequestLog struct {
	Service    string `json:"service"`
	Path       string `json:"path"`
	StationID  string `json:"stationId,omitempty"`
	Status     int    `json:"status"`
	DurationMS int64  `json:"duration_ms"`
	CacheHit   bool   `json:"cache_hit"`
	Error      string `json:"error,omitempty"`
}

type wuUsageReport struct {
	RequestsLast1m      int64            `json:"requests_last_1m"`
	RequestsLast60m     int64            `json:"requests_last_60m"`
	RequestsToday       int64            `json:"requests_today"`
	CacheHitsLast60m    int64            `json:"cache_hits_last_60m"`
	CacheHitsToday      int64            `json:"cache_hits_today"`
	ExpiredResponses    int64            `json:"expired_responses_today"`
	ErrorRateLast60m    float64          `json:"error_rate_last_60m"`
	AvgLatencyMsLast60m float64          `json:"avg_latency_ms_last_60m"`
	TopStationIDs       []wuStationCount `json:"top_station_ids"`
	PerMinuteTotal      map[string]int64 `json:"per_minute_total"`
	PerMinuteByStation  map[string]int64 `json:"per_minute_by_station"`
	PerDayTotal         int64            `json:"per_day_total"`
	PerDayByStation     map[string]int64 `json:"per_day_by_station"`
}

type wuStationCount struct {
	StationID string `json:"station_id"`
	Count     int64  `json:"count"`
}

type wuUsageEvent struct {
	At         time.Time
	StationID  string
	Status     int
	DurationMS int64
	Error      bool
	Expired    bool
	CacheHit   bool
}

type wuUsageTracker struct {
	mu     sync.RWMutex
	events []wuUsageEvent
}

var (
	wuClientOnce sync.Once
	wuClientInst *wuClient
)

func getWUClient() *wuClient {
	wuClientOnce.Do(func() {
		cfg := wuClientConfig{
			cacheTTL:         time.Duration(envInt("WU_CACHE_TTL_SECONDS", 60)) * time.Second,
			staleFallbackMax: time.Duration(envInt("WU_STALE_FALLBACK_MAX_AGE_SECONDS", 120)) * time.Second,
			rateLimitPerMin:  envInt("WU_RATELIMIT_PER_MIN", 25),
			rateLimitBurst:   envInt("WU_RATELIMIT_BURST", 5),
			usageLogInterval: time.Minute,
		}
		wuClientInst = newWUClient(cfg, &http.Client{Timeout: defaultTimeout})
		if cfg.usageLogInterval > 0 {
			wuClientInst.usage.startSummaryLogger(cfg.usageLogInterval)
		}
	})
	return wuClientInst
}

func GetWUUsageReport() wuUsageReport {
	return getWUClient().usage.snapshot(time.Now().UTC())
}

func newWUClient(cfg wuClientConfig, httpClient *http.Client) *wuClient {
	perSecond := float64(cfg.rateLimitPerMin) / 60.0
	if perSecond <= 0 {
		perSecond = 1
	}
	burst := cfg.rateLimitBurst
	if burst <= 0 {
		burst = 1
	}
	if cfg.cacheTTL <= 0 {
		cfg.cacheTTL = 60 * time.Second
	}
	if cfg.staleFallbackMax <= 0 {
		cfg.staleFallbackMax = 120 * time.Second
	}
	return &wuClient{
		httpClient:  httpClient,
		apiKey:      strings.TrimSpace(os.Getenv("WEATHER_COM_API_KEY")),
		baseURL:     defaultBaseURL,
		cacheTTL:    cfg.cacheTTL,
		fallbackTTL: cfg.staleFallbackMax,
		limiter:     rate.NewLimiter(rate.Limit(perSecond), burst),
		usage:       &wuUsageTracker{},
		cache:       make(map[string]wuCacheEntry),
	}
}

func (c *wuClient) fetchCurrent(ctx context.Context, stationID, units, format string) ([]byte, wuFetchMeta, error) {
	meta := wuFetchMeta{}
	cacheKey := fmt.Sprintf("%s|%s|%s", stationID, units, format)
	now := time.Now().UTC()

	if body, hit, expired := c.getFreshCached(cacheKey, now); hit {
		meta.CacheHit = true
		meta.Expired = expired
		c.usage.record(wuUsageEvent{At: now, StationID: stationID, Status: 200, CacheHit: true, Expired: expired})
		return body, meta, nil
	}

	if !c.limiter.Allow() {
		if body, found, expired, age := c.getAnyCached(cacheKey, now); found && age <= c.fallbackTTL {
			meta.CacheHit = true
			meta.Expired = expired
			c.usage.record(wuUsageEvent{At: now, StationID: stationID, Status: 200, CacheHit: true, Expired: expired})
			return body, meta, nil
		}
		if err := c.limiter.Wait(ctx); err != nil {
			logWURequest(wuRequestLog{
				Service:    "wu",
				Path:       "/v2/pws/observations/current",
				StationID:  stationID,
				Status:     429,
				DurationMS: 0,
				CacheHit:   false,
				Error:      errWURateLimited.Error(),
			})
			return nil, meta, errWURateLimited
		}
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, meta, err
	}
	q := u.Query()
	q.Set("stationId", stationID)
	q.Set("format", format)
	q.Set("units", units)
	q.Set("apiKey", c.apiKey)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, meta, err
	}

	start := time.Now()
	resp, err := c.httpClient.Do(req)
	duration := time.Since(start)
	if err != nil {
		logWURequest(wuRequestLog{
			Service:    "wu",
			Path:       "/v2/pws/observations/current",
			StationID:  stationID,
			Status:     0,
			DurationMS: duration.Milliseconds(),
			CacheHit:   false,
			Error:      err.Error(),
		})
		c.usage.record(wuUsageEvent{
			At:         now,
			StationID:  stationID,
			Status:     0,
			DurationMS: duration.Milliseconds(),
			Error:      true,
		})
		return nil, meta, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	status := resp.StatusCode
	expired := isExpiredBody(body)

	logWURequest(wuRequestLog{
		Service:    "wu",
		Path:       "/v2/pws/observations/current",
		StationID:  stationID,
		Status:     status,
		DurationMS: duration.Milliseconds(),
		CacheHit:   false,
		Error:      statusError(status, body),
	})
	c.usage.record(wuUsageEvent{
		At:         now,
		StationID:  stationID,
		Status:     status,
		DurationMS: duration.Milliseconds(),
		Error:      status >= http.StatusBadRequest,
		Expired:    expired,
	})

	if status != http.StatusOK {
		// Data expired should not be treated as fatal.
		if expired {
			meta.Expired = true
			c.storeCache(cacheKey, body, status, true, now)
			return body, meta, nil
		}
		return nil, meta, fmt.Errorf("status %d: %s", status, strings.TrimSpace(string(body)))
	}

	if isExpiredObservation(body) {
		expired = true
		meta.Expired = true
	}
	c.storeCache(cacheKey, body, status, expired, now)
	meta.Expired = expired
	return body, meta, nil
}

func (c *wuClient) getFreshCached(cacheKey string, now time.Time) ([]byte, bool, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.cache[cacheKey]
	if !ok {
		return nil, false, false
	}
	if now.Sub(entry.stored) > c.cacheTTL {
		return nil, false, false
	}
	body := append([]byte(nil), entry.body...)
	return body, true, entry.expired
}

func (c *wuClient) getAnyCached(cacheKey string, now time.Time) ([]byte, bool, bool, time.Duration) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.cache[cacheKey]
	if !ok {
		return nil, false, false, 0
	}
	body := append([]byte(nil), entry.body...)
	return body, true, entry.expired, now.Sub(entry.stored)
}

func (c *wuClient) storeCache(cacheKey string, body []byte, status int, expired bool, now time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[cacheKey] = wuCacheEntry{
		body:    append([]byte(nil), body...),
		stored:  now,
		status:  status,
		expired: expired,
	}
}

func (t *wuUsageTracker) record(ev wuUsageEvent) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.events = append(t.events, ev)
	cutoff := ev.At.Add(-24 * time.Hour)
	idx := 0
	for idx < len(t.events) && t.events[idx].At.Before(cutoff) {
		idx++
	}
	if idx > 0 {
		t.events = append([]wuUsageEvent(nil), t.events[idx:]...)
	}
}

func (t *wuUsageTracker) snapshot(now time.Time) wuUsageReport {
	t.mu.RLock()
	defer t.mu.RUnlock()

	last1mFrom := now.Add(-1 * time.Minute)
	last60From := now.Add(-60 * time.Minute)
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	report := wuUsageReport{
		PerMinuteTotal:     make(map[string]int64),
		PerMinuteByStation: make(map[string]int64),
		PerDayByStation:    make(map[string]int64),
	}

	var errLast60 int64
	var durLast60 int64
	var reqLast60 int64
	for _, ev := range t.events {
		if ev.CacheHit {
			if ev.At.After(last60From) {
				report.CacheHitsLast60m++
			}
			if !ev.At.Before(dayStart) {
				report.CacheHitsToday++
			}
			continue
		}

		if ev.At.After(last1mFrom) {
			report.RequestsLast1m++
		}
		if ev.At.After(last60From) {
			report.RequestsLast60m++
			reqLast60++
			durLast60 += ev.DurationMS
			if ev.Error {
				errLast60++
			}
		}
		if !ev.At.Before(dayStart) {
			report.RequestsToday++
			report.PerDayByStation[ev.StationID]++
			if ev.Expired {
				report.ExpiredResponses++
			}
		}

		minuteKey := ev.At.Truncate(time.Minute).Format(time.RFC3339)
		report.PerMinuteTotal[minuteKey]++
		report.PerMinuteByStation[fmt.Sprintf("%s|%s", minuteKey, ev.StationID)]++
	}

	report.PerDayTotal = report.RequestsToday
	if reqLast60 > 0 {
		report.ErrorRateLast60m = float64(errLast60) / float64(reqLast60)
		report.AvgLatencyMsLast60m = float64(durLast60) / float64(reqLast60)
	}
	report.TopStationIDs = topStations(report.PerDayByStation, 10)
	return report
}

func (t *wuUsageTracker) startSummaryLogger(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			s := t.snapshot(time.Now().UTC())
			log.Printf("wu usage summary: requests_last_1m=%d requests_last_60m=%d requests_today=%d cache_hits_last_60m=%d error_rate_last_60m=%.3f avg_latency_ms_last_60m=%.1f expired_today=%d",
				s.RequestsLast1m,
				s.RequestsLast60m,
				s.RequestsToday,
				s.CacheHitsLast60m,
				s.ErrorRateLast60m,
				s.AvgLatencyMsLast60m,
				s.ExpiredResponses,
			)
		}
	}()
}

func topStations(m map[string]int64, limit int) []wuStationCount {
	out := make([]wuStationCount, 0, len(m))
	for k, v := range m {
		if k == "" {
			continue
		}
		out = append(out, wuStationCount{StationID: k, Count: v})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count == out[j].Count {
			return out[i].StationID < out[j].StationID
		}
		return out[i].Count > out[j].Count
	})
	if len(out) > limit {
		out = out[:limit]
	}
	return out
}

func statusError(status int, body []byte) string {
	if status < http.StatusBadRequest {
		return ""
	}
	b := strings.TrimSpace(string(body))
	if len(b) > 180 {
		b = b[:180]
	}
	return b
}

func logWURequest(entry wuRequestLog) {
	line, err := json.Marshal(entry)
	if err != nil {
		log.Printf("wu request: service=%s path=%s stationId=%s status=%d duration_ms=%d cache_hit=%t error=%s",
			entry.Service, entry.Path, entry.StationID, entry.Status, entry.DurationMS, entry.CacheHit, entry.Error)
		return
	}
	log.Printf("%s", line)
}

func isExpiredBody(body []byte) bool {
	b := strings.ToLower(string(body))
	return strings.Contains(b, "data expired") || strings.Contains(b, "\"dataexpired\"")
}

func isExpiredObservation(body []byte) bool {
	var payload weatherPWSResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return isExpiredBody(body)
	}
	if len(payload.Observations) == 0 {
		return isExpiredBody(body)
	}
	obs := payload.Observations[0]
	obsTime := parseObsTime(obs.ObsTimeUTC, obs.Epoch)
	if obsTime == nil {
		return isExpiredBody(body)
	}
	return time.Since(obsTime.UTC()) > 60*time.Minute
}

func envInt(name string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}
