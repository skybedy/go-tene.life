package web

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

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
}

func NewHandler(ws *store.WeatherStore) *Handler {
	return &Handler{
		WeatherStore: ws,
	}
}

func (h *Handler) IndexHandler(c echo.Context) error {
	// 1. Get Weather Data from JSON
	var weather *models.WeatherData
	file, err := os.Open("public/files/weather.json")
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		weather = &models.WeatherData{}
		if err := decoder.Decode(weather); err != nil {
			log.Println("Error decoding weather.json:", err)
			weather = nil
		}
	} else {
		log.Println("Error opening weather.json:", err)
	}

	// 2. Get Sea Temperature from DB via Store
	seaTemp, err := h.WeatherStore.GetLatestSeaTemperature()
	if err != nil {
		log.Println("Error fetching sea temperature (not critical):", err)
	}

	// 3. Format Date/Time
	// Use timestamp from weather data or current time
	ts := time.Now()
	if weather != nil && weather.Timestamp > 0 {
		ts = time.Unix(weather.Timestamp, 0)
	}

	// Use explicit Czech timezone or UTC+1/PCT if server is local (assumed local for PoC)
	formattedDate := ts.Format("2. 1. 2006")
	formattedTime := ts.Format("15:04")

	// Dereference seaTemp for template if it exists
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
	}

	return c.Render(http.StatusOK, "index.html", data)
}
