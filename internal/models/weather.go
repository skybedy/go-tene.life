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
	PageTitle         string
}

type WeatherHourly struct {
	Date           string  `db:"date"`
	Hour           int     `db:"hour"`
	AvgTemperature float64 `db:"avg_temperature"`
	AvgPressure    float64 `db:"avg_pressure"`
	AvgHumidity    float64 `db:"avg_humidity"`
	SamplesCount   int     `db:"samples_count"`
}

type ChartResponse struct {
	Labels   []string                 `json:"labels"`
	Datasets map[string][]interface{} `json:"datasets"`
}

// Or more specific for the hourly chart
type HourlyChartResponse struct {
	Labels   []string `json:"labels"`
	Datasets struct {
		Temperature []float64 `json:"temperature"`
		Pressure    []float64 `json:"pressure"`
		Humidity    []float64 `json:"humidity"`
	} `json:"datasets"`
}

type WeatherDaily struct {
	Date           string   `db:"date"`
	SeaTemperature *float64 `db:"sea_temperature"`
	AvgTemperature float64  `db:"avg_temperature"`
	MinTemperature float64  `db:"min_temperature"`
	MaxTemperature float64  `db:"max_temperature"`
	AvgPressure    float64  `db:"avg_pressure"`
	MinPressure    float64  `db:"min_pressure"`
	MaxPressure    float64  `db:"max_pressure"`
	AvgHumidity    float64  `db:"avg_humidity"`
	MinHumidity    float64  `db:"min_humidity"`
	MaxHumidity    float64  `db:"max_humidity"`
	SamplesCount   int      `db:"samples_count"`
}

type WeatherWeekly struct {
	Year           int     `db:"year"`
	Week           int     `db:"week"`
	WeekStart      string  `db:"week_start"`
	WeekEnd        string  `db:"week_end"`
	AvgTemperature float64 `db:"avg_temperature"`
	MinTemperature float64 `db:"min_temperature"`
	MaxTemperature float64 `db:"max_temperature"`
	AvgPressure    float64 `db:"avg_pressure"`
	MinPressure    float64 `db:"min_pressure"`
	MaxPressure    float64 `db:"max_pressure"`
	AvgHumidity    float64 `db:"avg_humidity"`
	MinHumidity    float64 `db:"min_humidity"`
	MaxHumidity    float64 `db:"max_humidity"`
	SamplesCount   int     `db:"samples_count"`
}

type WeatherMonthly struct {
	Year           int     `db:"year"`
	Month          int     `db:"month"`
	AvgTemperature float64 `db:"avg_temperature"`
	MinTemperature float64 `db:"min_temperature"`
	MaxTemperature float64 `db:"max_temperature"`
	AvgPressure    float64 `db:"avg_pressure"`
	MinPressure    float64 `db:"min_pressure"`
	MaxPressure    float64 `db:"max_pressure"`
	AvgHumidity    float64 `db:"avg_humidity"`
	MinHumidity    float64 `db:"min_humidity"`
	MaxHumidity    float64 `db:"max_humidity"`
	SamplesCount   int     `db:"samples_count"`
}

type StatsPageData struct {
	DailyStats   []WeatherDaily
	WeeklyStats  []WeatherWeekly
	MonthlyStats []WeatherMonthly
	AnnualStats  []WeatherMonthly
	PageTitle    string
}
