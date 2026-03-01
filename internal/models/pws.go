package models

import "time"

type PWSStation struct {
	StationID string
	Name      string
	Lat       *float64
	Lon       *float64
}

type PWSLatestRecord struct {
	StationID    string
	TempC        *float64
	Humidity     *float64
	Lat          *float64
	Lon          *float64
	ObsTimeUTC   *time.Time
	FetchedAtUTC time.Time
	Stale        bool
	Invalid      bool
	ErrorMessage string
}

type PWSMapPoint struct {
	StationID    string   `json:"stationId"`
	Name         string   `json:"name"`
	Lat          *float64 `json:"lat,omitempty"`
	Lon          *float64 `json:"lon,omitempty"`
	TempC        *float64 `json:"temp_c,omitempty"`
	Humidity     *float64 `json:"humidity,omitempty"`
	ObsTimeUTC   string   `json:"obs_time_utc,omitempty"`
	FetchedAtUTC string   `json:"fetched_at_utc,omitempty"`
	Stale        bool     `json:"stale"`
	Invalid      bool     `json:"invalid"`
	Error        string   `json:"error,omitempty"`
}

type PWSMapPageData struct {
	PageTitle       string
	Locale          string
	LocalePrefix    string
	CurrentPath     string
	CurrentSection  string
	Languages       []LanguageOption
	I18n            map[string]string
	GAEnabled       bool
	GAMeasurementID string
}
