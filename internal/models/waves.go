package models

// WavesLatest is the cached measured wave payload used by homepage/API.
type WavesLatest struct {
	Location      string  `json:"location"`
	StationID     int     `json:"station_id"`
	DataKind      string  `json:"data_kind"`
	Source        string  `json:"source"`
	MeasuredAtUTC string  `json:"measured_at_utc"`
	Hm0M          float64 `json:"hm0_m"`
	PeriodS       float64 `json:"period_s"`
	MeanDirDeg    float64 `json:"mean_dir_deg"`
	FetchedAtUTC  string  `json:"fetched_at_utc"`
}
