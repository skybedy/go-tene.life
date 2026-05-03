package models

type SoundTrack struct {
	Title    string
	FileName string
	URL      string
	Icon     string
}

type SoundsPageData struct {
	Tracks          []SoundTrack
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
