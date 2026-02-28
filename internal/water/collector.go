package water

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/skybedy/laravel-tene.life/internal/models"
)

const (
	wmsBaseURL           = "https://idecan2.grafcan.es/ServicioWMS/ZB_PM?"
	defaultOutputPath    = "data/water_quality_latest.json"
	defaultLocationLabel = "Los Cristianos"
)

var (
	reTag            = regexp.MustCompile(`<[^>]+>`)
	reDenominacion   = regexp.MustCompile(`(?is)Denominación:\s*</td>\s*<td[^>]*>([^<]+)`)
	reMunicipio      = regexp.MustCompile(`(?is)Municipio:\s*</td>\s*<td[^>]*>([^<]+)`)
	reCodigo         = regexp.MustCompile(`(?is)Código del punto:\s*</td>\s*<td[^>]*>([^<]+)`)
	reCalidadAnual   = regexp.MustCompile(`(?is)Calidad anual\s*([0-9]{4})?:\s*</td>\s*<td[^>]*>(.*?)</td>`)
	reStatusWord     = regexp.MustCompile(`(?i)\b(EXCELENTE|BUENA|SUFICIENTE|INSUFICIENTE)\b`)
	reWhitespace     = regexp.MustCompile(`\s+`)
	reLinkCandidates = regexp.MustCompile(`(?is)href="([^"]+)"`)
)

type pmRecord struct {
	Denominacion string
	Municipio    string
	Codigo       string
	Status       string
	Year         string
}

func CollectLatestToDefaultPath(ctx context.Context) error {
	lat := envFloat("WATER_LAT", 28.043493664088118)
	lon := envFloat("WATER_LON", -16.710318262067073)
	radius := envFloat("WATER_BBOX_RADIUS", 0.08)
	nameFilter := strings.TrimSpace(os.Getenv("WATER_NAME_FILTER"))
	if nameFilter == "" {
		nameFilter = "CRISTIANOS"
	}
	location := strings.TrimSpace(os.Getenv("WATER_LOCATION"))
	if location == "" {
		location = defaultLocationLabel
	}
	return CollectLatestOfficial(ctx, lat, lon, radius, nameFilter, location, defaultOutputPath)
}

func CollectLatestOfficial(ctx context.Context, lat, lon, radius float64, nameFilter, location, outPath string) error {
	featureInfoURL, err := buildFeatureInfoURL(lat, lon, radius)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, featureInfoURL, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("wms getfeatureinfo status %d", resp.StatusCode)
	}

	body, err := ioReadAllLimit(resp, 2<<20)
	if err != nil {
		return err
	}

	records := parsePMRecords(string(body))
	if len(records) == 0 {
		return fmt.Errorf("wms returned no PM records in GetFeatureInfo window")
	}
	best := selectBestRecord(records, nameFilter)

	notes := "Clasificación anual oficial de aguas de baño (capa WMS PM); no incluye fecha del último análisis."
	if best.Year != "" {
		notes = fmt.Sprintf("Clasificación anual %s (capa WMS PM); no incluye fecha del último análisis.", best.Year)
	}
	link := wmsBaseURL
	if href := firstMeaningfulLink(string(body)); href != "" {
		link = href
	}

	payload := &models.WaterQualityLatest{
		Location:     location,
		DataKind:     "official",
		Source:       "IDECanarias / GRAFCAN WMS (ZB_PM, layer PM)",
		Status:       normalizeStatus(best.Status),
		SampleDate:   "",
		UpdatedAtUTC: time.Now().UTC().Format(time.RFC3339),
		Notes:        notes,
		Link:         link,
	}

	return writeAtomicJSON(outPath, payload)
}

func LoadLatestFromFile(path string) (*models.WaterQualityLatest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out models.WaterQualityLatest
	if err := json.NewDecoder(f).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func buildFeatureInfoURL(lat, lon, radius float64) (string, error) {
	if radius <= 0 {
		return "", fmt.Errorf("radius must be positive")
	}
	minLat := lat - radius
	minLon := lon - radius
	maxLat := lat + radius
	maxLon := lon + radius

	u, err := url.Parse(wmsBaseURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("SERVICE", "WMS")
	q.Set("VERSION", "1.3.0")
	q.Set("REQUEST", "GetFeatureInfo")
	q.Set("LAYERS", "PM")
	q.Set("QUERY_LAYERS", "PM")
	q.Set("CRS", "EPSG:4326")
	q.Set("BBOX", fmt.Sprintf("%.6f,%.6f,%.6f,%.6f", minLat, minLon, maxLat, maxLon))
	q.Set("WIDTH", "200")
	q.Set("HEIGHT", "200")
	q.Set("I", "100")
	q.Set("J", "100")
	q.Set("INFO_FORMAT", "text/html")
	q.Set("FEATURE_COUNT", "50")
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func parsePMRecords(raw string) []pmRecord {
	dens := reDenominacion.FindAllStringSubmatch(raw, -1)
	muns := reMunicipio.FindAllStringSubmatch(raw, -1)
	cods := reCodigo.FindAllStringSubmatch(raw, -1)
	cals := reCalidadAnual.FindAllStringSubmatch(raw, -1)

	n := minInt(len(dens), len(muns), len(cods), len(cals))
	if n == 0 {
		return nil
	}

	out := make([]pmRecord, 0, n)
	for i := 0; i < n; i++ {
		status := statusFromCalidadBlock(cals[i][2])
		if status == "" {
			continue
		}
		out = append(out, pmRecord{
			Denominacion: cleanText(dens[i][1]),
			Municipio:    cleanText(muns[i][1]),
			Codigo:       cleanText(cods[i][1]),
			Status:       status,
			Year:         cleanText(cals[i][1]),
		})
	}
	return out
}

func selectBestRecord(records []pmRecord, nameFilter string) pmRecord {
	filter := strings.ToUpper(strings.TrimSpace(nameFilter))
	bestIdx := 0
	bestScore := -1
	for i, r := range records {
		score := 0
		den := strings.ToUpper(r.Denominacion)
		mun := strings.ToUpper(r.Municipio)
		if filter != "" && strings.Contains(den, filter) {
			score += 5
		}
		if strings.Contains(mun, "ARONA") {
			score += 2
		}
		if strings.Contains(den, "CRISTIANOS") {
			score += 3
		}
		if score > bestScore {
			bestScore = score
			bestIdx = i
		}
	}
	return records[bestIdx]
}

func statusFromCalidadBlock(calidad string) string {
	text := strings.ToUpper(cleanText(calidad))
	m := reStatusWord.FindStringSubmatch(text)
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

func normalizeStatus(raw string) string {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "EXCELENTE":
		return "Excelente"
	case "BUENA":
		return "Buena"
	case "SUFICIENTE":
		return "Suficiente"
	case "INSUFICIENTE":
		return "Insuficiente"
	default:
		return strings.TrimSpace(raw)
	}
}

func firstMeaningfulLink(raw string) string {
	matches := reLinkCandidates.FindAllStringSubmatch(raw, -1)
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		href := strings.TrimSpace(html.UnescapeString(m[1]))
		if href == "" {
			continue
		}
		// Ignore generic law link from quality description.
		if strings.Contains(strings.ToLower(href), "boe.es") {
			continue
		}
		return href
	}
	return ""
}

func cleanText(s string) string {
	s = reTag.ReplaceAllString(s, " ")
	s = html.UnescapeString(s)
	s = strings.TrimSpace(reWhitespace.ReplaceAllString(s, " "))
	return s
}

func writeAtomicJSON(path string, payload *models.WaterQualityLatest) error {
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

func envFloat(key string, fallback float64) float64 {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return fallback
	}
	return v
}

func minInt(values ...int) int {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

func ioReadAllLimit(resp *http.Response, max int64) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(io.LimitReader(resp.Body, max))
}
