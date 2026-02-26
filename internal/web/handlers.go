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
	"github.com/skybedy/laravel-tene.life/internal/utils"
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
	seaTempCache *float64
	cacheMu      sync.RWMutex
	lastCache    time.Time
	cacheTimeout time.Duration
}

func NewHandler(ws *store.WeatherStore) *Handler {
	return &Handler{
		WeatherStore: ws,
		cacheTimeout: 30 * time.Second, // Default cache timeout
	}
}

func (h *Handler) IndexHandler(c echo.Context) error {
	// 1. Get Weather Data (with improved caching)
	weather, seaTemp, err := h.getCachedWeatherData()
	if err != nil {
		log.Printf("Error getting weather data: %v", err)
		// Continue with cached data if available, or empty data
		if weather == nil && seaTemp == nil {
			return c.Render(http.StatusOK, "index.html", models.PageData{
				FormattedDate: time.Now().Format("2. 1. 2006"),
				FormattedTime: time.Now().Format("15:04"),
			})
		}
	}

	// 2. Format Date/Time
	ts := time.Now()
	if weather != nil && weather.Timestamp > 0 {
		ts = time.Unix(weather.Timestamp, 0)
	}

	var seaTempVal float64
	if seaTemp != nil {
		seaTempVal = *seaTemp
	}

	data := models.PageData{
		Weather:           weather,
		SeaTemperature:    seaTemp,
		SeaTemperatureVal: seaTempVal,
		FormattedDate:     ts.Format("2. 1. 2006"),
		FormattedTime:     ts.Format("15:04"),
		PageTitle:         "",
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
		defer h.cacheMu.Unlock()
		
		// Double check after acquiring lock
		if h.weatherCache == nil || h.seaTempCache == nil || time.Since(h.lastCache) > h.cacheTimeout {
			// Update Weather JSON
			weatherPath := os.Getenv("WEATHER_JSON_PATH")
			if weatherPath == "" {
				weatherPath = "public/files/weather.json"
			}
			
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
			
			// Update Sea Temp from DB
			newSeaTemp, err := h.WeatherStore.GetLatestSeaTemperature()
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

func (h *Handler) GetHourlyDataHandler(c echo.Context) error {
	date := c.QueryParam("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
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

func (h *Handler) WebcamBigHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "webcam-big.html", nil)
}

func (h *Handler) WebcamImageHandler(c echo.Context) error {
	webcamPath := os.Getenv("WEBCAM_IMAGE_PATH")
	if webcamPath == "" {
		webcamPath = "public/images/tenelife.jpg"
	}

	// Validate webcam path
	if !utils.IsSafePath(webcamPath) {
		appErr := utils.NewForbiddenError("Invalid webcam image path", nil)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "forbidden",
			Message: appErr.Message,
		})
	}

	// Disable caching for the webcam image to ensure it's always fresh
	c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Response().Header().Set("Pragma", "no-cache")
	c.Response().Header().Set("Expires", "0")

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
	weatherPath := os.Getenv("WEATHER_JSON_PATH")
	if weatherPath == "" {
		weatherPath = "public/files/weather.json"
	}
	
	if _, err := os.Stat(weatherPath); err != nil {
		appErr := utils.NewInternalServerError("Weather data file not accessible", err)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "weather_data_unhealthy",
			Message: appErr.Message,
		})
	}

	// Check if we can read webcam image
	webcamPath := os.Getenv("WEBCAM_IMAGE_PATH")
	if webcamPath == "" {
		webcamPath = "public/images/tenelife.jpg"
	}
	
	if _, err := os.Stat(webcamPath); err != nil {
		appErr := utils.NewInternalServerError("Webcam image not accessible", err)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "webcam_unhealthy",
			Message: appErr.Message,
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":  "healthy",
		"database": "ok",
		"weather":  "ok",
		"webcam":   "ok",
	})
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
		appErr := utils.NewInternalServerError("Failed to get daily stats", err)
		return c.JSON(appErr.Code, utils.ErrorResponse{
			Error:   "internal_server_error",
			Message: appErr.Message,
		})
	}
	// Format for Chart.js
	return c.JSON(http.StatusOK, stats)
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
	return c.JSON(http.StatusOK, stats)
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
	return c.JSON(http.StatusOK, stats)
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
	return c.JSON(http.StatusOK, stats)
}
