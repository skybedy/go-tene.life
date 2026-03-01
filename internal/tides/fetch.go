package tides

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ExtremeEvent struct {
	Type      string
	TimeLocal time.Time
	HeightM   float64
}

type ExtremeFetcher func(ctx context.Context, date time.Time, loc LocationConfig) ([]ExtremeEvent, string, error)

func FetchPuertosExtremes(ctx context.Context, date time.Time, loc LocationConfig) ([]ExtremeEvent, string, error) {
	stationID := loc.PuertosStationID
	if raw := strings.TrimSpace(os.Getenv("PUERTOS_TIDE_STATION_ID")); raw != "" {
		v, err := strconv.Atoi(raw)
		if err == nil && v > 0 {
			stationID = v
		}
	}
	if stationID <= 0 {
		return nil, "", fmt.Errorf("invalid Puertos station id")
	}

	param := strings.TrimSpace(os.Getenv("PUERTOS_TIDE_PARAM"))
	if param == "" {
		param = loc.PuertosParam
	}
	if param == "" {
		return nil, "", fmt.Errorf("missing Puertos tide param")
	}

	locTZ, err := time.LoadLocation(loc.Timezone)
	if err != nil {
		return nil, "", err
	}
	startLocal := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, locTZ)
	endLocal := startLocal.Add(24*time.Hour - time.Minute)

	u := url.URL{
		Scheme: "https",
		Host:   "poem.puertos.es",
		Path:   "/portus/StationData",
	}
	q := u.Query()
	q.Set("code", strconv.Itoa(stationID))
	q.Set("params", param)
	q.Set("from", formatPortusDate(startLocal.UTC()))
	q.Set("to", formatPortusDate(endLocal.UTC()))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, "", err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("puertos status %d", resp.StatusCode)
	}

	var payload []any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, "", err
	}

	raw, _ := json.Marshal(payload)
	series, err := parsePuertosSeries(payload, locTZ)
	if err != nil {
		return nil, string(raw), err
	}

	events := detectExtremes(series)
	if len(events) == 0 {
		return nil, string(raw), fmt.Errorf("no Puertos extremes found")
	}
	return events, string(raw), nil
}

func FetchOpenMeteoExtremes(ctx context.Context, date time.Time, loc LocationConfig) ([]ExtremeEvent, string, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "marine-api.open-meteo.com",
		Path:   "/v1/marine",
	}
	q := u.Query()
	q.Set("latitude", strconv.FormatFloat(loc.Lat, 'f', 6, 64))
	q.Set("longitude", strconv.FormatFloat(loc.Lon, 'f', 6, 64))
	q.Set("minutely_15", "sea_level_height_msl")
	q.Set("timezone", loc.Timezone)
	day := date.Format("2006-01-02")
	q.Set("start_date", day)
	q.Set("end_date", day)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, "", err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("open-meteo status %d", resp.StatusCode)
	}

	var payload struct {
		Minutely15 struct {
			Time   []string  `json:"time"`
			Height []float64 `json:"sea_level_height_msl"`
		} `json:"minutely_15"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, "", err
	}

	raw, _ := json.Marshal(payload)
	if len(payload.Minutely15.Time) == 0 || len(payload.Minutely15.Time) != len(payload.Minutely15.Height) {
		return nil, string(raw), fmt.Errorf("open-meteo payload missing minutely_15 sea level data")
	}

	locTZ, err := time.LoadLocation(loc.Timezone)
	if err != nil {
		return nil, string(raw), err
	}

	series := make([]seriesPoint, 0, len(payload.Minutely15.Time))
	for i := range payload.Minutely15.Time {
		ts, err := time.ParseInLocation("2006-01-02T15:04", payload.Minutely15.Time[i], locTZ)
		if err != nil {
			return nil, string(raw), fmt.Errorf("open-meteo time parse: %w", err)
		}
		series = append(series, seriesPoint{Time: ts, Height: payload.Minutely15.Height[i]})
	}

	events := detectExtremes(series)
	if len(events) == 0 {
		return nil, string(raw), fmt.Errorf("no open-meteo extremes found")
	}
	return events, string(raw), nil
}

type seriesPoint struct {
	Time   time.Time
	Height float64
}

func parsePuertosSeries(payload []any, loc *time.Location) ([]seriesPoint, error) {
	if len(payload) < 2 {
		return nil, fmt.Errorf("puertos payload missing sections")
	}
	rowsRaw, ok := payload[1].([]any)
	if !ok || len(rowsRaw) == 0 {
		return nil, fmt.Errorf("puertos payload missing rows")
	}

	series := make([]seriesPoint, 0, len(rowsRaw))
	for _, item := range rowsRaw {
		row, ok := item.([]any)
		if !ok || len(row) < 2 {
			continue
		}
		unixTs, ok := parseUnixSeconds(row[0])
		if !ok {
			continue
		}
		pair, ok := row[1].([]any)
		if !ok || len(pair) == 0 {
			continue
		}
		h, ok := pair[0].(float64)
		if !ok {
			continue
		}
		series = append(series, seriesPoint{
			Time:   time.Unix(unixTs, 0).In(loc),
			Height: h,
		})
	}

	if len(series) < 3 {
		return nil, fmt.Errorf("puertos series too short")
	}

	sort.Slice(series, func(i, j int) bool {
		return series[i].Time.Before(series[j].Time)
	})
	return series, nil
}

func parseUnixSeconds(v any) (int64, bool) {
	f, ok := v.(float64)
	if !ok {
		return 0, false
	}
	return int64(f), true
}

func detectExtremes(series []seriesPoint) []ExtremeEvent {
	if len(series) < 3 {
		return nil
	}

	events := make([]ExtremeEvent, 0, 4)
	for i := 1; i < len(series)-1; i++ {
		prev := series[i-1]
		curr := series[i]
		next := series[i+1]
		if math.Abs(curr.Height-prev.Height) < 0.001 && math.Abs(curr.Height-next.Height) < 0.001 {
			continue
		}

		if curr.Height >= prev.Height && curr.Height > next.Height {
			events = append(events, ExtremeEvent{Type: "HIGH", TimeLocal: curr.Time, HeightM: curr.Height})
			continue
		}
		if curr.Height <= prev.Height && curr.Height < next.Height {
			events = append(events, ExtremeEvent{Type: "LOW", TimeLocal: curr.Time, HeightM: curr.Height})
		}
	}

	if len(events) == 0 {
		return nil
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].TimeLocal.Before(events[j].TimeLocal)
	})

	out := make([]ExtremeEvent, 0, len(events))
	for _, e := range events {
		if len(out) == 0 {
			out = append(out, e)
			continue
		}
		last := out[len(out)-1]
		if e.Type == last.Type && e.TimeLocal.Sub(last.TimeLocal) <= 30*time.Minute {
			if e.Type == "HIGH" && e.HeightM > last.HeightM {
				out[len(out)-1] = e
			}
			if e.Type == "LOW" && e.HeightM < last.HeightM {
				out[len(out)-1] = e
			}
			continue
		}
		out = append(out, e)
	}

	return out
}

func formatPortusDate(ts time.Time) string {
	return ts.UTC().Format("20060102@1504")
}
