package web

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"sync"

	"github.com/labstack/echo/v4"
	"github.com/skybedy/laravel-tene.life/internal/i18n"
	"github.com/skybedy/laravel-tene.life/internal/models"
	"github.com/skybedy/laravel-tene.life/internal/store"
	"github.com/skybedy/laravel-tene.life/internal/tides"
	"github.com/skybedy/laravel-tene.life/internal/utils"
	"github.com/skybedy/laravel-tene.life/internal/water"
	"github.com/skybedy/laravel-tene.life/internal/waves"
)

// TemplateRenderer implements Echo's Renderer interface
type TemplateRenderer struct {
	Templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

type Handler struct {
	WeatherStore  *store.WeatherStore
	TideCollector *tides.Collector
	weatherCache  *models.WeatherData
	seaTempCache  *float64
	cacheMu       sync.RWMutex
	lastCache     time.Time
	cacheTimeout  time.Duration
}

func NewHandler(ws *store.WeatherStore) *Handler {
	return &Handler{
		WeatherStore:  ws,
		TideCollector: tides.NewCollector(ws),
		cacheTimeout:  30 * time.Second, // Default cache timeout
	}
}

func (h *Handler) getLocale(c echo.Context) string {
	return i18n.NormalizeLocale(c.Param("locale"))
}

func (h *Handler) getGAConfig() (bool, string) {
	measurementID := strings.TrimSpace(os.Getenv("GA_MEASUREMENT_ID"))
	if measurementID == "" {
		return false, ""
	}
	return true, measurementID
}

func (h *Handler) getCommonViewData(c echo.Context) (string, string, []models.LanguageOption, map[string]string, bool, string) {
	locale := h.getLocale(c)
	currentPath := i18n.StripLocalePrefix(c.Request().URL.Path)
	gaEnabled, gaMeasurementID := h.getGAConfig()
	return locale, currentPath, i18n.SupportedLanguages(), i18n.Messages(locale), gaEnabled, gaMeasurementID
}

func (h *Handler) webcamImageURL() string {
	webcamPath := utils.EnvPathOrDefault("WEBCAM_IMAGE_PATH", "public/images/tenelife.jpg")
	info, err := os.Stat(webcamPath)
	if err != nil {
		return "/webcam/image.jpg"
	}
	return fmt.Sprintf("/webcam/image.jpg?v=%d", info.ModTime().Unix())
}

func hasDailyValues(day models.WeatherDaily) bool {
	return day.AvgTemperature != nil ||
		day.MinTemperature != nil ||
		day.MaxTemperature != nil ||
		day.AvgPressure != nil ||
		day.AvgHumidity != nil
}

func filterDailyStatsWithValues(stats []models.WeatherDaily) []models.WeatherDaily {
	filtered := make([]models.WeatherDaily, 0, len(stats))
	for _, day := range stats {
		if hasDailyValues(day) {
			filtered = append(filtered, day)
		}
	}
	return filtered
}

func (h *Handler) IndexHandler(c echo.Context) error {
	locale, currentPath, languages, messages, gaEnabled, gaMeasurementID := h.getCommonViewData(c)
	webcamImageURL := h.webcamImageURL()

	// 1. Get Weather Data (with improved caching)
	weather, seaTemp, err := h.getCachedWeatherData()
	if err != nil {
		log.Printf("Error getting weather data: %v", err)
		// Continue with cached data if available, or empty data
		if weather == nil && seaTemp == nil {
			loc, _ := time.LoadLocation("Atlantic/Canary")
			now := time.Now().In(loc)
			return c.Render(http.StatusOK, "index.html", models.PageData{
				FormattedDate:   now.Format("2. 1. 2006"),
				FormattedTime:   now.Format("15:04"),
				WebcamImageURL:  webcamImageURL,
				Locale:          locale,
				LocalePrefix:    i18n.LocalePrefix(locale),
				CurrentPath:     currentPath,
				CurrentSection:  "home",
				Languages:       languages,
				I18n:            messages,
				GAEnabled:       gaEnabled,
				GAMeasurementID: gaMeasurementID,
			})
		}
	}

	// 2. Format Date/Time (force Atlantic/Canary timezone)
	loc, _ := time.LoadLocation("Atlantic/Canary")
	ts := time.Now().In(loc)
	if weather != nil && weather.Timestamp > 0 {
		ts = time.Unix(weather.Timestamp, 0).In(loc)
	}

	var dayMaxTemperature *float64
	var dayMinTemperature *float64
	dayMaxTime := ""
	dayMinTime := ""
	dayMaxTempText := ""
	dayMinTempText := ""
	if maxTemp, maxAt, minTemp, minAt, extErr := h.WeatherStore.GetDailyTemperatureExtremes(ts.Format("2006-01-02")); extErr == nil {
		dayMaxTemperature = maxTemp
		dayMinTemperature = minTemp
		if maxTemp != nil {
			dayMaxTempText = fmt.Sprintf("%.1f", *maxTemp)
		}
		if minTemp != nil {
			dayMinTempText = fmt.Sprintf("%.1f", *minTemp)
		}
		if maxAt != nil {
			dayMaxTime = maxAt.In(loc).Format("15:04")
		}
		if minAt != nil {
			dayMinTime = minAt.In(loc).Format("15:04")
		}
	} else {
		log.Printf("Error fetching daily temperature extremes: %v", extErr)
	}

	var seaTempVal float64
	if seaTemp != nil {
		seaTempVal = *seaTemp
	}

	tideHighEvents, tideLowEvents := h.getCachedTideData(ts)
	nextHighTide := ""
	nextLowTide := ""
	if len(tideHighEvents) > 0 {
		nextHighTide = strings.Join(tideHighEvents, ", ")
	}
	if len(tideLowEvents) > 0 {
		nextLowTide = strings.Join(tideLowEvents, ", ")
	}
	var waveData *models.WavesLatest
	var waterData *models.WaterQualityLatest
	wavePath := utils.EnvPathOrDefault("WAVES_JSON_PATH", "data/waves_latest.json")
	waveData, err = waves.LoadLatestFromFile(wavePath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Error reading wave cache '%s': %v", wavePath, err)
		}
		waveData = nil
	}
	waterPath := utils.EnvPathOrDefault("WATER_JSON_PATH", "data/water_quality_latest.json")
	waterData, err = water.LoadLatestFromFile(waterPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Water cache file not found at '%s' (set WATER_JSON_PATH to an absolute path on production)", waterPath)
		} else {
			log.Printf("Error reading water cache '%s': %v", waterPath, err)
		}
		waterData = nil
	}

	data := models.PageData{
		Weather:           weather,
		WebcamImageURL:    webcamImageURL,
		SeaTemperature:    seaTemp,
		SeaTemperatureVal: seaTempVal,
		NextHighTide:      nextHighTide,
		NextLowTide:       nextLowTide,
		TideHighEvents:    tideHighEvents,
		TideLowEvents:     tideLowEvents,
		Waves:             waveData,
		WaterQuality:      waterData,
		DayMaxTemperature: dayMaxTemperature,
		DayMinTemperature: dayMinTemperature,
		DayMaxTempText:    dayMaxTempText,
		DayMinTempText:    dayMinTempText,
		DayMaxTime:        dayMaxTime,
		DayMinTime:        dayMinTime,
		FormattedDate:     ts.Format("2. 1. 2006"),
		FormattedTime:     ts.Format("15:04"),
		PageTitle:         "",
		Locale:            locale,
		LocalePrefix:      i18n.LocalePrefix(locale),
		CurrentPath:       currentPath,
		CurrentSection:    "home",
		Languages:         languages,
		I18n:              messages,
		GAEnabled:         gaEnabled,
		GAMeasurementID:   gaMeasurementID,
	}

	return c.Render(http.StatusOK, "index.html", data)
}

// getCachedWeatherData gets weather data with improved caching and error handling
func (h *Handler) getCachedWeatherData() (*models.WeatherData, *float64, error) {
	h.cacheMu.RLock()
	weather := h.weatherCache
	seaTemp := h.seaTempCache
	cacheAge := time.Since(h.lastCache)
	h.cacheMu.RUnlock()

	// Refresh cache if older than timeout or empty
	if weather == nil || seaTemp == nil || cacheAge > h.cacheTimeout {
		h.cacheMu.Lock()

		// Ensure unlock happens when we return
		defer h.cacheMu.Unlock()

		// Double check after acquiring lock
		if h.weatherCache == nil || h.seaTempCache == nil || time.Since(h.lastCache) > h.cacheTimeout {
			// Update Weather JSON
			weatherPath := utils.EnvPathOrDefault("WEATHER_JSON_PATH", "public/files/weather.json")

			file, err := os.Open(weatherPath)
			if err != nil {
				return h.weatherCache, h.seaTempCache, utils.NewInternalServerError(
					"Failed to open weather file", err)
			}
			defer file.Close()

			decoder := json.NewDecoder(file)
			newWeather := &models.WeatherData{}
			if err := decoder.Decode(newWeather); err != nil {
				return h.weatherCache, h.seaTempCache, utils.NewInternalServerError(
					"Failed to decode weather data", err)
			}
			h.weatherCache = newWeather
			weather = newWeather

			// Update Sea Temp from DB based on weather timestamp date.
			refDate := time.Now().Format("2006-01-02")
			if newWeather.Timestamp > 0 {
				refDate = time.Unix(newWeather.Timestamp, 0).Format("2006-01-02")
			}
			newSeaTemp, err := h.WeatherStore.GetLatestSeaTemperature(refDate)
			if err != nil {
				return weather, h.seaTempCache, utils.NewInternalServerError(
					"Failed to get sea temperature", err)
			}
			h.seaTempCache = newSeaTemp
			seaTemp = newSeaTemp

			h.lastCache = time.Now()
		} else {
			weather = h.weatherCache
			seaTemp = h.seaTempCache
		}
	}

	return weather, seaTemp, nil
}

func (h *Handler) getCachedTideData(reference time.Time) ([]string, []string) {
	_ = reference
	locCfg, ok := tides.ResolveLocation("los_cristianos")
	if !ok {
		return nil, nil
	}

	tz, err := time.LoadLocation(locCfg.Timezone)
	if err != nil {
		return nil, nil
	}
	// Tide events are expected for today's local day, independent of weather feed timestamp freshness.
	dateLocal := time.Now().In(tz).Format("2006-01-02")
	preferredSource := tides.ServingSource()

	ctx, cancel := context.WithTimeout(context.Background(), 1200*time.Millisecond)
	events, err := h.WeatherStore.GetTideEvents(ctx, dateLocal, locCfg.Key)
	cancel()
	if err != nil {
		return nil, nil
	}

	selected := selectTideEventsBySource(events, preferredSource)
	if len(selected) == 0 {
		collectCtx, collectCancel := context.WithTimeout(context.Background(), 4*time.Second)
		_ = h.TideCollector.CollectTides(collectCtx, dateLocal, locCfg.Key)
		collectCancel()

		refetchCtx, refetchCancel := context.WithTimeout(context.Background(), 1200*time.Millisecond)
		events, err = h.WeatherStore.GetTideEvents(refetchCtx, dateLocal, locCfg.Key)
		refetchCancel()
		if err != nil {
			return nil, nil
		}
		selected = selectTideEventsBySource(events, preferredSource)
	}

	if len(selected) == 0 {
		return nil, nil
	}

	high := make([]string, 0, 2)
	low := make([]string, 0, 2)
	for _, ev := range selected {
		timeLabel := strings.TrimPrefix(ev.EventTimeLocal.Format("15:04"), "0")
		item := fmt.Sprintf("%5s (%s m)", timeLabel, formatTideHeightSigned(ev.HeightM))
		switch ev.EventType {
		case "HIGH":
			high = append(high, item)
		case "LOW":
			low = append(low, item)
		}
	}

	return high, low
}

func formatTideHeightSigned(v float64) string {
	if v < 0 {
		return "−" + fmt.Sprintf("%.2f", math.Abs(v))
	}
	return "+" + fmt.Sprintf("%.2f", v)
}

func (h *Handler) GetHourlyDataHandler(c echo.Context) error {
	date := c.QueryParam("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		appErr := utils.NewBadRequestError("Date must be in YYYY-MM-DD format", err)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "bad_request",
			Message: appErr.Message,
		})
	}

	data, err := h.WeatherStore.GetHourlyData(date)
	if err != nil {
		appErr := utils.NewInternalServerError("Failed to get hourly data", err)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "internal_server_error",
			Message: appErr.Message,
		})
	}

	response := models.HourlyChartResponse{}
	for _, record := range data {
		response.Labels = append(response.Labels, fmt.Sprintf("%02d:00", record.Hour))
		response.Datasets.Temperature = append(response.Datasets.Temperature, record.AvgTemperature)
		response.Datasets.Pressure = append(response.Datasets.Pressure, record.AvgPressure)
		response.Datasets.Humidity = append(response.Datasets.Humidity, record.AvgHumidity)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetHomeDataHandler(c echo.Context) error {
	weather, seaTemp, err := h.getCachedWeatherData()
	if err != nil {
		log.Printf("Error getting weather data for /api/home: %v", err)
	}

	var waveData *models.WavesLatest
	var waterData *models.WaterQualityLatest
	wavePath := utils.EnvPathOrDefault("WAVES_JSON_PATH", "data/waves_latest.json")
	waveData, err = waves.LoadLatestFromFile(wavePath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Error reading wave cache '%s': %v", wavePath, err)
		}
		waveData = nil
	}
	waterPath := utils.EnvPathOrDefault("WATER_JSON_PATH", "data/water_quality_latest.json")
	waterData, err = water.LoadLatestFromFile(waterPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Water cache file not found at '%s' (set WATER_JSON_PATH to an absolute path on production)", waterPath)
		} else {
			log.Printf("Error reading water cache '%s': %v", waterPath, err)
		}
		waterData = nil
	}

	return c.JSON(http.StatusOK, models.HomeAPIResponse{
		Weather:        weather,
		SeaTemperature: seaTemp,
		Waves:          waveData,
		WaterQuality:   waterData,
	})
}

func (h *Handler) GetTidesHandler(c echo.Context) error {
	locKey := strings.TrimSpace(c.QueryParam("loc"))
	locCfg, ok := tides.ResolveLocation(locKey)
	if !ok {
		appErr := utils.NewBadRequestError("Unknown location key", nil)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "bad_request",
			Message: appErr.Message,
		})
	}

	tz, err := time.LoadLocation(locCfg.Timezone)
	if err != nil {
		appErr := utils.NewInternalServerError("Failed to load tide timezone", err)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "internal_server_error",
			Message: appErr.Message,
		})
	}

	dateLocal := strings.TrimSpace(c.QueryParam("date"))
	if dateLocal == "" {
		dateLocal = time.Now().In(tz).Format("2006-01-02")
	}
	if _, err := time.Parse("2006-01-02", dateLocal); err != nil {
		appErr := utils.NewBadRequestError("Date must be in YYYY-MM-DD format", err)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "bad_request",
			Message: appErr.Message,
		})
	}

	events, err := h.WeatherStore.GetTideEvents(c.Request().Context(), dateLocal, locCfg.Key)
	if err != nil {
		appErr := utils.NewInternalServerError("Failed to query tide data", err)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "internal_server_error",
			Message: appErr.Message,
		})
	}

	preferredSource := tides.ServingSource()
	selected := selectTideEventsBySource(events, preferredSource)
	if len(selected) == 0 {
		collectCtx, cancel := context.WithTimeout(c.Request().Context(), 8*time.Second)
		err = h.TideCollector.CollectTides(collectCtx, dateLocal, locCfg.Key)
		cancel()
		if err != nil {
			log.Printf("tides collect failed date=%s location=%s error=%v", dateLocal, locCfg.Key, err)
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"error":   "try_later",
				"message": "Tide data is being collected, try later.",
			})
		}

		events, err = h.WeatherStore.GetTideEvents(c.Request().Context(), dateLocal, locCfg.Key)
		if err != nil {
			appErr := utils.NewInternalServerError("Failed to query tide data after collect", err)
			return c.JSON(appErr.Code, utils.ErrorResponse{
				Error:   "internal_server_error",
				Message: appErr.Message,
			})
		}
		if len(events) == 0 {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"error":   "try_later",
				"message": "Tide data not ready yet, try later.",
			})
		}
		selected = selectTideEventsBySource(events, preferredSource)
		if len(selected) == 0 {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"error":   "try_later",
				"message": "Requested tide source not ready yet, try later.",
			})
		}
	}

	return c.JSON(http.StatusOK, buildTidesResponse(dateLocal, locCfg.Key, selected))
}

func buildTidesResponse(dateLocal, locationKey string, events []models.TideEvent) models.TideAPIResponse {
	response := models.TideAPIResponse{
		DateLocal: dateLocal,
		Location:  locationKey,
		Events:    make([]models.TideEventResponse, 0, len(events)),
	}
	if len(events) == 0 {
		return response
	}

	selectedSource := events[0].Source
	response.Source = selectedSource
	response.Confidence = events[0].Confidence
	response.FetchedAt = events[0].FetchedAt.UTC().Format(time.RFC3339)

	for _, ev := range events {
		if ev.Source != selectedSource {
			continue
		}
		response.Events = append(response.Events, models.TideEventResponse{
			Type:      ev.EventType,
			TimeLocal: ev.EventTimeLocal.Format("2006-01-02T15:04:05"),
			HeightM:   ev.HeightM,
		})
	}

	return response
}

func selectTideEventsBySource(events []models.TideEvent, source string) []models.TideEvent {
	if len(events) == 0 {
		return nil
	}

	target := source
	if target == "hybrid" {
		target = "puertos"
	}
	if target == "" {
		target = "open_meteo"
	}

	out := make([]models.TideEvent, 0, len(events))
	for _, ev := range events {
		if ev.Source == target {
			out = append(out, ev)
		}
	}
	return out
}

func (h *Handler) WebcamBigHandler(c echo.Context) error {
	locale, currentPath, languages, messages, gaEnabled, gaMeasurementID := h.getCommonViewData(c)
	data := models.PageData{
		Locale:          locale,
		LocalePrefix:    i18n.LocalePrefix(locale),
		CurrentPath:     currentPath,
		CurrentSection:  "home",
		WebcamImageURL:  h.webcamImageURL(),
		Languages:       languages,
		I18n:            messages,
		GAEnabled:       gaEnabled,
		GAMeasurementID: gaMeasurementID,
	}
	return c.Render(http.StatusOK, "webcam-big.html", data)
}

func (h *Handler) WebcamImageHandler(c echo.Context) error {
	webcamPath := utils.EnvPathOrDefault("WEBCAM_IMAGE_PATH", "public/images/tenelife.jpg")

	info, err := os.Stat(webcamPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Webcam image not found")
	}

	modTime := info.ModTime().UTC()
	etag := fmt.Sprintf("\"%x-%x\"", info.Size(), modTime.UnixNano())

	// Browser may cache, but must always revalidate before using cached copy.
	c.Response().Header().Set("Cache-Control", "public, no-cache, must-revalidate")
	c.Response().Header().Set("ETag", etag)
	c.Response().Header().Set("Last-Modified", modTime.Format(http.TimeFormat))

	if inm := c.Request().Header.Get("If-None-Match"); inm != "" && inm == etag {
		return c.NoContent(http.StatusNotModified)
	}

	if ims := c.Request().Header.Get("If-Modified-Since"); ims != "" {
		if t, parseErr := time.Parse(http.TimeFormat, ims); parseErr == nil && !modTime.After(t) {
			return c.NoContent(http.StatusNotModified)
		}
	}

	return c.File(webcamPath)
}

// HealthCheckHandler provides a health check endpoint
func (h *Handler) HealthCheckHandler(c echo.Context) error {
	// Check database connection
	err := h.WeatherStore.DB.Ping()
	if err != nil {
		appErr := utils.NewInternalServerError("Database connection failed", err)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "database_unhealthy",
			Message: appErr.Message,
		})
	}

	// Check if we can read weather data
	weatherPath := utils.EnvPathOrDefault("WEATHER_JSON_PATH", "public/files/weather.json")

	if _, err := os.Stat(weatherPath); err != nil {
		appErr := utils.NewInternalServerError("Weather data file not accessible", err)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "weather_data_unhealthy",
			Message: appErr.Message,
		})
	}

	// Check if we can read webcam image
	webcamPath := utils.EnvPathOrDefault("WEBCAM_IMAGE_PATH", "public/images/tenelife.jpg")

	if _, err := os.Stat(webcamPath); err != nil {
		appErr := utils.NewInternalServerError("Webcam image not accessible", err)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "webcam_unhealthy",
			Message: appErr.Message,
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":   "healthy",
		"database": "ok",
		"weather":  "ok",
		"webcam":   "ok",
	})
}

func (h *Handler) DailyStatisticsHandler(c echo.Context) error {
	locale, currentPath, languages, messages, gaEnabled, gaMeasurementID := h.getCommonViewData(c)
	stats, err := h.WeatherStore.GetDailyStats(30)
	if err != nil {
		log.Println("Error fetching daily stats:", err)
	}
	data := models.StatsPageData{
		DailyStats:      stats,
		PageTitle:       i18n.T(locale, "daily_statistics"),
		StatsSection:    "daily",
		Locale:          locale,
		LocalePrefix:    i18n.LocalePrefix(locale),
		CurrentPath:     currentPath,
		CurrentSection:  "statistics",
		Languages:       languages,
		I18n:            messages,
		GAEnabled:       gaEnabled,
		GAMeasurementID: gaMeasurementID,
	}
	err = c.Render(http.StatusOK, "daily.html", data)
	if err != nil {
		log.Println("RENDER ERROR:", err)
		return err
	}
	return nil
}

func (h *Handler) WeeklyStatisticsHandler(c echo.Context) error {
	locale, currentPath, languages, messages, gaEnabled, gaMeasurementID := h.getCommonViewData(c)
	stats, err := h.WeatherStore.GetWeeklyStats()
	if err != nil {
		log.Println("Error fetching weekly stats:", err)
	}
	data := models.StatsPageData{
		WeeklyStats:     stats,
		PageTitle:       i18n.T(locale, "weekly_statistics"),
		StatsSection:    "weekly",
		Locale:          locale,
		LocalePrefix:    i18n.LocalePrefix(locale),
		CurrentPath:     currentPath,
		CurrentSection:  "statistics",
		Languages:       languages,
		I18n:            messages,
		GAEnabled:       gaEnabled,
		GAMeasurementID: gaMeasurementID,
	}
	return c.Render(http.StatusOK, "weekly.html", data)
}

func (h *Handler) RecentStatisticsHandler(c echo.Context) error {
	locale, currentPath, languages, messages, gaEnabled, gaMeasurementID := h.getCommonViewData(c)
	stats, err := h.WeatherStore.GetDailyStats(10)
	if err != nil {
		log.Println("Error fetching recent stats:", err)
	}
	stats = filterDailyStatsWithValues(stats)
	data := models.StatsPageData{
		DailyStats:      stats,
		PageTitle:       i18n.T(locale, "recent_statistics"),
		StatsSection:    "recent",
		Locale:          locale,
		LocalePrefix:    i18n.LocalePrefix(locale),
		CurrentPath:     currentPath,
		CurrentSection:  "statistics",
		Languages:       languages,
		I18n:            messages,
		GAEnabled:       gaEnabled,
		GAMeasurementID: gaMeasurementID,
	}
	return c.Render(http.StatusOK, "recent.html", data)
}

func (h *Handler) GetPWSLatestHandler(c echo.Context) error {
	points, err := h.WeatherStore.GetPWSLatestPoints()
	if err != nil {
		appErr := utils.NewInternalServerError("Failed to get PWS latest data", err)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "internal_server_error",
			Message: appErr.Message,
		})
	}
	return c.JSON(http.StatusOK, points)
}

func (h *Handler) MonthlyStatisticsHandler(c echo.Context) error {
	locale, currentPath, languages, messages, gaEnabled, gaMeasurementID := h.getCommonViewData(c)
	stats, err := h.WeatherStore.GetMonthlyStats(12)
	if err != nil {
		log.Println("Error fetching monthly stats:", err)
	}
	data := models.StatsPageData{
		MonthlyStats:    stats,
		PageTitle:       i18n.T(locale, "monthly_statistics"),
		StatsSection:    "monthly",
		Locale:          locale,
		LocalePrefix:    i18n.LocalePrefix(locale),
		CurrentPath:     currentPath,
		CurrentSection:  "statistics",
		Languages:       languages,
		I18n:            messages,
		GAEnabled:       gaEnabled,
		GAMeasurementID: gaMeasurementID,
	}
	return c.Render(http.StatusOK, "monthly.html", data)
}

func (h *Handler) AnnualStatisticsHandler(c echo.Context) error {
	locale, currentPath, languages, messages, gaEnabled, gaMeasurementID := h.getCommonViewData(c)
	stats, err := h.WeatherStore.GetAnnualStats()
	if err != nil {
		log.Println("Error fetching annual stats:", err)
	}
	data := models.StatsPageData{
		AnnualStats:     stats,
		PageTitle:       i18n.T(locale, "annual_statistics"),
		StatsSection:    "annual",
		Locale:          locale,
		LocalePrefix:    i18n.LocalePrefix(locale),
		CurrentPath:     currentPath,
		CurrentSection:  "statistics",
		Languages:       languages,
		I18n:            messages,
		GAEnabled:       gaEnabled,
		GAMeasurementID: gaMeasurementID,
	}
	return c.Render(http.StatusOK, "annual.html", data)
}

// API for Statistics Charts

// API for Statistics Charts

func (h *Handler) GetDailyDataHandler(c echo.Context) error {
	daysStr := c.QueryParam("days")
	days := 7
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	endDate := c.QueryParam("end_date")
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	startTs := time.Now()
	if t, err := time.Parse("2006-01-02", endDate); err == nil {
		startTs = t
	}
	startDate := startTs.AddDate(0, 0, -(days - 1)).Format("2006-01-02")

	stats, err := h.WeatherStore.GetDailyStatsByRange(startDate, endDate)
	if err != nil {
		appErr := utils.NewInternalServerError("Failed to get daily stats", err)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "internal_server_error",
			Message: appErr.Message,
		})
	}

	response := models.DailyChartResponse{}
	for _, s := range stats {
		response.Labels = append(response.Labels, formatDayMonthLabel(s.Date))
		response.Datasets.AvgTemperature = append(response.Datasets.AvgTemperature, s.AvgTemperature)
		response.Datasets.MinTemperature = append(response.Datasets.MinTemperature, s.MinTemperature)
		response.Datasets.MaxTemperature = append(response.Datasets.MaxTemperature, s.MaxTemperature)
		response.Datasets.AvgPressure = append(response.Datasets.AvgPressure, s.AvgPressure)
		response.Datasets.AvgHumidity = append(response.Datasets.AvgHumidity, s.AvgHumidity)
		response.Datasets.SeaTemperature = append(response.Datasets.SeaTemperature, s.SeaTemperature)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetMonthlyDailyDataHandler(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")

	year, _ := strconv.Atoi(yearStr)
	month, _ := strconv.Atoi(monthStr)

	if year == 0 {
		year = time.Now().Year()
	}
	if month == 0 {
		month = int(time.Now().Month())
	}

	startDate := fmt.Sprintf("%04d-%02d-01", year, month)
	// Calculate end of month
	t := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC)
	endDate := t.Format("2006-01-02")

	stats, err := h.WeatherStore.GetDailyStatsByRange(startDate, endDate)
	if err != nil {
		appErr := utils.NewInternalServerError("Failed to get monthly daily stats", err)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "internal_server_error",
			Message: appErr.Message,
		})
	}

	response := models.DailyChartResponse{}
	daysInMonth := t.Day()

	// Pre-fill with nulls/zeros or just loop through records
	// To match Laravel's behavior of showing all days in month:
	for day := 1; day <= daysInMonth; day++ {
		dateStr := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
		response.Labels = append(response.Labels, strconv.Itoa(day))

		var found *models.WeatherDaily
		for _, s := range stats {
			if s.Date == dateStr {
				found = &s
				break
			}
		}

		if found != nil {
			response.Datasets.AvgTemperature = append(response.Datasets.AvgTemperature, found.AvgTemperature)
			response.Datasets.AvgPressure = append(response.Datasets.AvgPressure, found.AvgPressure)
			response.Datasets.AvgHumidity = append(response.Datasets.AvgHumidity, found.AvgHumidity)
			response.Datasets.SeaTemperature = append(response.Datasets.SeaTemperature, found.SeaTemperature)
		} else {
			response.Datasets.AvgTemperature = append(response.Datasets.AvgTemperature, nil)
			response.Datasets.AvgPressure = append(response.Datasets.AvgPressure, nil)
			response.Datasets.AvgHumidity = append(response.Datasets.AvgHumidity, nil)
			response.Datasets.SeaTemperature = append(response.Datasets.SeaTemperature, nil)
		}
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetMonthlyDataHandler(c echo.Context) error {
	stats, err := h.WeatherStore.GetMonthlyStats(12)
	if err != nil {
		appErr := utils.NewInternalServerError("Failed to get monthly stats", err)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "internal_server_error",
			Message: appErr.Message,
		})
	}

	response := models.DailyChartResponse{} // Reusing structure or could use Generic
	for i := len(stats) - 1; i >= 0; i-- {  // Reverse to show chronological
		s := stats[i]
		response.Labels = append(response.Labels, fmt.Sprintf("%d/%d", s.Month, s.Year))
		response.Datasets.SeaTemperature = append(response.Datasets.SeaTemperature, s.AvgSeaTemperature)
		response.Datasets.AvgTemperature = append(response.Datasets.AvgTemperature, s.AvgTemperature)
		response.Datasets.AvgPressure = append(response.Datasets.AvgPressure, s.AvgPressure)
		response.Datasets.AvgHumidity = append(response.Datasets.AvgHumidity, s.AvgHumidity)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetWeeklyDataHandler(c echo.Context) error {
	stats, err := h.WeatherStore.GetWeeklyStats()
	if err != nil {
		appErr := utils.NewInternalServerError("Failed to get weekly stats", err)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "internal_server_error",
			Message: appErr.Message,
		})
	}

	response := models.DailyChartResponse{}
	// Take only last 20 weeks and reverse
	limit := 20
	if len(stats) < limit {
		limit = len(stats)
	}
	for i := limit - 1; i >= 0; i-- {
		s := stats[i]
		label := fmt.Sprintf("%d/W%d", s.Year, s.Week)
		if s.WeekStart != "" && s.WeekEnd != "" {
			label = fmt.Sprintf("%d/%d (%s-%s)", s.Week, s.Year, formatDayMonthLabel(s.WeekStart), formatDayMonthLabel(s.WeekEnd))
		}
		response.Labels = append(response.Labels, label)
		response.Datasets.SeaTemperature = append(response.Datasets.SeaTemperature, s.AvgSeaTemperature)
		response.Datasets.AvgTemperature = append(response.Datasets.AvgTemperature, s.AvgTemperature)
		response.Datasets.AvgPressure = append(response.Datasets.AvgPressure, s.AvgPressure)
		response.Datasets.AvgHumidity = append(response.Datasets.AvgHumidity, s.AvgHumidity)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetAnnualDataHandler(c echo.Context) error {
	stats, err := h.WeatherStore.GetAnnualStats()
	if err != nil {
		appErr := utils.NewInternalServerError("Failed to get annual stats", err)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "internal_server_error",
			Message: appErr.Message,
		})
	}

	response := models.DailyChartResponse{}
	for i := len(stats) - 1; i >= 0; i-- {
		s := stats[i]
		response.Labels = append(response.Labels, fmt.Sprintf("%d/%d", s.Month, s.Year))
		response.Datasets.SeaTemperature = append(response.Datasets.SeaTemperature, s.AvgSeaTemperature)
		response.Datasets.AvgTemperature = append(response.Datasets.AvgTemperature, s.AvgTemperature)
		response.Datasets.AvgPressure = append(response.Datasets.AvgPressure, s.AvgPressure)
		response.Datasets.AvgHumidity = append(response.Datasets.AvgHumidity, s.AvgHumidity)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) StoreSeaTemperatureHandler(c echo.Context) error {
	type SeaTempRequest struct {
		Date        string  `json:"date"`
		Temperature float64 `json:"temperature"`
	}

	req := new(SeaTempRequest)
	if err := c.Bind(req); err != nil {
		appErr := utils.NewBadRequestError("Invalid request format", err)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "bad_request",
			Message: appErr.Message,
		})
	}

	if req.Date == "" || req.Temperature < -10 || req.Temperature > 50 {
		appErr := utils.NewBadRequestError("Invalid date or temperature value", nil)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "bad_request",
			Message: appErr.Message,
		})
	}
	if _, err := time.Parse("2006-01-02", req.Date); err != nil {
		appErr := utils.NewBadRequestError("Date must be in YYYY-MM-DD format", err)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "bad_request",
			Message: appErr.Message,
		})
	}

	err := h.WeatherStore.StoreSeaTemperature(req.Date, req.Temperature)
	if err != nil {
		appErr := utils.NewInternalServerError("Failed to store sea temperature", err)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "internal_server_error",
			Message: appErr.Message,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Sea temperature saved successfully",
	})
}

func formatDayMonthLabel(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	layouts := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		time.RFC3339,
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t.Format("2.1.")
		}
	}

	if len(raw) >= 10 {
		if t, err := time.Parse("2006-01-02", raw[:10]); err == nil {
			return t.Format("2.1.")
		}
	}

	return raw
}

func (h *Handler) CameraUploadHandler(c echo.Context) error {
	// Protect upload endpoint with a static token.
	// Accepted locations: Authorization: Bearer <token> or X-API-Key header.
	if !isAuthorizedUploadRequest(c) {
		appErr := utils.NewUnauthorizedError("Unauthorized upload request", nil)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "unauthorized",
			Message: appErr.Message,
		})
	}

	// 2. Determine where to save
	webcamPath := utils.EnvPathOrDefault("WEBCAM_IMAGE_PATH", "public/images/tenelife.jpg")

	// 3. Handle different upload methods (multipart vs raw body)
	file, err := c.FormFile("image")
	if err == nil {
		// Multipart file
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.Create(webcamPath)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
	} else {
		// Try raw body
		body, err := io.ReadAll(io.LimitReader(c.Request().Body, 10*1024*1024))
		if err != nil || len(body) == 0 {
			appErr := utils.NewBadRequestError("No image data received", err)
			return c.JSON(appErr.Code, utils.ErrorResponse{
				Error:   "bad_request",
				Message: appErr.Message,
			})
		}
		if err := os.WriteFile(webcamPath, body, 0644); err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"path":    webcamPath,
	})
}

func isAuthorizedUploadRequest(c echo.Context) bool {
	expectedToken := os.Getenv("CAMERA_UPLOAD_TOKEN")
	if expectedToken == "" {
		// Secure-by-default: if token is not configured, reject uploads.
		return false
	}

	gotToken := strings.TrimSpace(c.Request().Header.Get("X-API-Key"))
	if gotToken == "" {
		auth := strings.TrimSpace(c.Request().Header.Get("Authorization"))
		if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			gotToken = strings.TrimSpace(auth[7:])
		}
	}
	if gotToken == "" {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(gotToken), []byte(expectedToken)) == 1
}
