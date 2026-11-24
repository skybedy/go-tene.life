package web

import (
	"crypto/sha256"
	"crypto/subtle"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/skybedy/laravel-tene.life/internal/store"
)

const seaTemperatureFormPath = "/mereni-more"
const measurementTimeLayout = "2006-01-02T15:04"

type seaTemperatureForm struct {
	CSRF, Temperature, MeasuredAt, Error string
	Saved                                bool
}

func (h *Handler) RegisterSeaTemperatureForm(e *echo.Echo, password string) {
	if password == "" {
		return
	}

	passwordHash := sha256.Sum256([]byte(password))
	auth := middleware.BasicAuth(func(username, password string, _ echo.Context) (bool, error) {
		candidate := sha256.Sum256([]byte(password))
		return subtle.ConstantTimeCompare(passwordHash[:], candidate[:]) == 1 && username == "mereni", nil
	})
	e.GET(seaTemperatureFormPath, h.showSeaTemperatureForm, auth)
	e.POST(seaTemperatureFormPath, h.saveSeaTemperatureForm, auth)
}

func (h *Handler) showSeaTemperatureForm(c echo.Context) error {
	now, err := store.ParseMeasuredAtOrDate(time.Now().Format(time.RFC3339), "")
	if err != nil {
		return err
	}
	return renderSeaTemperatureForm(c, http.StatusOK, seaTemperatureForm{
		MeasuredAt: now.Format(measurementTimeLayout),
		Saved:      c.QueryParam("ulozeno") == "1",
	})
}

func (h *Handler) saveSeaTemperatureForm(c echo.Context) error {
	form := seaTemperatureForm{
		Temperature: strings.TrimSpace(c.FormValue("temperature")),
		MeasuredAt:  c.FormValue("measured_at"),
	}

	temperature, err := strconv.ParseFloat(strings.ReplaceAll(form.Temperature, ",", "."), 64)
	if err != nil || math.IsNaN(temperature) || temperature < -10 || temperature > 50 || math.Abs(temperature*10-math.Round(temperature*10)) > 1e-9 {
		form.Error = "Zadejte teplotu od −10 do 50 °C, nejvýše na jedno desetinné místo."
		return renderSeaTemperatureForm(c, http.StatusUnprocessableEntity, form)
	}

	measuredAt, err := store.ParseMeasuredAtOrDate(form.MeasuredAt, "")
	if err != nil || measuredAt.Year() < 1000 || measuredAt.Format(measurementTimeLayout) != form.MeasuredAt {
		form.Error = "Zadejte platné datum a čas měření na Tenerife."
		return renderSeaTemperatureForm(c, http.StatusUnprocessableEntity, form)
	}

	if err := h.WeatherStore.StoreWaterTemperatureMeasurement(measuredAt, temperature, "manual", nil); err != nil {
		c.Logger().Errorf("could not save sea temperature measurement: %v", err)
		form.Error = "Měření se nepodařilo uložit. Zkuste to znovu."
		return renderSeaTemperatureForm(c, http.StatusInternalServerError, form)
	}

	c.Response().Header().Set("Cache-Control", "no-store")
	return c.Redirect(http.StatusSeeOther, seaTemperatureFormPath+"?ulozeno=1")
}

func renderSeaTemperatureForm(c echo.Context, status int, form seaTemperatureForm) error {
	form.CSRF, _ = c.Get("csrf").(string)
	c.Response().Header().Set("Cache-Control", "no-store")
	c.Response().Header().Set("X-Robots-Tag", "noindex")
	return c.Render(status, "sea-temperature-form.html", form)
}
