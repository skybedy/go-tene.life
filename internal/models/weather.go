package models

type WeatherData struct {
	Timestamp   int64   `json:"timestamp"`
	Temperature float64 `json:"temperature"`
	Pressure    float64 `json:"pressure"`
	Humidity    float64 `json:"humidity"`
}

type PageData struct {
	Weather           *WeatherData
	SeaTemperature    *float64
	SeaTemperatureVal float64
	FormattedDate     string
	FormattedTime     string
}
