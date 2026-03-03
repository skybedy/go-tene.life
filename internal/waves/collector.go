package waves

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/skybedy/laravel-tene.life/internal/models"
)

const (
	defaultLocation  = "Los Cristianos"
	defaultStationID = 2446
	defaultSource    = "Puertos del Estado (PORTUS)"
	defaultOutPath   = "data/waves_latest.json"
)

func CollectLatestToDefaultPath(ctx context.Context) error {
	return CollectLatestMeasured(ctx, defaultStationID, defaultLocation, outputPath())
}

func CollectLatestMeasured(ctx context.Context, stationID int, location, outPath string) error {
	if stationID <= 0 {
		return fmt.Errorf("invalid station id: %d", stationID)
	}

	now := time.Now().UTC()
	from := now.Add(-48 * time.Hour)

	payload, err := fetchStationData(ctx, stationID, from, now)
	if err != nil {
		return err
	}

	waves, err := extractLatestMeasured(payload, stationID, location, now)
	if err != nil {
		return err
	}

	return writeAtomicJSON(outPath, waves)
}

func LoadLatestFromFile(path string) (*models.WavesLatest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out models.WavesLatest
	if err := json.NewDecoder(f).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func fetchStationData(ctx context.Context, stationID int, from, to time.Time) ([]any, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "poem.puertos.es",
		Path:   "/portus/StationData",
	}

	q := u.Query()
	q.Set("code", strconv.Itoa(stationID))
	q.Set("params", "Hm0,Tp,Tm02,MeanDir180")
	q.Set("from", formatStationDataTime(from))
	q.Set("to", formatStationDataTime(to))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 12 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("stationdata status %d", resp.StatusCode)
	}

	var payload []any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if len(payload) < 2 {
		return nil, fmt.Errorf("stationdata payload missing sections")
	}

	return payload, nil
}

func formatStationDataTime(ts time.Time) string {
	utc := ts.UTC()
	return utc.Format("20060102@1504")
}

func extractLatestMeasured(payload []any, stationID int, location string, fetchedAt time.Time) (*models.WavesLatest, error) {
	headersRaw, ok := payload[0].([]any)
	if !ok {
		return nil, fmt.Errorf("invalid headers format")
	}
	rowsRaw, ok := payload[1].([]any)
	if !ok {
		return nil, fmt.Errorf("invalid rows format")
	}
	if len(rowsRaw) == 0 {
		return nil, fmt.Errorf("stationdata has no measured rows")
	}

	headers := make([]string, 0, len(headersRaw))
	for _, h := range headersRaw {
		headers = append(headers, fmt.Sprint(h))
	}

	idxHm0 := findHeaderIndex(headers, "hm0")
	idxTp := findHeaderIndex(headers, "tp")
	idxTm02 := findHeaderIndex(headers, "tm02")
	idxMeanDir := findHeaderIndex(headers, "meandir")
	if idxHm0 == -1 {
		return nil, fmt.Errorf("missing Hm0 column")
	}

	for i := len(rowsRaw) - 1; i >= 0; i-- {
		row, ok := rowsRaw[i].([]any)
		if !ok || len(row) == 0 {
			continue
		}

		ts, ok := parseUnixSeconds(row[0])
		if !ok {
			continue
		}

		hm0, ok := extractMeasureValue(row, idxHm0)
		if !ok {
			continue
		}

		period, periodOK := extractMeasureValue(row, idxTp)
		if !periodOK {
			period, periodOK = extractMeasureValue(row, idxTm02)
		}
		if !periodOK {
			continue
		}

		meanDir, meanOK := extractMeasureValue(row, idxMeanDir)
		if !meanOK {
			continue
		}

		return &models.WavesLatest{
			Location:      location,
			StationID:     stationID,
			DataKind:      "measured",
			Source:        defaultSource,
			MeasuredAtUTC: time.Unix(ts, 0).UTC().Format(time.RFC3339),
			Hm0M:          hm0,
			PeriodS:       period,
			MeanDirDeg:    meanDir,
			FetchedAtUTC:  fetchedAt.UTC().Format(time.RFC3339),
		}, nil
	}

	return nil, fmt.Errorf("no valid measured row with Hm0/period/MeanDir found")
}

func parseUnixSeconds(v any) (int64, bool) {
	f, ok := v.(float64)
	if !ok {
		return 0, false
	}
	return int64(f), true
}

func extractMeasureValue(row []any, idx int) (float64, bool) {
	if idx < 0 || idx >= len(row) {
		return 0, false
	}
	pair, ok := row[idx].([]any)
	if !ok || len(pair) == 0 {
		return 0, false
	}
	v, ok := pair[0].(float64)
	if !ok {
		return 0, false
	}
	return v, true
}

func findHeaderIndex(headers []string, needle string) int {
	needle = strings.ToLower(strings.TrimSpace(needle))
	for i, h := range headers {
		hh := strings.ToLower(strings.TrimSpace(h))
		if hh == needle || strings.HasPrefix(hh, needle) || strings.HasPrefix(hh, needle+" ") || strings.HasPrefix(hh, needle+"(") {
			return i
		}
	}
	return -1
}

func writeAtomicJSON(path string, payload *models.WavesLatest) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	return nil
}

func outputPath() string {
	if raw := strings.TrimSpace(os.Getenv("WAVES_JSON_PATH")); raw != "" {
		return raw
	}
	return defaultOutPath
}
