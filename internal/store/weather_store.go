package store

import (
	"database/sql"

	"github.com/skybedy/laravel-tene.life/internal/models"
)

type WeatherStore struct {
	DB *sql.DB
}

func NewWeatherStore(db *sql.DB) *WeatherStore {
	return &WeatherStore{DB: db}
}

func (s *WeatherStore) GetLatestSeaTemperature() (*float64, error) {
	var seaTemp *float64
	var temp float64

	// Try today's temperature first
	query := "SELECT sea_temperature FROM weather_daily WHERE date = CURRENT_DATE AND sea_temperature IS NOT NULL LIMIT 1"
	err := s.DB.QueryRow(query).Scan(&temp)
	if err == nil {
		seaTemp = &temp
		return seaTemp, nil
	}

	// Fallback to latest available
	query = "SELECT sea_temperature FROM weather_daily WHERE sea_temperature IS NOT NULL ORDER BY date DESC LIMIT 1"
	err = s.DB.QueryRow(query).Scan(&temp)
	if err == nil {
		seaTemp = &temp
		return seaTemp, nil
	}

	return nil, nil // Return nil if no temperature found, not strictly an error for the view
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

func (s *WeatherStore) GetWeeklyStats() ([]models.WeatherWeekly, error) {
	var results []models.WeatherWeekly
	query := `SELECT year, week, week_start, week_end, avg_temperature, min_temperature, max_temperature, 
	                 avg_pressure, min_pressure, max_pressure, avg_humidity, min_humidity, max_humidity, samples_count 
	          FROM weather_weekly ORDER BY year DESC, week DESC`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var w models.WeatherWeekly
		err := rows.Scan(&w.Year, &w.Week, &w.WeekStart, &w.WeekEnd, &w.AvgTemperature, &w.MinTemperature, &w.MaxTemperature,
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
	query := `SELECT year, month, avg_temperature, min_temperature, max_temperature, 
	                 avg_pressure, min_pressure, max_pressure, avg_humidity, min_humidity, max_humidity, samples_count 
	          FROM weather_monthly ORDER BY year DESC, month DESC LIMIT ?`
	rows, err := s.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m models.WeatherMonthly
		err := rows.Scan(&m.Year, &m.Month, &m.AvgTemperature, &m.MinTemperature, &m.MaxTemperature,
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
	query := `SELECT year, month, avg_temperature, min_temperature, max_temperature, 
	                 avg_pressure, min_pressure, max_pressure, avg_humidity, min_humidity, max_humidity, samples_count 
	          FROM weather_monthly ORDER BY year DESC, month DESC`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m models.WeatherMonthly
		err := rows.Scan(&m.Year, &m.Month, &m.AvgTemperature, &m.MinTemperature, &m.MaxTemperature,
			&m.AvgPressure, &m.MinPressure, &m.MaxPressure, &m.AvgHumidity, &m.MinHumidity, &m.MaxHumidity, &m.SamplesCount)
		if err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, nil
}
