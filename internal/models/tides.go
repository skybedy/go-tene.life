package models

import "time"

// TideEvent stores one high/low event for a specific location and local day.
type TideEvent struct {
	ID             int64     `db:"id" json:"id"`
	DateLocal      string    `db:"date_local" json:"date_local"`
	LocationKey    string    `db:"location_key" json:"location_key"`
	EventType      string    `db:"event_type" json:"event_type"`
	EventTimeLocal time.Time `db:"event_time_local" json:"event_time_local"`
	HeightM        float64   `db:"height_m" json:"height_m"`
	Source         string    `db:"source" json:"source"`
	Confidence     int       `db:"confidence" json:"confidence"`
	FetchedAt      time.Time `db:"fetched_at" json:"fetched_at"`
	RawJSON        string    `db:"raw_json" json:"raw_json"`
}

type TideEventResponse struct {
	Type      string  `json:"type"`
	TimeLocal string  `json:"time_local"`
	HeightM   float64 `json:"height_m"`
}

type TideAPIResponse struct {
	DateLocal  string              `json:"date_local"`
	Location   string              `json:"location"`
	Source     string              `json:"source"`
	Events     []TideEventResponse `json:"events"`
	FetchedAt  string              `json:"fetched_at"`
	Confidence int                 `json:"confidence"`
}
