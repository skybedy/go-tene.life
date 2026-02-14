package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"sync"

	"github.com/labstack/echo/v4"
	"github.com/skybedy/laravel-tene.life/internal/models"
	"github.com/skybedy/laravel-tene.life/internal/store"
)

// TemplateRenderer implements Echo's Renderer interface
type TemplateRenderer struct {
	Templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

type Handler struct {
	WeatherStore *store.WeatherStore
	weatherCache *models.WeatherData
	cacheMu      sync.RWMutex
	lastCache    time.Time
}

func NewHandler(ws *store.WeatherStore) *Handler {
	return &Handler{
		WeatherStore: ws,
	}
}

func (h *Handler) IndexHandler(c echo.Context) error {
	// 1. Get Weather Data (with simple caching)
	h.cacheMu.RLock()
	weather := h.weatherCache
	cacheAge := time.Since(h.lastCache)
	h.cacheMu.RUnlock()

	// Refresh cache if older than 30 seconds or empty
	if weather == nil || cacheAge > 30*time.Second {
		h.cacheMu.Lock()
		// Double check after acquiring lock
		if h.weatherCache == nil || time.Since(h.lastCache) > 30*time.Second {
			weatherPath := os.Getenv("WEATHER_JSON_PATH")
			if weatherPath == "" {
				weatherPath = "public/files/weather.json"
			}
			file, err := os.Open(weatherPath)
			if err == nil {
				defer file.Close()
				decoder := json.NewDecoder(file)
				newWeather := &models.WeatherData{}
				if err := decoder.Decode(newWeather); err == nil {
					h.weatherCache = newWeather
					h.lastCache = time.Now()
					weather = newWeather
				} else {
					log.Println("Error decoding weather.json:", err)
				}
			} else {
				log.Println("Error opening weather.json:", err)
			}
		} else {
			weather = h.weatherCache
		}
		h.cacheMu.Unlock()
	}

	// 2. Get Sea Temperature from DB via Store
	seaTemp, err := h.WeatherStore.GetLatestSeaTemperature()
	if err != nil {
		log.Println("Error fetching sea temperature (not critical):", err)
	}

	// 3. Format Date/Time
	ts := time.Now()
	if weather != nil && weather.Timestamp > 0 {
		ts = time.Unix(weather.Timestamp, 0)
	}

	formattedDate := ts.Format("2. 1. 2006")
	formattedTime := ts.Format("15:04")

	var seaTempVal float64
	if seaTemp != nil {
		seaTempVal = *seaTemp
	}

	data := models.PageData{
		Weather:           weather,
		SeaTemperature:    seaTemp,
		SeaTemperatureVal: seaTempVal,
		FormattedDate:     formattedDate,
		FormattedTime:     formattedTime,
		PageTitle:         "",
	}

	return c.Render(http.StatusOK, "index.html", data)
}

func (h *Handler) GetHourlyDataHandler(c echo.Context) error {
	date := c.QueryParam("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	data, err := h.WeatherStore.GetHourlyData(date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
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

func (h *Handler) WebcamBigHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "webcam-big.html", nil)
}

func (h *Handler) WebcamImageHandler(c echo.Context) error {
	webcamPath := os.Getenv("WEBCAM_IMAGE_PATH")
	if webcamPath == "" {
		webcamPath = "public/images/tenelife.jpg"
	}

	// Disable caching for the webcam image to ensure it's always fresh
	c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Response().Header().Set("Pragma", "no-cache")
	c.Response().Header().Set("Expires", "0")

	return c.File(webcamPath)
}

func (h *Handler) DailyStatisticsHandler(c echo.Context) error {
	stats, err := h.WeatherStore.GetDailyStats(30)
	if err != nil {
		log.Println("Error fetching daily stats:", err)
	}
	data := models.StatsPageData{
		DailyStats: stats,
		PageTitle:  "Denní statistiky",
	}
	return c.Render(http.StatusOK, "daily.html", data)
}

func (h *Handler) WeeklyStatisticsHandler(c echo.Context) error {
	stats, err := h.WeatherStore.GetWeeklyStats()
	if err != nil {
		log.Println("Error fetching weekly stats:", err)
	}
	data := models.StatsPageData{
		WeeklyStats: stats,
		PageTitle:   "Týdenní statistiky",
	}
	return c.Render(http.StatusOK, "weekly.html", data)
}

func (h *Handler) MonthlyStatisticsHandler(c echo.Context) error {
	stats, err := h.WeatherStore.GetMonthlyStats(12)
	if err != nil {
		log.Println("Error fetching monthly stats:", err)
	}
	data := models.StatsPageData{
		MonthlyStats: stats,
		PageTitle:    "Měsíční statistiky",
	}
	return c.Render(http.StatusOK, "monthly.html", data)
}

func (h *Handler) AnnualStatisticsHandler(c echo.Context) error {
	stats, err := h.WeatherStore.GetAnnualStats()
	if err != nil {
		log.Println("Error fetching annual stats:", err)
	}
	data := models.StatsPageData{
		AnnualStats: stats,
		PageTitle:   "Roční statistiky",
	}
	return c.Render(http.StatusOK, "annual.html", data)
}

// API for Statistics Charts

func (h *Handler) GetDailyDataHandler(c echo.Context) error {
	stats, err := h.WeatherStore.GetDailyStats(7) // Default to 7 days
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	// Format for Chart.js
	return c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetMonthlyDataHandler(c echo.Context) error {
	stats, err := h.WeatherStore.GetMonthlyStats(12)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetWeeklyDataHandler(c echo.Context) error {
	stats, err := h.WeatherStore.GetWeeklyStats()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetAnnualDataHandler(c echo.Context) error {
	stats, err := h.WeatherStore.GetAnnualStats()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, stats)
}
