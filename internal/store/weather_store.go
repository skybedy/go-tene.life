package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/skybedy/laravel-tene.life/internal/models"
)

type WeatherStore struct {
	DB *sql.DB
}

func NewWeatherStore(db *sql.DB) *WeatherStore {
	return &WeatherStore{DB: db}
}

func canaryLocation() *time.Location {
	loc, err := time.LoadLocation("Atlantic/Canary")
	if err != nil {
		return time.UTC
	}
	return loc
}

func (s *WeatherStore) GetLatestSeaTemperature(referenceDate string) (*float64, string, error) {
	var seaTemp *float64
	var temp float64
	var recDateTime time.Time

	// Prefer the latest point-in-time measurement from the dedicated table.
	query := "SELECT temperature, measured_at FROM water_temperatures WHERE measured_at <= ? ORDER BY measured_at DESC LIMIT 1"
	endOfDay := referenceDate + " 23:59:59"
	err := s.DB.QueryRow(query, endOfDay).Scan(&temp, &recDateTime)
	if err == nil {
		seaTemp = &temp
		return seaTemp, recDateTime.Format("2006-01-02 15:04:05"), nil
	}

	// Fallback to latest point-in-time measurement regardless of date.
	query = "SELECT temperature, measured_at FROM water_temperatures ORDER BY measured_at DESC LIMIT 1"
	err = s.DB.QueryRow(query).Scan(&temp, &recDateTime)
	if err == nil {
		seaTemp = &temp
		return seaTemp, recDateTime.Format("2006-01-02 15:04:05"), nil
	}

	// Transitional fallback to legacy daily column while old data still exists there.
	var recDate string
	query = "SELECT sea_temperature, date FROM weather_daily WHERE date <= ? AND sea_temperature IS NOT NULL ORDER BY date DESC LIMIT 1"
	err = s.DB.QueryRow(query, referenceDate).Scan(&temp, &recDate)
	if err == nil {
		seaTemp = &temp
		return seaTemp, recDate + " 10:00:00", nil
	}

	query = "SELECT sea_temperature, date FROM weather_daily WHERE sea_temperature IS NOT NULL ORDER BY date DESC LIMIT 1"
	err = s.DB.QueryRow(query).Scan(&temp, &recDate)
	if err == nil {
		seaTemp = &temp
		return seaTemp, recDate + " 10:00:00", nil
	}

	return nil, "", nil // Return nil if no temperature found, not strictly an error for the view
}

func (s *WeatherStore) GetLatestWeather() (*models.WeatherData, error) {
	query := `SELECT measured_at, temperature, pressure, humidity
	          FROM weather
	          ORDER BY measured_at DESC
	          LIMIT 1`

	var measuredAt time.Time
	var weather models.WeatherData
	if err := s.DB.QueryRow(query).Scan(&measuredAt, &weather.Temperature, &weather.Pressure, &weather.Humidity); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	weather.Timestamp = measuredAt.In(canaryLocation()).Unix()
	return &weather, nil
}

func (s *WeatherStore) GetHourlyData(date string) ([]models.WeatherHourly, error) {
	var results []models.WeatherHourly

	query := "SELECT date, hour, avg_temperature, avg_pressure, avg_humidity FROM weather_hourly WHERE date = ? ORDER BY hour ASC"
	rows, err := s.DB.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var h models.WeatherHourly
		err := rows.Scan(&h.Date, &h.Hour, &h.AvgTemperature, &h.AvgPressure, &h.AvgHumidity)
		if err != nil {
			return nil, err
		}
		results = append(results, h)
	}

	return results, nil
}

func (s *WeatherStore) GetDailyStats(limit int) ([]models.WeatherDaily, error) {
	var results []models.WeatherDaily
	query := `SELECT date, sea_temperature, avg_temperature, min_temperature, max_temperature, 
	                 avg_pressure, min_pressure, max_pressure, avg_humidity, min_humidity, max_humidity, samples_count 
	          FROM weather_daily ORDER BY date DESC LIMIT ?`
	rows, err := s.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d models.WeatherDaily
		err := rows.Scan(&d.Date, &d.SeaTemperature, &d.AvgTemperature, &d.MinTemperature, &d.MaxTemperature,
			&d.AvgPressure, &d.MinPressure, &d.MaxPressure, &d.AvgHumidity, &d.MinHumidity, &d.MaxHumidity, &d.SamplesCount)
		if err != nil {
			return nil, err
		}
		results = append(results, d)
	}
	return results, nil
}

func (s *WeatherStore) GetDailyStatsByRange(startDate, endDate string) ([]models.WeatherDaily, error) {
	var results []models.WeatherDaily
	query := `SELECT date, sea_temperature, avg_temperature, min_temperature, max_temperature, 
	                 avg_pressure, min_pressure, max_pressure, avg_humidity, min_humidity, max_humidity, samples_count 
	          FROM weather_daily 
	          WHERE date BETWEEN ? AND ? 
	          ORDER BY date ASC`
	rows, err := s.DB.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d models.WeatherDaily
		err := rows.Scan(&d.Date, &d.SeaTemperature, &d.AvgTemperature, &d.MinTemperature, &d.MaxTemperature,
			&d.AvgPressure, &d.MinPressure, &d.MaxPressure, &d.AvgHumidity, &d.MinHumidity, &d.MaxHumidity, &d.SamplesCount)
		if err != nil {
			return nil, err
		}
		results = append(results, d)
	}
	return results, nil
}

func (s *WeatherStore) GetWeeklyStats() ([]models.WeatherWeekly, error) {
	var results []models.WeatherWeekly
	query := `SELECT year, week, week_start, week_end, avg_sea_temperature, avg_temperature, min_temperature, max_temperature, 
	                 avg_pressure, min_pressure, max_pressure, avg_humidity, min_humidity, max_humidity, samples_count 
	          FROM weather_weekly ORDER BY year DESC, week DESC`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var w models.WeatherWeekly
		err := rows.Scan(&w.Year, &w.Week, &w.WeekStart, &w.WeekEnd, &w.AvgSeaTemperature, &w.AvgTemperature, &w.MinTemperature, &w.MaxTemperature,
			&w.AvgPressure, &w.MinPressure, &w.MaxPressure, &w.AvgHumidity, &w.MinHumidity, &w.MaxHumidity, &w.SamplesCount)
		if err != nil {
			return nil, err
		}
		results = append(results, w)
	}
	return results, nil
}

func (s *WeatherStore) GetMonthlyStats(limit int) ([]models.WeatherMonthly, error) {
	var results []models.WeatherMonthly
	query := `SELECT year, month, avg_sea_temperature, avg_temperature, min_temperature, max_temperature, 
	                 avg_pressure, min_pressure, max_pressure, avg_humidity, min_humidity, max_humidity, samples_count 
	          FROM weather_monthly ORDER BY year DESC, month DESC LIMIT ?`
	rows, err := s.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m models.WeatherMonthly
		err := rows.Scan(&m.Year, &m.Month, &m.AvgSeaTemperature, &m.AvgTemperature, &m.MinTemperature, &m.MaxTemperature,
			&m.AvgPressure, &m.MinPressure, &m.MaxPressure, &m.AvgHumidity, &m.MinHumidity, &m.MaxHumidity, &m.SamplesCount)
		if err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, nil
}

func (s *WeatherStore) GetAnnualStats() ([]models.WeatherMonthly, error) {
	var results []models.WeatherMonthly
	query := `SELECT year, month, avg_sea_temperature, avg_temperature, min_temperature, max_temperature, 
	                 avg_pressure, min_pressure, max_pressure, avg_humidity, min_humidity, max_humidity, samples_count 
	          FROM weather_monthly ORDER BY year DESC, month DESC`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m models.WeatherMonthly
		err := rows.Scan(&m.Year, &m.Month, &m.AvgSeaTemperature, &m.AvgTemperature, &m.MinTemperature, &m.MaxTemperature,
			&m.AvgPressure, &m.MinPressure, &m.MaxPressure, &m.AvgHumidity, &m.MinHumidity, &m.MaxHumidity, &m.SamplesCount)
		if err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, nil
}

func (s *WeatherStore) StoreSeaTemperature(date string, temp float64) error {
	query := `INSERT INTO weather_daily (date, sea_temperature) 
	          VALUES (?, ?) 
	          ON DUPLICATE KEY UPDATE sea_temperature = VALUES(sea_temperature)`
	_, err := s.DB.Exec(query, date, temp)
	return err
}

func (s *WeatherStore) StoreWaterTemperatureMeasurement(measuredAt time.Time, temp float64, source string, note *string) error {
	if strings.TrimSpace(source) == "" {
		source = "manual"
	}

	query := `INSERT INTO water_temperatures (measured_at, temperature, source, note)
	          VALUES (?, ?, ?, ?)`
	_, err := s.DB.Exec(query, measuredAt.In(canaryLocation()), temp, source, note)
	return err
}

func (s *WeatherStore) GetWaterTemperaturesByRange(start, end time.Time, limit int) ([]models.WaterTemperatureMeasurement, error) {
	query := `SELECT id, measured_at, temperature, source, note, legacy_weather_daily_id
	          FROM water_temperatures
	          WHERE measured_at BETWEEN ? AND ?
	          ORDER BY measured_at ASC`
	args := []interface{}{start.In(canaryLocation()), end.In(canaryLocation())}
	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.WaterTemperatureMeasurement, 0)
	for rows.Next() {
		var rec models.WaterTemperatureMeasurement
		var measuredAt time.Time
		if err := rows.Scan(&rec.ID, &measuredAt, &rec.Temperature, &rec.Source, &rec.Note, &rec.LegacyWeatherDailyID); err != nil {
			return nil, err
		}
		rec.MeasuredAt = measuredAt.In(canaryLocation()).Format("2006-01-02 15:04:05")
		out = append(out, rec)
	}
	return out, nil
}

func (s *WeatherStore) BuildWaterTemperatureHistory(start, end time.Time, limit int) (*models.WaterTemperatureHistoryResponse, error) {
	records, err := s.GetWaterTemperaturesByRange(start, end, limit)
	if err != nil {
		return nil, err
	}
	resp := &models.WaterTemperatureHistoryResponse{
		Labels:       make([]string, 0, len(records)),
		Temperatures: make([]*float64, 0, len(records)),
		MeasuredAt:   make([]string, 0, len(records)),
	}

	for _, rec := range records {
		resp.MeasuredAt = append(resp.MeasuredAt, rec.MeasuredAt)
		resp.Temperatures = append(resp.Temperatures, rec.Temperature)

		if ts, err := time.ParseInLocation("2006-01-02 15:04:05", rec.MeasuredAt, canaryLocation()); err == nil {
			resp.Labels = append(resp.Labels, ts.Format("2.1. 15:04"))
		} else {
			resp.Labels = append(resp.Labels, rec.MeasuredAt)
		}
	}

	return resp, nil
}

func (s *WeatherStore) UpsertLegacySeaTemperatureFromMeasurement(measuredAt time.Time, temp float64) error {
	day := measuredAt.In(canaryLocation()).Format("2006-01-02")
	return s.StoreSeaTemperature(day, temp)
}

func WaterMeasurementAtLegacyDefault(date string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", strings.TrimSpace(date))
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 10, 0, 0, 0, canaryLocation()), nil
}

func ParseMeasuredAtOrDate(measuredAtRaw, dateRaw string) (time.Time, error) {
	measuredAtRaw = strings.TrimSpace(measuredAtRaw)
	if measuredAtRaw != "" {
		if t, err := time.Parse(time.RFC3339, measuredAtRaw); err == nil {
			return t.In(canaryLocation()), nil
		}

		layouts := []string{
			"2006-01-02 15:04:05",
			"2006-01-02 15:04",
			"2006-01-02T15:04:05",
			"2006-01-02T15:04",
		}
		for _, layout := range layouts {
			if t, err := time.ParseInLocation(layout, measuredAtRaw, canaryLocation()); err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("invalid measured_at format")
	}

	return WaterMeasurementAtLegacyDefault(dateRaw)
}

func (s *WeatherStore) GetDailyTemperatureExtremes(date string) (*float64, *time.Time, *float64, *time.Time, error) {
	query := `
		SELECT
			(SELECT temperature
			 FROM weather
			 WHERE DATE(measured_at) = ?
			 ORDER BY temperature DESC, measured_at ASC
			 LIMIT 1) AS max_temperature,
			(SELECT measured_at
			 FROM weather
			 WHERE DATE(measured_at) = ?
			 ORDER BY temperature DESC, measured_at ASC
			 LIMIT 1) AS max_measured_at,
			(SELECT temperature
			 FROM weather
			 WHERE DATE(measured_at) = ?
			 ORDER BY temperature ASC, measured_at ASC
			 LIMIT 1) AS min_temperature,
			(SELECT measured_at
			 FROM weather
			 WHERE DATE(measured_at) = ?
			 ORDER BY temperature ASC, measured_at ASC
			 LIMIT 1) AS min_measured_at
	`

	var maxTemp sql.NullFloat64
	var maxAt sql.NullTime
	var minTemp sql.NullFloat64
	var minAt sql.NullTime

	if err := s.DB.QueryRow(query, date, date, date, date).Scan(&maxTemp, &maxAt, &minTemp, &minAt); err != nil {
		return nil, nil, nil, nil, err
	}

	var maxPtr *float64
	var minPtr *float64
	var maxTime *time.Time
	var minTime *time.Time

	if maxTemp.Valid {
		v := maxTemp.Float64
		maxPtr = &v
	}
	if minTemp.Valid {
		v := minTemp.Float64
		minPtr = &v
	}
	if maxAt.Valid {
		maxTime = &maxAt.Time
	}
	if minAt.Valid {
		minTime = &minAt.Time
	}

	return maxPtr, maxTime, minPtr, minTime, nil
}

func formatTimeHM(ts time.Time) string {
	return strings.TrimPrefix(ts.Format("15:04"), "0")
}

func (s *WeatherStore) GetTideEvents(ctx context.Context, dateLocal, locationKey string) ([]models.TideEvent, error) {
	query := `
		SELECT id, date_local, location_key, event_type, event_time_local, height_m, source, confidence, fetched_at, raw_json
		FROM tide_events
		WHERE date_local = ? AND location_key = ?
		ORDER BY
			CASE source WHEN 'puertos' THEN 0 ELSE 1 END,
			event_time_local ASC
	`

	rows, err := s.DB.QueryContext(ctx, query, dateLocal, locationKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.TideEvent, 0, 4)
	for rows.Next() {
		var ev models.TideEvent
		if err := rows.Scan(
			&ev.ID,
			&ev.DateLocal,
			&ev.LocationKey,
			&ev.EventType,
			&ev.EventTimeLocal,
			&ev.HeightM,
			&ev.Source,
			&ev.Confidence,
			&ev.FetchedAt,
			&ev.RawJSON,
		); err != nil {
			return nil, err
		}
		out = append(out, ev)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *WeatherStore) HasFreshPuertosData(ctx context.Context, dateLocal, locationKey string, sinceUTC time.Time) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM tide_events
		WHERE date_local = ?
			AND location_key = ?
			AND source = 'puertos'
			AND fetched_at >= ?
	`

	var count int
	if err := s.DB.QueryRowContext(ctx, query, dateLocal, locationKey, sinceUTC).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *WeatherStore) ReplaceTideEvents(ctx context.Context, dateLocal, locationKey string, events []models.TideEvent) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err = tx.ExecContext(ctx, "DELETE FROM tide_events WHERE date_local = ? AND location_key = ?", dateLocal, locationKey); err != nil {
		return err
	}

	insertQuery := `
		INSERT INTO tide_events (
			date_local, location_key, event_type, event_time_local, height_m, source, confidence, fetched_at, raw_json
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	for _, ev := range events {
		// Store local wall-clock datetime as literal string (DATETIME has no timezone).
		// This avoids driver timezone conversion shifting the local tide time.
		eventTimeLocal := ev.EventTimeLocal.Format("2006-01-02 15:04:05")
		if _, err = tx.ExecContext(ctx, insertQuery,
			ev.DateLocal,
			ev.LocationKey,
			ev.EventType,
			eventTimeLocal,
			ev.HeightM,
			ev.Source,
			ev.Confidence,
			ev.FetchedAt.UTC(),
			ev.RawJSON,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}
