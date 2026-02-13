package store

import (
	"database/sql"
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
