package models

type SoundTrack struct {
	Title        string
	FileName     string
	URL          string
	Icon         string
	Lesson       string
	Spanish      string
	Czech        string
	Segments     []SoundSegment
	SegmentsJSON string
}

type SoundSegment struct {
	Start   float64 `json:"start"`
	End     float64 `json:"end"`
	Spanish string  `json:"spanish"`
	Czech   string  `json:"czech"`
}

type SoundsPageData struct {
	Tracks          []SoundTrack
	Track250        *SoundTrack
	Track500        *SoundTrack
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
