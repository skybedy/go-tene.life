package i18n

import (
	"fmt"
	"math"
	"strings"

	"github.com/skybedy/laravel-tene.life/internal/models"
)

const DefaultLocale = "cs"

var languages = []models.LanguageOption{
	{Code: "cs", Name: "Čeština"},
	{Code: "en", Name: "English"},
	{Code: "es", Name: "Español"},
	{Code: "pl", Name: "Polski"},
	{Code: "de", Name: "Deutsch"},
	{Code: "fr", Name: "Français"},
	{Code: "it", Name: "Italiano"},
	{Code: "hu", Name: "Magyar"},
}

var supportedLocales = func() map[string]struct{} {
	m := make(map[string]struct{}, len(languages))
	for _, lang := range languages {
		m[lang.Code] = struct{}{}
	}
	return m
}()

var dictionary = map[string]map[string]string{
	"cs": {
		"home":                        "Webkamera",
		"statistics":                  "Rozšířené statistiky",
		"daily_statistics":            "Denní statistiky",
		"recent_statistics":           "Posledních 10 dní",
		"tenerife_temperatures":       "Teplotní mapa Kanárských ostrovů",
		"weekly_statistics":           "Týdenní statistiky",
		"monthly_statistics":          "Měsíční statistiky",
		"annual_statistics":           "Roční statistiky",
		"daily_short":                 "Denní",
		"recent_short":                "10 dní",
		"weekly_short":                "Týdenní",
		"monthly_short":               "Měsíční",
		"annual_short":                "Roční",
		"webcam_title":                "Webkamera – Tenerife, Los Cristianos",
		"webcam_view_alt":             "Webkamera výhled",
		"location_heading":            "Umístění a směr pohledu",
		"location_description":        "Avenida Ámsterdam, severovýchod – výhled na Montaña el Mojón 250 m/nm a Roque de Ichasagua 1001, dále, úplně vpravo za stromem, na Morros del Viento 406 a při jasné obloze pak v pozadí i na Pico del Teide 3715, Pico Viejo 3135 a Alto de Guajara 2715.",
		"weather_source_heading":      "O zdroji meteorologických dat",
		"weather_source_description":  "Data o aktuální teplotě, tlaku a vlhkosti jsou odebírána z vlastní meteostanice a teplotního čidla v celodenně zastíněném místě, bez dosahu přímého slunce, takže se jedná čistě o hodnoty ve stínu.",
		"hobby_disclaimer":            "Vezměte prosím také na vědomí, že informace o počasí zde uváděné jsou pouze hobby projektem s daty získávanými amatérskou meteorologickou technikou a metodami, bez jakýchkoliv ambicí konkurovat tradičním meteorologickým zdrojům.",
		"temperature":                 "Teplota",
		"air_temperature":             "Teplota vzduchu",
		"sea_temperature":             "Teplota moře",
		"tide_high":                   "Příliv",
		"tide_low":                    "Odliv",
		"wave_height":                 "Výška vln",
		"wave_period":                 "Perioda vln",
		"wave_direction":              "Směr vln",
		"water_quality":               "Kvalita vody",
		"water_sample_date":           "Datum vzorku",
		"pressure":                    "Tlak",
		"atmospheric_pressure":        "Atmosférický tlak",
		"humidity":                    "Vlhkost",
		"relative_humidity":           "Relativní vlhkost",
		"weather_unavailable":         "Počasí nedostupné",
		"weather_data_title":          "Meteorologická data",
		"weather_data_subtitle":       "Grafy zobrazují hodinové průměry od půlnoci dnešního dne",
		"temperature_chart":           "Teplota (°C)",
		"pressure_chart":              "Atmosférický tlak (hPa)",
		"humidity_chart":              "Relativní vlhkost (%)",
		"daily_stats_subtitle":        "Přehled meteorologických dat za posledních 30 dní.",
		"weekly_stats_subtitle":       "Týdenní průměry a extrémy meteorologických dat.",
		"monthly_stats_subtitle":      "Měsíční přehledy za poslední rok.",
		"annual_stats_subtitle":       "Kompletní historie měsíčních průměrů.",
		"daily_temp_chart_title":      "Teplota za posledních 7 dní (°C)",
		"daily_pressure_chart_title":  "Tlak za posledních 7 dní (hPa)",
		"daily_humidity_chart_title":  "Vlhkost za posledních 7 dní (%)",
		"recent_stats_subtitle":       "Grafy zobrazují denní průměry za posledních 10 dní.",
		"recent_temp_chart_title":     "Teplota za posledních 10 dní (°C)",
		"recent_pressure_chart_title": "Tlak za posledních 10 dní (hPa)",
		"recent_humidity_chart_title": "Vlhkost za posledních 10 dní (%)",
		"weekly_temp_chart_title":     "Teplota (týdenní průměry) °C",
		"monthly_temp_chart_title":    "Teplota (měsíční průměry) °C",
		"go_to_detail_charts":         "Přejít na detailní grafy dne:",
		"go_to_detail_charts_hint":    "Načte hodinové grafy pro vybraný den přímo níže na stránce.",
		"detail_charts_load_error":    "Nepodařilo se načíst hodinové grafy pro vybrané datum.",
		"show":                        "Zobrazit",
		"table_overview":              "Tabulkový přehled",
		"table_overview_30":           "Tabulkový přehled (posledních 30 dní)",
		"date":                        "Datum",
		"click_to_change_date":        "Kliknutím můžete změnit datum",
		"week_year":                   "Týden / Rok",
		"month_year":                  "Měsíc / Rok",
		"period":                      "Období",
		"avg_temp":                    "Ø Teplota",
		"min_max":                     "Min / Max",
		"avg_pressure":                "Ø Tlak",
		"avg_humidity":                "Ø Vlhkost",
		"sea":                         "Moře",
		"no_data":                     "Zatím nejsou k dispozici žádná data.",
		"average":                     "Průměr:",
		"min":                         "Min",
		"max":                         "Max",
		"back":                        "Zpět",
		"webcam_big_title":            "Webkamera – Velký náhled",
		"webcam_big_alt":              "Webkamera - velký náhled",
		"site_title_suffix":           "Tenerife | Los Cristianos | Webcam",
		"pws_map_subtitle":            "Aktuální teploty z vybraných meteostanic na Kanárských ostrovech.",
		"pws_last_update":             "Poslední aktualizace",
		"pws_no_data":                 "Zatím nejsou dostupná data ze stanic.",
		"pws_observed_at":             "Naměřeno",
		"pws_fetched_at":              "Staženo",
		"pws_stale":                   "Neaktuální (starší než 30 minut)",
		"pws_invalid":                 "Podezřelá hodnota",
	},
	"en": {
		"home":                        "Webcam",
		"statistics":                  "Extended Statistics",
		"daily_statistics":            "Daily Statistics",
		"recent_statistics":           "Last 10 Days",
		"tenerife_temperatures":       "Canary Islands Temperature Map",
		"weekly_statistics":           "Weekly Statistics",
		"monthly_statistics":          "Monthly Statistics",
		"annual_statistics":           "Annual Statistics",
		"daily_short":                 "Daily",
		"recent_short":                "10 days",
		"weekly_short":                "Weekly",
		"monthly_short":               "Monthly",
		"annual_short":                "Annual",
		"webcam_title":                "Webcam – Tenerife, Los Cristianos",
		"webcam_view_alt":             "Webcam view",
		"location_heading":            "Location and View Direction",
		"location_description":        "Avenida Ámsterdam, northeast – view of Montaña el Mojón 250 m asl and Roque de Ichasagua 1001, further, far right behind the tree, towards Morros del Viento 406 and in clear weather also in the background Pico del Teide 3715, Pico Viejo 3135 and Alto de Guajara 2715.",
		"weather_source_heading":      "About Weather Data Source",
		"weather_source_description":  "Data on current temperature, pressure and humidity are collected from our own weather station and temperature sensor in a permanently shaded location, without direct sunlight, so these are pure shade values.",
		"hobby_disclaimer":            "Please also note that the weather information provided here is only a hobby project with data obtained using amateur meteorological equipment and methods, without any ambition to compete with traditional meteorological sources.",
		"temperature":                 "Temperature",
		"air_temperature":             "Air Temperature",
		"sea_temperature":             "Sea Temperature",
		"tide_high":                   "High Tide",
		"tide_low":                    "Low Tide",
		"wave_height":                 "Wave Height",
		"wave_period":                 "Wave Period",
		"wave_direction":              "Wave Direction",
		"water_quality":               "Water Quality",
		"water_sample_date":           "Sample Date",
		"pressure":                    "Pressure",
		"atmospheric_pressure":        "Atmospheric Pressure",
		"humidity":                    "Humidity",
		"relative_humidity":           "Relative Humidity",
		"weather_unavailable":         "Weather data is currently unavailable",
		"weather_data_title":          "Meteorological Data",
		"weather_data_subtitle":       "Charts show hourly averages from midnight today",
		"temperature_chart":           "Temperature (°C)",
		"pressure_chart":              "Atmospheric Pressure (hPa)",
		"humidity_chart":              "Relative Humidity (%)",
		"daily_stats_subtitle":        "Overview of meteorological data for the last 30 days.",
		"weekly_stats_subtitle":       "Weekly averages and extremes of meteorological data.",
		"monthly_stats_subtitle":      "Monthly overviews for the past year.",
		"annual_stats_subtitle":       "Complete history of monthly averages.",
		"daily_temp_chart_title":      "Temperature over the last 7 days (°C)",
		"daily_pressure_chart_title":  "Pressure over the last 7 days (hPa)",
		"daily_humidity_chart_title":  "Humidity over the last 7 days (%)",
		"recent_stats_subtitle":       "Charts show daily averages over the last 10 days.",
		"recent_temp_chart_title":     "Temperature over the last 10 days (°C)",
		"recent_pressure_chart_title": "Pressure over the last 10 days (hPa)",
		"recent_humidity_chart_title": "Humidity over the last 10 days (%)",
		"weekly_temp_chart_title":     "Temperature (weekly averages) °C",
		"monthly_temp_chart_title":    "Temperature (monthly averages) °C",
		"go_to_detail_charts":         "Go to detailed charts of day:",
		"go_to_detail_charts_hint":    "Loads hourly charts for the selected day directly below on this page.",
		"detail_charts_load_error":    "Failed to load hourly charts for the selected date.",
		"show":                        "Show",
		"table_overview":              "Table Overview",
		"table_overview_30":           "Table Overview (last 30 days)",
		"date":                        "Date",
		"click_to_change_date":        "Click to change the date",
		"week_year":                   "Week / Year",
		"month_year":                  "Month / Year",
		"period":                      "Period",
		"avg_temp":                    "Avg. Temp",
		"min_max":                     "Min / Max",
		"avg_pressure":                "Avg. Pressure",
		"avg_humidity":                "Avg. Humidity",
		"sea":                         "Sea",
		"no_data":                     "No data available yet.",
		"average":                     "Average:",
		"min":                         "Min",
		"max":                         "Max",
		"back":                        "Back",
		"webcam_big_title":            "Webcam – Large View",
		"webcam_big_alt":              "Webcam - large view",
		"site_title_suffix":           "Tenerife | Los Cristianos | Webcam",
		"pws_map_subtitle":            "Current temperatures from selected weather stations across the Canary Islands.",
		"pws_last_update":             "Last update",
		"pws_no_data":                 "No station data available yet.",
		"pws_observed_at":             "Observed",
		"pws_fetched_at":              "Fetched",
		"pws_stale":                   "Stale (older than 30 minutes)",
		"pws_invalid":                 "Suspicious value",
	},
}

func SupportedLanguages() []models.LanguageOption {
	out := make([]models.LanguageOption, len(languages))
	copy(out, languages)
	return out
}

func IsSupportedLocale(locale string) bool {
	_, ok := supportedLocales[locale]
	return ok
}

func NormalizeLocale(locale string) string {
	if IsSupportedLocale(locale) {
		return locale
	}
	return DefaultLocale
}

func Messages(locale string) map[string]string {
	loc := NormalizeLocale(locale)
	fallback := dictionary[DefaultLocale]
	chosen := dictionary[loc]

	merged := make(map[string]string, len(fallback))
	for k, v := range fallback {
		merged[k] = v
	}
	for k, v := range chosen {
		merged[k] = v
	}

	for k, v := range extraLocaleMessages[loc] {
		merged[k] = v
	}
	for k, v := range waterQualityLocaleMessages[loc] {
		merged[k] = v
	}
	for k, v := range recentLocaleMessages[loc] {
		merged[k] = v
	}
	for k, v := range waveSourceLocaleMessages[loc] {
		merged[k] = v
	}

	return merged
}

func T(locale, key string) string {
	msgs := Messages(locale)
	if val, ok := msgs[key]; ok {
		return val
	}
	return key
}

func WaterQualityStatusLabel(locale, raw string) string {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "EXCELENTE":
		return T(locale, "water_quality_status_excellent")
	case "BUENA":
		return T(locale, "water_quality_status_good")
	case "SUFICIENTE":
		return T(locale, "water_quality_status_sufficient")
	case "INSUFICIENTE":
		return T(locale, "water_quality_status_poor")
	default:
		return strings.TrimSpace(raw)
	}
}

func WaterQualityTooltip(locale string) string {
	return fmt.Sprintf(
		"%s: %s, %s, %s, %s",
		T(locale, "water_quality_scale"),
		T(locale, "water_quality_status_excellent"),
		T(locale, "water_quality_status_good"),
		T(locale, "water_quality_status_sufficient"),
		T(locale, "water_quality_status_poor"),
	)
}

func WaveDirectionLabel(locale string, deg float64) string {
	loc := NormalizeLocale(locale)
	dirsByLocale := map[string][]string{
		"cs": {"sever", "severovýchod", "východ", "jihovýchod", "jih", "jihozápad", "západ", "severozápad"},
		"en": {"north", "northeast", "east", "southeast", "south", "southwest", "west", "northwest"},
		"es": {"norte", "noreste", "este", "sureste", "sur", "suroeste", "oeste", "noroeste"},
		"pl": {"północ", "północny wschód", "wschód", "południowy wschód", "południe", "południowy zachód", "zachód", "północny zachód"},
		"de": {"Norden", "Nordosten", "Osten", "Südosten", "Süden", "Südwesten", "Westen", "Nordwesten"},
		"fr": {"nord", "nord-est", "est", "sud-est", "sud", "sud-ouest", "ouest", "nord-ouest"},
		"it": {"nord", "nord-est", "est", "sud-est", "sud", "sud-ovest", "ovest", "nord-ovest"},
		"hu": {"észak", "északkelet", "kelet", "délkelet", "dél", "délnyugat", "nyugat", "északnyugat"},
	}
	dirs, ok := dirsByLocale[loc]
	if !ok || len(dirs) != 8 {
		dirs = dirsByLocale["en"]
	}

	normalized := math.Mod(deg, 360.0)
	if normalized < 0 {
		normalized += 360.0
	}
	idx := int(math.Floor((normalized+22.5)/45.0)) % 8
	return dirs[idx]
}

func LocalePrefix(locale string) string {
	loc := NormalizeLocale(locale)
	if loc == DefaultLocale {
		return ""
	}
	return "/" + loc
}

func LocaleURL(locale, path string) string {
	cleanPath := path
	if cleanPath == "" {
		cleanPath = "/"
	}
	if !strings.HasPrefix(cleanPath, "/") {
		cleanPath = "/" + cleanPath
	}
	if cleanPath != "/" {
		cleanPath = strings.TrimRight(cleanPath, "/")
	}

	prefix := LocalePrefix(locale)
	if cleanPath == "/" {
		if prefix == "" {
			return "/"
		}
		return prefix
	}

	return prefix + cleanPath
}

func StripLocalePrefix(path string) string {
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return "/"
	}
	if IsSupportedLocale(parts[0]) {
		if len(parts) == 1 {
			return "/"
		}
		return "/" + strings.Join(parts[1:], "/")
	}
	return path
}

func MonthName(locale string, month int) string {
	if month < 1 || month > 12 {
		return fmt.Sprintf("%d", month)
	}
	key := fmt.Sprintf("month_%d", month)
	return T(locale, key)
}

func LanguageFlag(locale string) string {
	switch NormalizeLocale(locale) {
	case "cs":
		return "🇨🇿"
	case "en":
		return "🇬🇧"
	case "es":
		return "🇪🇸"
	case "pl":
		return "🇵🇱"
	case "de":
		return "🇩🇪"
	case "fr":
		return "🇫🇷"
	case "it":
		return "🇮🇹"
	case "hu":
		return "🇭🇺"
	default:
		return "🌐"
	}
}

var extraLocaleMessages = map[string]map[string]string{
	"es": {
		"home": "Webcam", "statistics": "Estadísticas Extendidas", "daily_statistics": "Estadísticas Diarias", "weekly_statistics": "Estadísticas Semanales", "monthly_statistics": "Estadísticas Mensuales", "annual_statistics": "Estadísticas Anuales", "daily_short": "Diarias", "weekly_short": "Semanales", "monthly_short": "Mensuales", "annual_short": "Anuales", "webcam_title": "Cámara Web – Tenerife, Los Cristianos", "webcam_view_alt": "Vista de la webcam", "location_heading": "Ubicación y Dirección de Vista", "location_description": "Avenida Ámsterdam, noreste – vista de Montaña el Mojón 250 m snm y Roque de Ichasagua 1001, más adelante, a la derecha detrás del árbol, hacia Morros del Viento 406 y con cielo despejado también al fondo Pico del Teide 3715, Pico Viejo 3135 y Alto de Guajara 2715.", "weather_source_heading": "Sobre la Fuente de Datos Meteorológicos", "weather_source_description": "Los datos de temperatura, presión y humedad actuales se recopilan de nuestra propia estación meteorológica y sensor de temperatura en un lugar permanentemente sombreado, sin luz solar directa, por lo que son valores puros a la sombra.", "hobby_disclaimer": "Tenga en cuenta también que la información meteorológica aquí proporcionada es solo un proyecto de hobby con datos obtenidos mediante equipos y métodos meteorológicos amateur, sin ninguna ambición de competir con fuentes meteorológicas tradicionales.", "temperature": "Temperatura", "air_temperature": "Temperatura del Aire", "sea_temperature": "Temperatura del Mar", "tide_high": "Marea alta", "tide_low": "Marea baja", "wave_height": "Altura de ola", "wave_period": "Período de ola", "wave_direction": "Dirección de ola", "pressure": "Presión", "atmospheric_pressure": "Presión Atmosférica", "humidity": "Humedad", "relative_humidity": "Humedad Relativa", "weather_unavailable": "Los datos meteorológicos no están disponibles actualmente", "weather_data_title": "Datos Meteorológicos", "weather_data_subtitle": "Los gráficos muestran promedios por hora desde la medianoche de hoy", "temperature_chart": "Temperatura (°C)", "pressure_chart": "Presión Atmosférica (hPa)", "humidity_chart": "Humedad Relativa (%)", "daily_stats_subtitle": "Resumen de los datos meteorológicos de los últimos 30 días.", "weekly_stats_subtitle": "Promedios semanales y extremos de datos meteorológicos.", "monthly_stats_subtitle": "Resumen mensual del último año.", "annual_stats_subtitle": "Historial completo de promedios mensuales.", "daily_temp_chart_title": "Temperatura de los últimos 7 días (°C)", "daily_pressure_chart_title": "Presión de los últimos 7 días (hPa)", "daily_humidity_chart_title": "Humedad de los últimos 7 días (%)", "weekly_temp_chart_title": "Temperatura (promedios semanales) °C", "monthly_temp_chart_title": "Temperatura (promedios mensuales) °C", "go_to_detail_charts": "Ir a gráficos detallados del día:", "show": "Mostrar", "table_overview": "Resumen en Tabla", "table_overview_30": "Resumen en Tabla (últimos 30 días)", "date": "Fecha", "week_year": "Semana / Año", "month_year": "Mes / Año", "period": "Período", "avg_temp": "Temp. Prom.", "min_max": "Mín / Máx", "avg_pressure": "Presión Prom.", "avg_humidity": "Humedad Prom.", "sea": "Mar", "no_data": "No hay datos disponibles aún.", "average": "Promedio:", "min": "Mín", "max": "Máx", "back": "Volver", "webcam_big_title": "Webcam – Vista Grande", "webcam_big_alt": "Webcam - vista grande", "site_title_suffix": "Tenerife | Los Cristianos | Webcam", "tenerife_temperatures": "Mapa de temperaturas de las Islas Canarias", "pws_map_subtitle": "Temperaturas actuales de estaciones meteorológicas seleccionadas en las Islas Canarias.", "pws_last_update": "Última actualización", "pws_no_data": "Aún no hay datos disponibles de estaciones.", "pws_observed_at": "Medido", "pws_fetched_at": "Descargado", "pws_stale": "Desactualizado (más de 30 minutos)", "pws_invalid": "Valor sospechoso",
	},
	"pl": {
		"home": "Kamera", "statistics": "Rozszerzone Statystyki", "daily_statistics": "Statystyki Dzienne", "weekly_statistics": "Statystyki Tygodniowe", "monthly_statistics": "Statystyki Miesięczne", "annual_statistics": "Statystyki Roczne", "daily_short": "Dzienne", "weekly_short": "Tygodniowe", "monthly_short": "Miesięczne", "annual_short": "Roczne", "webcam_title": "Kamera internetowa – Teneryfa, Los Cristianos", "webcam_view_alt": "Widok z kamery", "location_heading": "Lokalizacja i Kierunek Widoku", "location_description": "Avenida Ámsterdam, północny wschód – widok na Montaña el Mojón 250 m npm i Roque de Ichasagua 1001, dalej, po prawej za drzewem, na Morros del Viento 406 i przy czystym niebie także w tle Pico del Teide 3715, Pico Viejo 3135 i Alto de Guajara 2715.", "weather_source_heading": "O Źródle Danych Pogodowych", "weather_source_description": "Dane o aktualnej temperaturze, ciśnieniu i wilgotności są zbierane z naszej własnej stacji meteorologicznej i czujnika temperatury w stale zacienionym miejscu, bez bezpośredniego działania słońca, więc są to czyste wartości w cieniu.", "hobby_disclaimer": "Proszę również zauważyć, że informacje pogodowe podane tutaj to tylko projekt hobbystyczny z danymi uzyskanymi przy użyciu amatorskiego sprzętu meteorologicznego i metod, bez żadnych ambicji konkurowania z tradycyjnymi źródłami meteorologicznymi.", "temperature": "Temperatura", "air_temperature": "Temperatura Powietrza", "sea_temperature": "Temperatura Morza", "tide_high": "Przypływ", "tide_low": "Odpływ", "wave_height": "Wysokość fal", "wave_period": "Okres fali", "wave_direction": "Kierunek fali", "pressure": "Ciśnienie", "atmospheric_pressure": "Ciśnienie Atmosferyczne", "humidity": "Wilgotność", "relative_humidity": "Wilgotność Względna", "weather_unavailable": "Dane pogodowe są obecnie niedostępne", "weather_data_title": "Dane Meteorologiczne", "weather_data_subtitle": "Wykresy pokazują średnie godzinowe od północy dzisiejszego dnia", "temperature_chart": "Temperatura (°C)", "pressure_chart": "Ciśnienie Atmosferyczne (hPa)", "humidity_chart": "Wilgotność Względna (%)", "daily_stats_subtitle": "Przegląd danych meteorologicznych z ostatnich 30 dni.", "weekly_stats_subtitle": "Tygodniowe średnie i ekstrema danych meteorologicznych.", "monthly_stats_subtitle": "Miesięczne podsumowania za ostatni rok.", "annual_stats_subtitle": "Pełna historia średnich miesięcznych.", "daily_temp_chart_title": "Temperatura z ostatnich 7 dni (°C)", "daily_pressure_chart_title": "Ciśnienie z ostatnich 7 dni (hPa)", "daily_humidity_chart_title": "Wilgotność z ostatnich 7 dni (%)", "weekly_temp_chart_title": "Temperatura (średnie tygodniowe) °C", "monthly_temp_chart_title": "Temperatura (średnie miesięczne) °C", "go_to_detail_charts": "Przejdź do szczegółowych wykresów dnia:", "show": "Pokaż", "table_overview": "Przegląd Tabelaryczny", "table_overview_30": "Przegląd Tabelaryczny (ostatnie 30 dni)", "date": "Data", "week_year": "Tydzień / Rok", "month_year": "Miesiąc / Rok", "period": "Okres", "avg_temp": "Śr. Temp.", "min_max": "Min / Max", "avg_pressure": "Śr. Ciśnienie", "avg_humidity": "Śr. Wilgotność", "sea": "Morze", "no_data": "Brak dostępnych danych.", "average": "Średnia:", "min": "Min", "max": "Max", "back": "Wstecz", "webcam_big_title": "Webcam – Duży Podgląd", "webcam_big_alt": "Webcam - duży podgląd", "site_title_suffix": "Tenerife | Los Cristianos | Webcam", "tenerife_temperatures": "Mapa temperatur Wysp Kanaryjskich", "pws_map_subtitle": "Aktualne temperatury z wybranych stacji pogodowych na Wyspach Kanaryjskich.", "pws_last_update": "Ostatnia aktualizacja", "pws_no_data": "Brak dostępnych danych ze stacji.", "pws_observed_at": "Pomiar", "pws_fetched_at": "Pobrano", "pws_stale": "Nieaktualne (starsze niż 30 minut)", "pws_invalid": "Podejrzana wartość",
	},
	"de": {
		"home": "Webcam", "statistics": "Erweiterte Statistiken", "daily_statistics": "Tägliche Statistiken", "weekly_statistics": "Wöchentliche Statistiken", "monthly_statistics": "Monatliche Statistiken", "annual_statistics": "Jährliche Statistiken", "daily_short": "Täglich", "weekly_short": "Wöchentlich", "monthly_short": "Monatlich", "annual_short": "Jährlich", "webcam_title": "Webcam – Teneriffa, Los Cristianos", "webcam_view_alt": "Webcam Ansicht", "location_heading": "Standort und Blickrichtung", "location_description": "Avenida Ámsterdam, Nordosten – Blick auf Montaña el Mojón 250 m ü.M. und Roque de Ichasagua 1001, weiter, ganz rechts hinter dem Baum, auf Morros del Viento 406 und bei klarem Wetter auch im Hintergrund auf Pico del Teide 3715, Pico Viejo 3135 und Alto de Guajara 2715.", "weather_source_heading": "Über die Wetterdatenquelle", "weather_source_description": "Daten zu aktueller Temperatur, Luftdruck und Luftfeuchtigkeit werden von unserer eigenen Wetterstation und einem Temperatursensor an einem dauerhaft beschatteten Ort ohne direkte Sonneneinstrahlung gesammelt, es handelt sich also um reine Schattenwerte.", "hobby_disclaimer": "Bitte beachten Sie auch, dass die hier bereitgestellten Wetterinformationen nur ein Hobbyprojekt mit Daten sind, die mithilfe von Amateur-Wettergeräten und -methoden gewonnen wurden, ohne jegliche Ambition, mit traditionellen meteorologischen Quellen zu konkurrieren.", "temperature": "Temperatur", "air_temperature": "Lufttemperatur", "sea_temperature": "Meerestemperatur", "tide_high": "Flut", "tide_low": "Ebbe", "wave_height": "Wellenhöhe", "wave_period": "Wellenperiode", "wave_direction": "Wellenrichtung", "pressure": "Druck", "atmospheric_pressure": "Atmosphärischer Druck", "humidity": "Luftfeuchtigkeit", "relative_humidity": "Relative Luftfeuchtigkeit", "weather_unavailable": "Wetterdaten sind derzeit nicht verfügbar", "weather_data_title": "Meteorologische Daten", "weather_data_subtitle": "Die Diagramme zeigen stündliche Durchschnittswerte ab Mitternacht heute", "temperature_chart": "Temperatur (°C)", "pressure_chart": "Atmosphärischer Druck (hPa)", "humidity_chart": "Relative Luftfeuchtigkeit (%)", "daily_stats_subtitle": "Übersicht der meteorologischen Daten der letzten 30 Tage.", "weekly_stats_subtitle": "Wöchentliche Durchschnittswerte und Extreme meteorologischer Daten.", "monthly_stats_subtitle": "Monatliche Übersicht für das letzte Jahr.", "annual_stats_subtitle": "Vollständige Historie monatlicher Durchschnittswerte.", "daily_temp_chart_title": "Temperatur der letzten 7 Tage (°C)", "daily_pressure_chart_title": "Druck der letzten 7 Tage (hPa)", "daily_humidity_chart_title": "Luftfeuchtigkeit der letzten 7 Tage (%)", "weekly_temp_chart_title": "Temperatur (wöchentliche Durchschnittswerte) °C", "monthly_temp_chart_title": "Temperatur (monatliche Durchschnittswerte) °C", "go_to_detail_charts": "Zu den Detaildiagrammen des Tages:", "show": "Anzeigen", "table_overview": "Tabellenübersicht", "table_overview_30": "Tabellenübersicht (letzte 30 Tage)", "date": "Datum", "week_year": "Woche / Jahr", "month_year": "Monat / Jahr", "period": "Zeitraum", "avg_temp": "Ø Temp.", "min_max": "Min / Max", "avg_pressure": "Ø Druck", "avg_humidity": "Ø Feuchte", "sea": "Meer", "no_data": "Noch keine Daten verfügbar.", "average": "Durchschnitt:", "min": "Min", "max": "Max", "back": "Zurück", "webcam_big_title": "Webcam – Große Ansicht", "webcam_big_alt": "Webcam - große Ansicht", "site_title_suffix": "Tenerife | Los Cristianos | Webcam", "tenerife_temperatures": "Temperaturkarte der Kanarischen Inseln", "pws_map_subtitle": "Aktuelle Temperaturen ausgewählter Wetterstationen auf den Kanarischen Inseln.", "pws_last_update": "Letzte Aktualisierung", "pws_no_data": "Noch keine Stationsdaten verfügbar.", "pws_observed_at": "Gemessen", "pws_fetched_at": "Abgerufen", "pws_stale": "Veraltet (älter als 30 Minuten)", "pws_invalid": "Verdächtiger Wert",
	},
	"fr": {
		"home": "Webcam", "statistics": "Statistiques Étendues", "daily_statistics": "Statistiques Quotidiennes", "weekly_statistics": "Statistiques Hebdomadaires", "monthly_statistics": "Statistiques Mensuelles", "annual_statistics": "Statistiques Annuelles", "daily_short": "Quotidien", "weekly_short": "Hebdo", "monthly_short": "Mensuel", "annual_short": "Annuel", "webcam_title": "Webcam – Tenerife, Los Cristianos", "webcam_view_alt": "Vue webcam", "location_heading": "Emplacement et Direction de Vue", "location_description": "Avenida Ámsterdam, nord-est – vue sur Montaña el Mojón 250 m d'altitude et Roque de Ichasagua 1001, plus loin, tout à droite derrière l'arbre, vers Morros del Viento 406 et par temps clair également en arrière-plan Pico del Teide 3715, Pico Viejo 3135 et Alto de Guajara 2715.", "weather_source_heading": "À Propos de la Source des Données Météorologiques", "weather_source_description": "Les données sur la température, la pression et l'humidité actuelles sont collectées à partir de notre propre station météorologique et capteur de température dans un endroit en permanence ombragé, sans lumière directe du soleil, il s'agit donc de valeurs pures à l'ombre.", "hobby_disclaimer": "Veuillez également noter que les informations météorologiques fournies ici ne sont qu'un projet de loisir avec des données obtenues à l'aide d'équipements et de méthodes météorologiques amateurs, sans aucune ambition de concurrencer les sources météorologiques traditionnelles.", "temperature": "Température", "air_temperature": "Température de l'Air", "sea_temperature": "Température de la Mer", "tide_high": "Marée haute", "tide_low": "Marée basse", "wave_height": "Hauteur des vagues", "wave_period": "Période des vagues", "wave_direction": "Direction des vagues", "pressure": "Pression", "atmospheric_pressure": "Pression Atmosphérique", "humidity": "Humidité", "relative_humidity": "Humidité Relative", "weather_unavailable": "Les données météorologiques sont actuellement indisponibles", "weather_data_title": "Données Météorologiques", "weather_data_subtitle": "Les graphiques montrent les moyennes horaires depuis minuit aujourd'hui", "temperature_chart": "Température (°C)", "pressure_chart": "Pression Atmosphérique (hPa)", "humidity_chart": "Humidité Relative (%)", "daily_stats_subtitle": "Aperçu des données météorologiques des 30 derniers jours.", "weekly_stats_subtitle": "Moyennes hebdomadaires et extrêmes des données météorologiques.", "monthly_stats_subtitle": "Aperçus mensuels pour l'année écoulée.", "annual_stats_subtitle": "Historique complet des moyennes mensuelles.", "daily_temp_chart_title": "Température des 7 derniers jours (°C)", "daily_pressure_chart_title": "Pression des 7 derniers jours (hPa)", "daily_humidity_chart_title": "Humidité des 7 derniers jours (%)", "weekly_temp_chart_title": "Température (moyennes hebdomadaires) °C", "monthly_temp_chart_title": "Température (moyennes mensuelles) °C", "go_to_detail_charts": "Voir les graphiques détaillés du jour :", "show": "Afficher", "table_overview": "Aperçu du Tableau", "table_overview_30": "Aperçu du Tableau (30 derniers jours)", "date": "Date", "week_year": "Semaine / Année", "month_year": "Mois / Année", "period": "Période", "avg_temp": "Temp. Moy.", "min_max": "Min / Max", "avg_pressure": "Press. Moy.", "avg_humidity": "Hum. Moy.", "sea": "Mer", "no_data": "Aucune donnée disponible pour le moment.", "average": "Moyenne:", "min": "Min", "max": "Max", "back": "Retour", "webcam_big_title": "Webcam – Grande Vue", "webcam_big_alt": "Webcam - grande vue", "site_title_suffix": "Tenerife | Los Cristianos | Webcam", "tenerife_temperatures": "Carte des températures des îles Canaries", "pws_map_subtitle": "Températures actuelles des stations météo sélectionnées aux îles Canaries.", "pws_last_update": "Dernière mise à jour", "pws_no_data": "Aucune donnée de station disponible pour l'instant.", "pws_observed_at": "Mesuré", "pws_fetched_at": "Téléchargé", "pws_stale": "Périmé (plus de 30 minutes)", "pws_invalid": "Valeur suspecte",
	},
	"it": {
		"home": "Webcam", "statistics": "Statistiche Estese", "daily_statistics": "Statistiche Giornaliere", "weekly_statistics": "Statistiche Settimanali", "monthly_statistics": "Statistiche Mensili", "annual_statistics": "Statistiche Annuali", "daily_short": "Giornaliere", "weekly_short": "Settimanali", "monthly_short": "Mensili", "annual_short": "Annuali", "webcam_title": "Webcam – Tenerife, Los Cristianos", "webcam_view_alt": "Vista webcam", "location_heading": "Posizione e Direzione della Vista", "location_description": "Avenida Ámsterdam, nord-est – vista su Montaña el Mojón 250 m slm e Roque de Ichasagua 1001, più avanti, a destra dietro l'albero, verso Morros del Viento 406 e con cielo sereno anche sullo sfondo Pico del Teide 3715, Pico Viejo 3135 e Alto de Guajara 2715.", "weather_source_heading": "Sulla Fonte dei Dati Meteorologici", "weather_source_description": "I dati su temperatura, pressione e umidità attuali sono raccolti dalla nostra stazione meteorologica e sensore di temperatura in un luogo permanentemente ombreggiato, senza luce solare diretta, quindi sono valori puri all'ombra.", "hobby_disclaimer": "Si prega inoltre di notare che le informazioni meteorologiche qui fornite sono solo un progetto hobbistico con dati ottenuti utilizzando attrezzature e metodi meteorologici amatoriali, senza alcuna ambizione di competere con fonti meteorologiche tradizionali.", "temperature": "Temperatura", "air_temperature": "Temperatura dell'Aria", "sea_temperature": "Temperatura del Mare", "tide_high": "Alta marea", "tide_low": "Bassa marea", "wave_height": "Altezza onde", "wave_period": "Periodo onde", "wave_direction": "Direzione onde", "pressure": "Pressione", "atmospheric_pressure": "Pressione Atmosferica", "humidity": "Umidità", "relative_humidity": "Umidità Relativa", "weather_unavailable": "I dati meteorologici non sono attualmente disponibili", "weather_data_title": "Dati Meteorologici", "weather_data_subtitle": "I grafici mostrano le medie orarie dalla mezzanotte di oggi", "temperature_chart": "Temperatura (°C)", "pressure_chart": "Pressione Atmosferica (hPa)", "humidity_chart": "Umidità Relativa (%)", "daily_stats_subtitle": "Panoramica dei dati meteorologici degli ultimi 30 giorni.", "weekly_stats_subtitle": "Medie settimanali ed estremi dei dati meteorologici.", "monthly_stats_subtitle": "Panoramiche mensili dell'ultimo anno.", "annual_stats_subtitle": "Cronologia completa delle medie mensili.", "daily_temp_chart_title": "Temperatura degli ultimi 7 giorni (°C)", "daily_pressure_chart_title": "Pressione degli ultimi 7 giorni (hPa)", "daily_humidity_chart_title": "Umidità degli ultimi 7 giorni (%)", "weekly_temp_chart_title": "Temperatura (medie settimanali) °C", "monthly_temp_chart_title": "Temperatura (medie mensili) °C", "go_to_detail_charts": "Vai ai grafici dettagliati del giorno:", "show": "Mostra", "table_overview": "Panoramica Tabellare", "table_overview_30": "Panoramica Tabellare (ultimi 30 giorni)", "date": "Data", "week_year": "Settimana / Anno", "month_year": "Mese / Anno", "period": "Periodo", "avg_temp": "Temp. Media", "min_max": "Min / Max", "avg_pressure": "Press. Media", "avg_humidity": "Umid. Media", "sea": "Mare", "no_data": "Nessun dato ancora disponibile.", "average": "Media:", "min": "Min", "max": "Max", "back": "Indietro", "webcam_big_title": "Webcam – Vista Grande", "webcam_big_alt": "Webcam - vista grande", "site_title_suffix": "Tenerife | Los Cristianos | Webcam", "tenerife_temperatures": "Mappa temperature delle Isole Canarie", "pws_map_subtitle": "Temperature attuali dalle stazioni meteo selezionate nelle Isole Canarie.", "pws_last_update": "Ultimo aggiornamento", "pws_no_data": "Dati delle stazioni non ancora disponibili.", "pws_observed_at": "Misurato", "pws_fetched_at": "Scaricato", "pws_stale": "Non aggiornato (più di 30 minuti)", "pws_invalid": "Valore sospetto",
	},
	"hu": {
		"home": "Webkamera", "statistics": "Bővített Statisztikák", "daily_statistics": "Napi Statisztikák", "weekly_statistics": "Heti Statisztikák", "monthly_statistics": "Havi Statisztikák", "annual_statistics": "Éves Statisztikák", "daily_short": "Napi", "weekly_short": "Heti", "monthly_short": "Havi", "annual_short": "Éves", "webcam_title": "Webkamera – Tenerife, Los Cristianos", "webcam_view_alt": "Webkamera nézet", "location_heading": "Helyszín és Nézeti Irány", "location_description": "Avenida Ámsterdam, északkelet – kilátás a Montaña el Mojón 250 m tszf és Roque de Ichasagua 1001, tovább, jobbra a fa mögött, a Morros del Viento 406 felé és tiszta időben a háttérben a Pico del Teide 3715, Pico Viejo 3135 és Alto de Guajara 2715.", "weather_source_heading": "Az Időjárási Adatok Forrásáról", "weather_source_description": "Az aktuális hőmérsékletről, légnyomásról és páratartalomról szóló adatok saját meteorológiai állomásunkról és hőmérséklet-érzékelőnkről származnak, amely állandóan árnyékolt helyen van, közvetlen napfény nélkül, ezért tiszta árnyékértékek.", "hobby_disclaimer": "Kérjük, vegye figyelembe azt is, hogy az itt megadott időjárási információk csak egy hobbi projekt, amely amatőr meteorológiai eszközökkel és módszerekkel gyűjtött adatokkal dolgozik, anélkül, hogy bármilyen ambíciója lenne a hagyományos meteorológiai forrásokkal való versenyzésre.", "temperature": "Hőmérséklet", "air_temperature": "Levegő Hőmérséklet", "sea_temperature": "Tenger Hőmérséklet", "tide_high": "Dagály", "tide_low": "Apály", "wave_height": "Hullámmagasság", "wave_period": "Hullámperiódus", "wave_direction": "Hullámirány", "pressure": "Nyomás", "atmospheric_pressure": "Légköri Nyomás", "humidity": "Páratartalom", "relative_humidity": "Relatív Páratartalom", "weather_unavailable": "Az időjárási adatok jelenleg nem érhetők el", "weather_data_title": "Meteorológiai Adatok", "weather_data_subtitle": "A diagramok a mai nap éjféltől számított óránkénti átlagokat mutatják", "temperature_chart": "Hőmérséklet (°C)", "pressure_chart": "Légköri Nyomás (hPa)", "humidity_chart": "Relatív Páratartalom (%)", "daily_stats_subtitle": "Meteorológiai adatok áttekintése az elmúlt 30 napról.", "weekly_stats_subtitle": "Heti átlagok és szélsőértékek a meteorológiai adatokban.", "monthly_stats_subtitle": "Havi áttekintés az elmúlt évről.", "annual_stats_subtitle": "A havi átlagok teljes története.", "daily_temp_chart_title": "Hőmérséklet az elmúlt 7 napban (°C)", "daily_pressure_chart_title": "Nyomás az elmúlt 7 napban (hPa)", "daily_humidity_chart_title": "Páratartalom az elmúlt 7 napban (%)", "weekly_temp_chart_title": "Hőmérséklet (heti átlagok) °C", "monthly_temp_chart_title": "Hőmérséklet (havi átlagok) °C", "go_to_detail_charts": "Ugrás a nap részletes grafikonjaihoz:", "show": "Mutat", "table_overview": "Táblázatos Áttekintés", "table_overview_30": "Táblázatos Áttekintés (utolsó 30 nap)", "date": "Dátum", "week_year": "Hét / Év", "month_year": "Hónap / Év", "period": "Időszak", "avg_temp": "Átl. Hőm.", "min_max": "Min / Max", "avg_pressure": "Átl. Nyomás", "avg_humidity": "Átl. Pára", "sea": "Tenger", "no_data": "Még nincsenek elérhető adatok.", "average": "Átlag:", "min": "Min", "max": "Max", "back": "Vissza", "webcam_big_title": "Webkamera – Nagy Nézet", "webcam_big_alt": "Webkamera - nagy nézet", "site_title_suffix": "Tenerife | Los Cristianos | Webcam", "tenerife_temperatures": "A Kanári-szigetek hőtérképe", "pws_map_subtitle": "Aktuális hőmérsékletek kiválasztott kanári-szigeteki időjárásállomásokról.", "pws_last_update": "Utolsó frissítés", "pws_no_data": "Még nincs elérhető állomásadat.", "pws_observed_at": "Mérve", "pws_fetched_at": "Letöltve", "pws_stale": "Elavult (30 percnél régebbi)", "pws_invalid": "Gyanús érték",
	},
}

var waterQualityLocaleMessages = map[string]map[string]string{
	"cs": {
		"water_quality":                   "Kvalita vody",
		"water_sample_date":               "Datum vzorku",
		"water_quality_scale":             "Stupnice kvality",
		"water_quality_status_excellent":  "Výborná",
		"water_quality_status_good":       "Dobrá",
		"water_quality_status_sufficient": "Dostatečná",
		"water_quality_status_poor":       "Nedostatečná",
	},
	"en": {
		"water_quality":                   "Water Quality",
		"water_sample_date":               "Sample Date",
		"water_quality_scale":             "Quality scale",
		"water_quality_status_excellent":  "Excellent",
		"water_quality_status_good":       "Good",
		"water_quality_status_sufficient": "Sufficient",
		"water_quality_status_poor":       "Poor",
	},
	"es": {
		"water_quality":                   "Calidad del agua",
		"water_sample_date":               "Fecha de muestra",
		"water_quality_scale":             "Escala de calidad",
		"water_quality_status_excellent":  "Excelente",
		"water_quality_status_good":       "Buena",
		"water_quality_status_sufficient": "Suficiente",
		"water_quality_status_poor":       "Insuficiente",
	},
	"pl": {
		"water_quality":                   "Jakość wody",
		"water_sample_date":               "Data próbki",
		"water_quality_scale":             "Skala jakości",
		"water_quality_status_excellent":  "Doskonała",
		"water_quality_status_good":       "Dobra",
		"water_quality_status_sufficient": "Dostateczna",
		"water_quality_status_poor":       "Niedostateczna",
	},
	"de": {
		"water_quality":                   "Wasserqualität",
		"water_sample_date":               "Probenahme",
		"water_quality_scale":             "Qualitätsskala",
		"water_quality_status_excellent":  "Ausgezeichnet",
		"water_quality_status_good":       "Gut",
		"water_quality_status_sufficient": "Ausreichend",
		"water_quality_status_poor":       "Mangelhaft",
	},
	"fr": {
		"water_quality":                   "Qualité de l'eau",
		"water_sample_date":               "Date d'échantillon",
		"water_quality_scale":             "Échelle de qualité",
		"water_quality_status_excellent":  "Excellente",
		"water_quality_status_good":       "Bonne",
		"water_quality_status_sufficient": "Suffisante",
		"water_quality_status_poor":       "Insuffisante",
	},
	"it": {
		"water_quality":                   "Qualità dell'acqua",
		"water_sample_date":               "Data campione",
		"water_quality_scale":             "Scala di qualità",
		"water_quality_status_excellent":  "Eccellente",
		"water_quality_status_good":       "Buona",
		"water_quality_status_sufficient": "Sufficiente",
		"water_quality_status_poor":       "Insufficiente",
	},
	"hu": {
		"water_quality":                   "Vízminőség",
		"water_sample_date":               "Minta dátuma",
		"water_quality_scale":             "Minőségi skála",
		"water_quality_status_excellent":  "Kiváló",
		"water_quality_status_good":       "Jó",
		"water_quality_status_sufficient": "Megfelelő",
		"water_quality_status_poor":       "Elégtelen",
	},
}

var recentLocaleMessages = map[string]map[string]string{
	"es": {
		"recent_statistics":           "Últimos 10 días",
		"recent_short":                "10 días",
		"recent_stats_subtitle":       "Los gráficos muestran promedios diarios de los últimos 10 días.",
		"recent_temp_chart_title":     "Temperatura de los últimos 10 días (°C)",
		"recent_pressure_chart_title": "Presión de los últimos 10 días (hPa)",
		"recent_humidity_chart_title": "Humedad de los últimos 10 días (%)",
		"detail_charts_load_error":    "No se pudieron cargar los gráficos horarios para la fecha seleccionada.",
		"click_to_change_date":        "Haga clic para cambiar la fecha",
	},
	"pl": {
		"recent_statistics":           "Ostatnie 10 dni",
		"recent_short":                "10 dni",
		"recent_stats_subtitle":       "Wykresy pokazują dzienne średnie z ostatnich 10 dni.",
		"recent_temp_chart_title":     "Temperatura z ostatnich 10 dni (°C)",
		"recent_pressure_chart_title": "Ciśnienie z ostatnich 10 dni (hPa)",
		"recent_humidity_chart_title": "Wilgotność z ostatnich 10 dni (%)",
		"detail_charts_load_error":    "Nie udało się wczytać wykresów godzinowych dla wybranej daty.",
		"click_to_change_date":        "Kliknij, aby zmienić datę",
	},
	"de": {
		"recent_statistics":           "Letzte 10 Tage",
		"recent_short":                "10 Tage",
		"recent_stats_subtitle":       "Die Diagramme zeigen tägliche Durchschnittswerte der letzten 10 Tage.",
		"recent_temp_chart_title":     "Temperatur der letzten 10 Tage (°C)",
		"recent_pressure_chart_title": "Druck der letzten 10 Tage (hPa)",
		"recent_humidity_chart_title": "Luftfeuchtigkeit der letzten 10 Tage (%)",
		"detail_charts_load_error":    "Die stündlichen Diagramme für das ausgewählte Datum konnten nicht geladen werden.",
		"click_to_change_date":        "Klicken Sie, um das Datum zu ändern",
	},
	"fr": {
		"recent_statistics":           "10 derniers jours",
		"recent_short":                "10 jours",
		"recent_stats_subtitle":       "Les graphiques montrent les moyennes quotidiennes des 10 derniers jours.",
		"recent_temp_chart_title":     "Température des 10 derniers jours (°C)",
		"recent_pressure_chart_title": "Pression des 10 derniers jours (hPa)",
		"recent_humidity_chart_title": "Humidité des 10 derniers jours (%)",
		"detail_charts_load_error":    "Impossible de charger les graphiques horaires pour la date sélectionnée.",
		"click_to_change_date":        "Cliquez pour changer la date",
	},
	"it": {
		"recent_statistics":           "Ultimi 10 giorni",
		"recent_short":                "10 giorni",
		"recent_stats_subtitle":       "I grafici mostrano le medie giornaliere degli ultimi 10 giorni.",
		"recent_temp_chart_title":     "Temperatura degli ultimi 10 giorni (°C)",
		"recent_pressure_chart_title": "Pressione degli ultimi 10 giorni (hPa)",
		"recent_humidity_chart_title": "Umidità degli ultimi 10 giorni (%)",
		"detail_charts_load_error":    "Impossibile caricare i grafici orari per la data selezionata.",
		"click_to_change_date":        "Fai clic per cambiare la data",
	},
	"hu": {
		"recent_statistics":           "Utolsó 10 nap",
		"recent_short":                "10 nap",
		"recent_stats_subtitle":       "A grafikonok az elmúlt 10 nap napi átlagait mutatják.",
		"recent_temp_chart_title":     "Hőmérséklet az elmúlt 10 napban (°C)",
		"recent_pressure_chart_title": "Nyomás az elmúlt 10 napban (hPa)",
		"recent_humidity_chart_title": "Páratartalom az elmúlt 10 napban (%)",
		"detail_charts_load_error":    "Nem sikerült betölteni az órás grafikonokat a kiválasztott dátumhoz.",
		"click_to_change_date":        "Kattintson a dátum módosításához",
	},
}

var waveSourceLocaleMessages = map[string]map[string]string{
	"cs": {
		"wave_source_note": "Údaj o výšce vln vychází z offshore stanice PORTUS v oblasti Golf del Sur (Hm0) a nemusí přesně odpovídat podmínkám na konkrétní pláži.",
	},
	"en": {
		"wave_source_note": "Wave height is taken from an offshore PORTUS station in the Golf del Sur area (Hm0) and may not exactly match conditions on a specific beach.",
	},
	"es": {
		"wave_source_note": "La altura de ola se toma de una estación PORTUS mar adentro en la zona de Golf del Sur (Hm0) y puede no coincidir exactamente con las condiciones en una playa concreta.",
	},
	"pl": {
		"wave_source_note": "Wysokość fal pochodzi z morskiej stacji PORTUS w rejonie Golf del Sur (Hm0) i może nie odpowiadać dokładnie warunkom na konkretnej plaży.",
	},
	"de": {
		"wave_source_note": "Die Wellenhöhe stammt von einer offshore PORTUS-Station im Bereich Golf del Sur (Hm0) und kann von den Bedingungen an einem konkreten Strand abweichen.",
	},
	"fr": {
		"wave_source_note": "La hauteur des vagues provient d'une station PORTUS au large dans la zone de Golf del Sur (Hm0) et peut ne pas correspondre exactement aux conditions sur une plage donnée.",
	},
	"it": {
		"wave_source_note": "L'altezza delle onde proviene da una stazione PORTUS offshore nell'area di Golf del Sur (Hm0) e potrebbe non corrispondere esattamente alle condizioni su una spiaggia specifica.",
	},
	"hu": {
		"wave_source_note": "A hullámmagasság egy parttól távoli, Golf del Sur térségében lévő PORTUS állomás (Hm0) adata, ezért nem mindig egyezik pontosan egy adott strand viszonyaival.",
	},
}

func init() {
	monthNames := map[string][]string{
		"cs": {"Leden", "Únor", "Březen", "Duben", "Květen", "Červen", "Červenec", "Srpen", "Září", "Říjen", "Listopad", "Prosinec"},
		"en": {"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"},
		"es": {"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"},
		"pl": {"Styczeń", "Luty", "Marzec", "Kwiecień", "Maj", "Czerwiec", "Lipiec", "Sierpień", "Wrzesień", "Październik", "Listopad", "Grudzień"},
		"de": {"Januar", "Februar", "März", "April", "Mai", "Juni", "Juli", "August", "September", "Oktober", "November", "Dezember"},
		"fr": {"Janvier", "Février", "Mars", "Avril", "Mai", "Juin", "Juillet", "Août", "Septembre", "Octobre", "Novembre", "Décembre"},
		"it": {"Gennaio", "Febbraio", "Marzo", "Aprile", "Maggio", "Giugno", "Luglio", "Agosto", "Settembre", "Ottobre", "Novembre", "Dicembre"},
		"hu": {"Január", "Február", "Március", "Április", "Május", "Június", "Július", "Augusztus", "Szeptember", "Október", "November", "December"},
	}

	for locale, months := range monthNames {
		if _, ok := dictionary[locale]; !ok {
			continue
		}
		for i, monthName := range months {
			dictionary[locale][fmt.Sprintf("month_%d", i+1)] = monthName
		}
	}
}
