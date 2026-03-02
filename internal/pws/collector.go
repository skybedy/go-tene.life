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
	"strconv"
	"strings"
	"time"

	"github.com/skybedy/laravel-tene.life/internal/models"
	"github.com/skybedy/laravel-tene.life/internal/store"
)

const (
	defaultBaseURL = "https://api.weather.com/v2/pws/observations/current"
	defaultTimeout = 12 * time.Second
	staleAfter     = 30 * time.Minute
	minValidTempC  = -5.0
	maxValidTempC  = 45.0
)

type weatherPWSResponse struct {
	Observations []struct {
		StationID  string   `json:"stationID"`
		ObsTimeUTC string   `json:"obsTimeUtc"`
		Epoch      int64    `json:"epoch"`
		Lat        *float64 `json:"lat"`
		Lon        *float64 `json:"lon"`
		Humidity   *float64 `json:"humidity"`
		Metric     struct {
			Temp *float64 `json:"temp"`
		} `json:"metric"`
	} `json:"observations"`
}

func CollectLatestToDB(ctx context.Context, weatherStore *store.WeatherStore) error {
	apiKey := strings.TrimSpace(os.Getenv("WEATHER_COM_API_KEY"))
	if apiKey == "" {
		return fmt.Errorf("WEATHER_COM_API_KEY is not set")
	}

	stations, err := weatherStore.GetActivePWSStations()
	if err != nil {
		return fmt.Errorf("load active stations: %w", err)
	}
	if len(stations) == 0 {
		log.Println("pws collector: no active stations in pws_stations")
		return nil
	}

	client := &http.Client{Timeout: defaultTimeout}
	fetchedAt := time.Now().UTC()

	okCount := 0
	failCount := 0
	for _, station := range stations {
		rec, fetchErr := fetchStationCurrent(ctx, client, station, apiKey, fetchedAt)
		if fetchErr != nil {
			if isWeatherAPIAccessDenied(fetchErr) {
				return fmt.Errorf("weather.com API access denied for key/host: %w", fetchErr)
			}
			failCount++
			log.Printf("pws collector: station %s (%s) failed: %v", station.StationID, station.Name, fetchErr)
			continue
		}

		if err := weatherStore.UpsertPWSLatest(rec); err != nil {
			failCount++
			log.Printf("pws collector: station %s (%s) db upsert failed: %v", station.StationID, station.Name, err)
			continue
		}

		okCount++
	}

	log.Printf("pws collector finished: total=%d ok=%d fail=%d fetched_at_utc=%s", len(stations), okCount, failCount, fetchedAt.Format(time.RFC3339))
	if okCount == 0 && failCount > 0 {
		return fmt.Errorf("pws collector failed for all stations")
	}
	return nil
}

func fetchStationCurrent(ctx context.Context, client *http.Client, station models.PWSStation, apiKey string, fetchedAt time.Time) (models.PWSLatestRecord, error) {
	rec := models.PWSLatestRecord{
		StationID:    station.StationID,
		FetchedAtUTC: fetchedAt,
	}

	u, err := url.Parse(defaultBaseURL)
	if err != nil {
		return rec, err
	}
	q := u.Query()
	q.Set("stationId", station.StationID)
	q.Set("format", "json")
	q.Set("units", "m")
	q.Set("apiKey", apiKey)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return rec, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return rec, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 768))
		return rec, fmt.Errorf("status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var payload weatherPWSResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return rec, err
	}
	if len(payload.Observations) == 0 {
		return rec, fmt.Errorf("empty observations")
	}

	obs := payload.Observations[0]
	rec.TempC = obs.Metric.Temp
	rec.Humidity = obs.Humidity
	rec.Lat = firstNonNil(obs.Lat, station.Lat)
	rec.Lon = firstNonNil(obs.Lon, station.Lon)

	obsTime := parseObsTime(obs.ObsTimeUTC, obs.Epoch)
	if obsTime != nil {
		utc := obsTime.UTC()
		rec.ObsTimeUTC = &utc
		rec.Stale = fetchedAt.Sub(utc) > staleAfter
	} else {
		rec.Stale = true
	}

	rec.Invalid = isInvalidTemperature(rec.TempC)
	if rec.Lat == nil || rec.Lon == nil {
		rec.ErrorMessage = "missing station coordinates"
	}
	if rec.Invalid {
		rec.ErrorMessage = "temperature outside expected range"
	}

	return rec, nil
}

func parseObsTime(raw string, epoch int64) *time.Time {
	raw = strings.TrimSpace(raw)
	if raw != "" {
		layouts := []string{
			time.RFC3339,
			"2006-01-02T15:04:05Z0700",
			"2006-01-02 15:04:05",
		}
		for _, layout := range layouts {
			if ts, err := time.Parse(layout, raw); err == nil {
				utc := ts.UTC()
				return &utc
			}
		}
	}

	if epoch > 0 {
		ts := time.Unix(epoch, 0).UTC()
		return &ts
	}
	return nil
}

func isInvalidTemperature(temp *float64) bool {
	if temp == nil {
		return true
	}
	return *temp < minValidTempC || *temp > maxValidTempC
}

func firstNonNil(primary, fallback *float64) *float64 {
	if primary != nil {
		v := *primary
		return &v
	}
	if fallback != nil {
		v := *fallback
		return &v
	}
	return nil
}

func CollectorInterval() time.Duration {
	interval := 10 * time.Minute
	raw := strings.TrimSpace(os.Getenv("PWS_COLLECT_INTERVAL_MINUTES"))
	if raw == "" {
		return interval
	}
	mins, err := strconv.Atoi(raw)
	if err != nil || mins <= 0 {
		return interval
	}
	return time.Duration(mins) * time.Minute
}

func APIKeyConfigured() bool {
	return strings.TrimSpace(os.Getenv("WEATHER_COM_API_KEY")) != ""
}

func isWeatherAPIAccessDenied(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "status 401") || strings.Contains(msg, "status 403") {
		return true
	}
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		nested := strings.ToLower(urlErr.Error())
		return strings.Contains(nested, "status 401") || strings.Contains(nested, "status 403")
	}
	return strings.Contains(msg, "access denied")
}
