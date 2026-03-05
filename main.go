package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/skybedy/laravel-tene.life/internal/alerts"
	"github.com/skybedy/laravel-tene.life/internal/i18n"
	"github.com/skybedy/laravel-tene.life/internal/pws"
	"github.com/skybedy/laravel-tene.life/internal/store"
	"github.com/skybedy/laravel-tene.life/internal/utils"
	"github.com/skybedy/laravel-tene.life/internal/water"
	"github.com/skybedy/laravel-tene.life/internal/waves"
	"github.com/skybedy/laravel-tene.life/internal/web"
)

//go:embed views
var viewsFS embed.FS

//go:embed public/js public/css public/images/tenelife-logo.png
var staticFS embed.FS

var db *sql.DB

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, checking parent directory...")
		// Try loading from parent for dev convenience in flat/nested structures
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println("No .env file found in parent either")
		}
	}

	if len(os.Args) > 1 && os.Args[1] == "collect:waves" {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err := waves.CollectLatestToDefaultPath(ctx); err != nil {
			log.Fatal("collect:waves failed: ", err)
		}
		log.Println("collect:waves completed: data/waves_latest.json")
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "collect:water" {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err := water.CollectLatestToDefaultPath(ctx); err != nil {
			log.Fatal("collect:water failed: ", err)
		}
		log.Println("collect:water completed: data/water_quality_latest.json")
		return
	}
	runPWSCollect := len(os.Args) > 1 && os.Args[1] == "collect:pws"

	// Validate environment variables
	utils.ValidateEnv()

	// Database Connection
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&timeout=5s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	log.Println("Connecting to database...")
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	// Verify connection
	if err := db.Ping(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Optimize connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Connected to database successfully!")

	// Initialize Store
	weatherStore := store.NewWeatherStore(db)
	emailNotifier := alerts.NewEmailNotifierFromEnv()

	if runPWSCollect {
		ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
		defer cancel()

		if err := pws.CollectLatestToDB(ctx, weatherStore); err != nil {
			sendPWSFailureAlert(emailNotifier, err, "manual")
			log.Fatal("collect:pws failed: ", err)
		}
		log.Println("collect:pws completed: pws_latest table updated")
		return
	}

	// In-app background collectors (interval-based only, no immediate run at startup).
	if envBool("WAVES_COLLECTOR_ENABLED", true) {
		startWavesCollectorLoop(context.Background())
	}
	if envBool("WATER_COLLECTOR_ENABLED", true) {
		startWaterCollectorLoop(context.Background())
	}
	if envBool("PWS_COLLECTOR_ENABLED", true) {
		startPWSCollectorLoop(context.Background(), weatherStore, emailNotifier)
	}

	// Initialize Handlers
	handler := web.NewHandler(weatherStore)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	// Simple logger for production (optional, can be disabled for max speed)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, latency=${latency_human}\n",
	}))

	// Security Middleware
	e.Use(middleware.Secure())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:csrf",
		Skipper: func(c echo.Context) bool {
			// API endpoints are consumed by scripts/JS and do not post form tokens.
			return len(c.Path()) >= 5 && c.Path()[:5] == "/api/"
		},
	}))

	// Dynamic Webcam Image serving
	webcamPath := utils.EnvPathOrDefault("WEBCAM_IMAGE_PATH", "public/images/tenelife.jpg")
	e.File("/images/tenelife.jpg", webcamPath)

	// Static Files from Embed
	publicFS, _ := fs.Sub(staticFS, "public")
	e.StaticFS("/js", echo.MustSubFS(publicFS, "js"))
	e.StaticFS("/css", echo.MustSubFS(publicFS, "css"))
	e.FileFS("/images/tenelife-logo.png", "images/tenelife-logo.png", publicFS)

	// We still need to serve local images for other files
	e.Static("/images", "public/images")
	e.Static("/files", "public/files")

	// Template Renderer using Embed
	renderer := &web.TemplateRenderer{
		Templates: template.Must(template.New("").Funcs(template.FuncMap{
			"localeURL":           i18n.LocaleURL,
			"monthName":           i18n.MonthName,
			"languageFlag":        i18n.LanguageFlag,
			"waterQualityStatus":  i18n.WaterQualityStatusLabel,
			"waterQualityTooltip": i18n.WaterQualityTooltip,
			"waveDirectionLabel":  i18n.WaveDirectionLabel,
			"f1": func(v *float64) string {
				if v == nil {
					return "--"
				}
				return fmt.Sprintf("%.1f", *v)
			},
			"f0": func(v *float64) string {
				if v == nil {
					return "--"
				}
				return fmt.Sprintf("%.0f", *v)
			},
			"dateOnly": func(v string) string {
				if len(v) >= 10 {
					return v[:10]
				}
				return v
			},
			"shortDate": func(v string) string {
				if len(v) < 10 {
					return v
				}
				var y, m, d int
				if _, err := fmt.Sscanf(v[:10], "%d-%d-%d", &y, &m, &d); err != nil {
					return v[:10]
				}
				return fmt.Sprintf("%d.%d", d, m)
			},
			"shortDateYear": func(v string) string {
				if len(v) < 10 {
					return v
				}
				var y, m, d int
				if _, err := fmt.Sscanf(v[:10], "%d-%d-%d", &y, &m, &d); err != nil {
					return v[:10]
				}
				return fmt.Sprintf("%d.%d.%02d", d, m, y%100)
			},
		}).ParseFS(viewsFS, "views/*.html", "views/statistics/*.html")),
	}

	e.Renderer = renderer

	// Routes (default locale: cs)
	e.GET("/", handler.IndexHandler)
	e.GET("/webcam/big", handler.WebcamBigHandler)
	e.GET("/canary-islands-temperatures", handler.PWSMapHandler)
	e.GET("/statistics", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/statistics/daily")
	})
	e.GET("/statistics/daily", handler.DailyStatisticsHandler)
	e.GET("/statistics/recent", handler.RecentStatisticsHandler)
	e.GET("/statistics/weekly", handler.WeeklyStatisticsHandler)
	e.GET("/statistics/monthly", handler.MonthlyStatisticsHandler)
	e.GET("/statistics/annual", handler.AnnualStatisticsHandler)

	// Routes with locale prefix: /en, /de, /fr ...
	localized := e.Group("/:locale")
	localized.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !i18n.IsSupportedLocale(c.Param("locale")) {
				return echo.ErrNotFound
			}
			return next(c)
		}
	})
	localized.GET("", handler.IndexHandler)
	localized.GET("/", handler.IndexHandler)
	localized.GET("/webcam/big", handler.WebcamBigHandler)
	localized.GET("/canary-islands-temperatures", handler.PWSMapHandler)
	localized.GET("/statistics", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, i18n.LocaleURL(c.Param("locale"), "/statistics/daily"))
	})
	localized.GET("/statistics/daily", handler.DailyStatisticsHandler)
	localized.GET("/statistics/recent", handler.RecentStatisticsHandler)
	localized.GET("/statistics/weekly", handler.WeeklyStatisticsHandler)
	localized.GET("/statistics/monthly", handler.MonthlyStatisticsHandler)
	localized.GET("/statistics/annual", handler.AnnualStatisticsHandler)

	// API and service routes
	e.GET("/webcam/image.jpg", handler.WebcamImageHandler) // New dynamic route
	e.GET("/api/weather/hourly", handler.GetHourlyDataHandler)
	e.GET("/api/home", handler.GetHomeDataHandler)
	e.GET("/api/tenerife/pws-latest", handler.GetPWSLatestHandler)
	e.GET("/api/tides", handler.GetTidesHandler)
	e.GET("/debug/wu-usage", func(c echo.Context) error {
		return c.JSON(http.StatusOK, pws.GetWUUsageReport())
	})

	// Health check endpoint
	e.GET("/health", handler.HealthCheckHandler)

	// API Statistics
	e.GET("/api/weather/daily", handler.GetDailyDataHandler)
	e.GET("/api/weather/monthly-daily", handler.GetMonthlyDailyDataHandler)
	e.GET("/api/weather/weekly", handler.GetWeeklyDataHandler)
	e.GET("/api/weather/monthly", handler.GetMonthlyDataHandler)
	e.GET("/api/weather/annual", handler.GetAnnualDataHandler)

	// API Data Ingestion
	e.POST("/api/weather/sea-temperature", handler.StoreSeaTemperatureHandler)
	e.POST("/api/camera/upload", handler.CameraUploadHandler, middleware.BodyLimit("10M"))

	// Start Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

func sendPWSFailureAlert(emailNotifier *alerts.EmailNotifier, err error, source string) {
	if emailNotifier == nil || err == nil {
		return
	}
	key := "pws_collect_failed"
	subject := "PWS collector failed"
	lower := strings.ToLower(err.Error())
	if strings.Contains(lower, "access denied") || strings.Contains(lower, "status 401") || strings.Contains(lower, "status 403") || strings.Contains(lower, "status 429") {
		key = "pws_access_denied"
		subject = "PWS API access denied/rate-limited"
	}
	body := fmt.Sprintf(
		"Time (UTC): %s\nSource: %s\nError: %v\nEnvironment: %s\n",
		time.Now().UTC().Format(time.RFC3339),
		source,
		err,
		os.Getenv("APP_ENV"),
	)
	emailNotifier.Notify(key, subject, body)
}

func startWavesCollectorLoop(ctx context.Context) {
	interval := time.Duration(envInt("WAVES_COLLECT_INTERVAL_MINUTES", 15)) * time.Minute
	if interval <= 0 {
		interval = 15 * time.Minute
	}
	log.Printf("waves collector loop enabled: interval=%s (first run on next tick)", interval)

	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("waves collector loop stopped")
				return
			case <-ticker.C:
				runCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
				if err := waves.CollectLatestToDefaultPath(runCtx); err != nil {
					log.Printf("waves collector failed: %v", err)
				} else {
					log.Println("waves collector updated: data/waves_latest.json")
				}
				cancel()
			}
		}
	}()
}

func startWaterCollectorLoop(ctx context.Context) {
	interval := time.Duration(envInt("WATER_COLLECT_INTERVAL_HOURS", 24)) * time.Hour
	if interval <= 0 {
		interval = 24 * time.Hour
	}
	log.Printf("water collector loop enabled: interval=%s (first run on next tick)", interval)

	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("water collector loop stopped")
				return
			case <-ticker.C:
				runCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
				if err := water.CollectLatestToDefaultPath(runCtx); err != nil {
					log.Printf("water collector failed: %v", err)
				} else {
					log.Println("water collector updated: data/water_quality_latest.json")
				}
				cancel()
			}
		}
	}()
}

func startPWSCollectorLoop(ctx context.Context, weatherStore *store.WeatherStore, emailNotifier *alerts.EmailNotifier) {
	if !pws.APIKeyConfigured() {
		log.Println("pws collector loop disabled: WEATHER_COM_API_KEY is empty")
		return
	}

	log.Println("pws paced collector loop enabled")
	go func() {
		if err := pws.RunPacedCollector(ctx, weatherStore); err != nil && err != context.Canceled {
			log.Printf("pws paced collector stopped: %v", err)
			sendPWSFailureAlert(emailNotifier, err, "scheduler")
		}
	}()
}

func envBool(key string, fallback bool) bool {
	raw := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if raw == "" {
		return fallback
	}
	switch raw {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func envInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return v
}
