package models

type WeatherData struct {
	Timestamp   int64   `json:"timestamp"`
	Temperature float64 `json:"temperature"`
	Pressure    float64 `json:"pressure"`
	Humidity    float64 `json:"humidity"`
}

type TideEventView struct {
	Time   string
	Height string
}

type LanguageOption struct {
	Code string
	Name string
}

type PageData struct {
	Weather            *WeatherData
	WebcamImageURL     string
	SeaTemperature     *float64
	SeaTemperatureDate string
	SeaTemperatureTime string
	SeaTemperatureVal  float64
	NextHighTide       string
	NextLowTide        string
	TideHighEvents     []TideEventView
	TideLowEvents      []TideEventView
	TideUpdatedAt      string
	Waves              *WavesLatest
	WaterQuality       *WaterQualityLatest
	DayMaxTemperature  *float64
	DayMinTemperature  *float64
	DayMaxTempText     string
	DayMinTempText     string
	DayMaxTime         string
	DayMinTime         string
	FormattedDate      string
	FormattedDateSmall string
	FormattedTime      string
	PageTitle          string
	Locale             string
	LocalePrefix       string
	CurrentPath        string
	CurrentSection     string
	Languages          []LanguageOption
	I18n               map[string]string
	GAEnabled          bool
	GAMeasurementID    string
}

type WeatherHourly struct {
	Date           string   `db:"date"`
	Hour           int      `db:"hour"`
	AvgTemperature *float64 `db:"avg_temperature"`
	AvgPressure    *float64 `db:"avg_pressure"`
	AvgHumidity    *float64 `db:"avg_humidity"`
	SamplesCount   *int     `db:"samples_count"`
}

type ChartResponse struct {
	Labels   []string                 `json:"labels"`
	Datasets map[string][]interface{} `json:"datasets"`
}

// Or more specific for the hourly chart
type HourlyChartResponse struct {
	Labels         []string `json:"labels"`
	SeaTemperature *float64 `json:"sea_temperature,omitempty"`
	SeaMeasuredAt  string   `json:"sea_measured_at,omitempty"`
	Datasets       struct {
		Temperature []*float64 `json:"temperature"`
		Pressure    []*float64 `json:"pressure"`
		Humidity    []*float64 `json:"humidity"`
	} `json:"datasets"`
}

type DailyChartResponse struct {
	Labels   []string `json:"labels"`
	Datasets struct {
		AvgTemperature []*float64 `json:"avg_temperature"`
		MinTemperature []*float64 `json:"min_temperature"`
		MaxTemperature []*float64 `json:"max_temperature"`
		AvgPressure    []*float64 `json:"avg_pressure"`
		AvgHumidity    []*float64 `json:"avg_humidity"`
		SeaTemperature []*float64 `json:"sea_temperature"`
	} `json:"datasets"`
}

type CurrentMonthSummaryResponse struct {
	Labels      []string `json:"labels"`
	ThroughDate string   `json:"through_date,omitempty"`
	Datasets    struct {
		AvgTemperature []*float64 `json:"avg_temperature"`
		AvgPressure    []*float64 `json:"avg_pressure"`
		AvgHumidity    []*float64 `json:"avg_humidity"`
		SeaTemperature []*float64 `json:"sea_temperature"`
	} `json:"datasets"`
}

type GenericChartResponse struct {
	Labels   []string                 `json:"labels"`
	Datasets map[string][]interface{} `json:"datasets"`
}

type WeatherDaily struct {
	Date           string   `db:"date"`
	SeaTemperature *float64 `db:"sea_temperature"`
	SeaMeasuredAt  *string  `db:"sea_measured_at"`
	AvgTemperature *float64 `db:"avg_temperature"`
	MinTemperature *float64 `db:"min_temperature"`
	MaxTemperature *float64 `db:"max_temperature"`
	AvgPressure    *float64 `db:"avg_pressure"`
	MinPressure    *float64 `db:"min_pressure"`
	MaxPressure    *float64 `db:"max_pressure"`
	AvgHumidity    *float64 `db:"avg_humidity"`
	MinHumidity    *float64 `db:"min_humidity"`
	MaxHumidity    *float64 `db:"max_humidity"`
	SamplesCount   *int     `db:"samples_count"`
}

type WeatherWeekly struct {
	Year              int      `db:"year"`
	Week              int      `db:"week"`
	WeekStart         string   `db:"week_start"`
	WeekEnd           string   `db:"week_end"`
	AvgSeaTemperature *float64 `db:"avg_sea_temperature"`
	AvgTemperature    *float64 `db:"avg_temperature"`
	MinTemperature    *float64 `db:"min_temperature"`
	MaxTemperature    *float64 `db:"max_temperature"`
	AvgPressure       *float64 `db:"avg_pressure"`
	MinPressure       *float64 `db:"min_pressure"`
	MaxPressure       *float64 `db:"max_pressure"`
	AvgHumidity       *float64 `db:"avg_humidity"`
	MinHumidity       *float64 `db:"min_humidity"`
	MaxHumidity       *float64 `db:"max_humidity"`
	SamplesCount      *int     `db:"samples_count"`
}

type WeatherMonthly struct {
	Year              int      `db:"year"`
	Month             int      `db:"month"`
	AvgSeaTemperature *float64 `db:"avg_sea_temperature"`
	AvgTemperature    *float64 `db:"avg_temperature"`
	MinTemperature    *float64 `db:"min_temperature"`
	MaxTemperature    *float64 `db:"max_temperature"`
	AvgPressure       *float64 `db:"avg_pressure"`
	MinPressure       *float64 `db:"min_pressure"`
	MaxPressure       *float64 `db:"max_pressure"`
	AvgHumidity       *float64 `db:"avg_humidity"`
	MinHumidity       *float64 `db:"min_humidity"`
	MaxHumidity       *float64 `db:"max_humidity"`
	SamplesCount      *int     `db:"samples_count"`
}

type CurrentMonthSummary struct {
	Year                  int      `json:"year"`
	Month                 int      `json:"month"`
	ThroughDate           string   `json:"through_date"`
	AvgTemperature        *float64 `json:"avg_temperature"`
	AvgPressure           *float64 `json:"avg_pressure"`
	AvgHumidity           *float64 `json:"avg_humidity"`
	AvgSeaTemperature     *float64 `json:"avg_sea_temperature"`
	WeatherSamplesCount   int      `json:"weather_samples_count"`
	SeaSamplesCount       int      `json:"sea_samples_count"`
	WeatherRangeStartUTC  string   `json:"weather_range_start_utc"`
	WeatherRangeEndUTC    string   `json:"weather_range_end_utc"`
	LocalRangeStartCanary string   `json:"local_range_start_canary"`
	LocalRangeEndCanary   string   `json:"local_range_end_canary"`
}

type MonthOption struct {
	Year  int
	Month int
}

type StatsPageData struct {
	DailyStats      []WeatherDaily
	WeeklyStats     []WeatherWeekly
	MonthlyStats    []WeatherMonthly
	AnnualStats     []WeatherMonthly
	SelectedDate    string
	AvailableMonths []MonthOption
	SelectedYear    int
	SelectedMonth   int
	PageTitle       string
	StatsSection    string
	Locale          string
	LocalePrefix    string
	CurrentPath     string
	Languages       []LanguageOption
	I18n            map[string]string
	CurrentSection  string
	GAEnabled       bool
	GAMeasurementID string
}

type HomeAPIResponse struct {
	Weather        *WeatherData        `json:"weather,omitempty"`
	SeaTemperature *float64            `json:"sea_temperature,omitempty"`
	SeaMeasuredAt  string              `json:"sea_measured_at,omitempty"`
	Waves          *WavesLatest        `json:"waves,omitempty"`
	WaterQuality   *WaterQualityLatest `json:"water_quality,omitempty"`
}

type WaterTemperatureMeasurement struct {
	ID                   uint64   `db:"id" json:"id"`
	MeasuredAt           string   `db:"measured_at" json:"measured_at"`
	Temperature          *float64 `db:"temperature" json:"temperature"`
	Source               string   `db:"source" json:"source"`
	Note                 *string  `db:"note" json:"note,omitempty"`
	LegacyWeatherDailyID *uint64  `db:"legacy_weather_daily_id" json:"legacy_weather_daily_id,omitempty"`
}

type WaterTemperatureHistoryResponse struct {
	Labels       []string   `json:"labels"`
	Temperatures []*float64 `json:"temperatures"`
	MeasuredAt   []string   `json:"measured_at"`
}
