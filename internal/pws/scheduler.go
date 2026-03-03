package pws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	neturl "net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/skybedy/laravel-tene.life/internal/models"
	"github.com/skybedy/laravel-tene.life/internal/store"
)

var errCircuitBreakerTripped = errors.New("pws circuit breaker tripped")

type pacedSchedulerConfig struct {
	window         time.Duration
	retryMax       int
	retryMin       time.Duration
	retryMaxDelay  time.Duration
	requestTimeout time.Duration
	cursorPath     string
}

type cursorState struct {
	Next         int    `json:"next"`
	UpdatedAtUTC string `json:"updated_at_utc"`
}

func RunPacedCollector(ctx context.Context, weatherStore *store.WeatherStore) error {
	if weatherStore == nil {
		return fmt.Errorf("weatherStore is nil")
	}
	if !APIKeyConfigured() {
		return fmt.Errorf("WEATHER_COM_API_KEY is not set")
	}

	cfg := loadPacedSchedulerConfig()
	wu := getWUClient()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	nextIdx := loadCursorIndex(cfg.cursorPath)

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		stations, err := weatherStore.GetActivePWSStations()
		if err != nil {
			log.Printf("pws paced collector: failed loading stations: %v", err)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(1 * time.Minute):
				continue
			}
		}
		if len(stations) == 0 {
			log.Printf("pws paced collector: no active stations, waiting")
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(1 * time.Minute):
				continue
			}
		}

		if nextIdx < 0 || nextIdx >= len(stations) {
			nextIdx = 0
		}

		pace := cfg.paceInterval(len(stations))
		log.Printf("pws paced collector: active_stations=%d pace=%s retry_max=%d", len(stations), pace, cfg.retryMax)
		ticker := time.NewTicker(pace)

		for {
			station := stations[nextIdx]
			if err := collectStationWithRetry(ctx, cfg, wu, weatherStore, station, rng); err != nil {
				ticker.Stop()
				if errors.Is(err, errCircuitBreakerTripped) {
					return err
				}
				log.Printf("pws paced collector: station %s failed permanently: %v", station.StationID, err)
			}

			nextIdx = (nextIdx + 1) % len(stations)
			if err := saveCursorIndex(cfg.cursorPath, nextIdx); err != nil {
				log.Printf("pws paced collector: save cursor failed: %v", err)
			}

			if nextIdx == 0 {
				ticker.Stop()
				break
			}

			select {
			case <-ctx.Done():
				ticker.Stop()
				return ctx.Err()
			case <-ticker.C:
			}
		}
	}
}

func collectStationWithRetry(
	ctx context.Context,
	cfg pacedSchedulerConfig,
	wu *wuClient,
	weatherStore *store.WeatherStore,
	station models.PWSStation,
	rng *rand.Rand,
) error {
	fetchedAt := time.Now().UTC()
	lastErr := error(nil)

	for attempt := 0; attempt <= cfg.retryMax; attempt++ {
		reqCtx, cancel := context.WithTimeout(ctx, cfg.requestTimeout)
		rec, err := fetchStationCurrent(reqCtx, wu, station, fetchedAt)
		cancel()

		if err == nil {
			if upsertErr := weatherStore.UpsertPWSLatest(rec); upsertErr != nil {
				log.Printf("pws paced collector: station %s db upsert failed: %v", station.StationID, upsertErr)
			}
			return nil
		}

		if isCircuitBreakerError(err) {
			return fmt.Errorf("%w: station=%s err=%v", errCircuitBreakerTripped, station.StationID, err)
		}

		lastErr = err
		if !isRetryableError(err) || attempt == cfg.retryMax {
			log.Printf("pws paced collector: station %s failed after retries, keeping previous DB value: %v", station.StationID, err)
			return nil
		}

		backoff := retryBackoffWithJitter(cfg.retryMin, cfg.retryMaxDelay, rng)
		log.Printf("pws paced collector: station %s retry %d/%d in %s due to: %v", station.StationID, attempt+1, cfg.retryMax, backoff, err)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
	}

	if lastErr != nil {
		return lastErr
	}
	return nil
}

func (c pacedSchedulerConfig) paceInterval(stations int) time.Duration {
	if stations <= 0 {
		return c.window
	}
	interval := c.window / time.Duration(stations)
	if interval < time.Second {
		return time.Second
	}
	return interval
}

func loadPacedSchedulerConfig() pacedSchedulerConfig {
	windowMins := envInt("PWS_PACED_WINDOW_MINUTES", 0)
	if windowMins <= 0 {
		// Backward-compatible fallback to legacy setting.
		windowMins = envInt("PWS_COLLECT_INTERVAL_MINUTES", 10)
	}
	if windowMins <= 0 {
		windowMins = 10
	}

	retryMinSec := envInt("PWS_RETRY_MIN_SECONDS", 20)
	retryMaxSec := envInt("PWS_RETRY_MAX_SECONDS", 45)
	if retryMinSec <= 0 {
		retryMinSec = 20
	}
	if retryMaxSec < retryMinSec {
		retryMaxSec = retryMinSec
	}

	cursorPath := strings.TrimSpace(os.Getenv("PWS_CURSOR_PATH"))
	if cursorPath == "" {
		cursorPath = "data/pws_cursor.json"
	}

	return pacedSchedulerConfig{
		window:         time.Duration(windowMins) * time.Minute,
		retryMax:       envInt("PWS_RETRY_MAX", 2),
		retryMin:       time.Duration(retryMinSec) * time.Second,
		retryMaxDelay:  time.Duration(retryMaxSec) * time.Second,
		requestTimeout: 20 * time.Second,
		cursorPath:     cursorPath,
	}
}

func retryBackoffWithJitter(minDelay, maxDelay time.Duration, rng *rand.Rand) time.Duration {
	if maxDelay <= minDelay {
		return minDelay
	}
	delta := maxDelay - minDelay
	return minDelay + time.Duration(rng.Int63n(int64(delta)+1))
}

func isCircuitBreakerError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, errWURateLimited) {
		return true
	}
	if code, ok := statusCodeFromError(err); ok {
		return code == 401 || code == 403 || code == 429
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "status 401") || strings.Contains(msg, "status 403") || strings.Contains(msg, "status 429")
}

func isRetryableError(err error) bool {
	if err == nil || isCircuitBreakerError(err) {
		return false
	}

	if code, ok := statusCodeFromError(err); ok {
		return code >= 500 && code <= 599
	}

	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	var urlErr *neturl.Error
	if errors.As(err, &urlErr) {
		return true
	}

	return false
}

func loadCursorIndex(path string) int {
	b, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	var cur cursorState
	if err := json.Unmarshal(b, &cur); err != nil {
		return 0
	}
	if cur.Next < 0 {
		return 0
	}
	return cur.Next
}

func saveCursorIndex(path string, next int) error {
	if next < 0 {
		next = 0
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	payload, err := json.Marshal(cursorState{
		Next:         next,
		UpdatedAtUTC: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}
	payload = append(payload, '\n')

	tmp := path + ".tmp." + strconv.FormatInt(time.Now().UnixNano(), 10)
	if err := os.WriteFile(tmp, payload, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}
