package models

type WaterQualityLatest struct {
	Location     string `json:"location"`
	DataKind     string `json:"data_kind"`
	Source       string `json:"source"`
	Status       string `json:"status"`
	SampleDate   string `json:"sample_date"`
	UpdatedAtUTC string `json:"updated_at_utc"`
	Notes        string `json:"notes"`
	Link         string `json:"link"`
}
