package main

import (
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/skybedy/laravel-tene.life/internal/store"
	"github.com/skybedy/laravel-tene.life/internal/web"
)

//go:embed views/*.html views/statistics/*.html
var viewsFS embed.FS

//go:embed public/js/*.js public/images/tenelife-logo.png
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

	// Initialize Handlers
	handler := web.NewHandler(weatherStore)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())

	// Simple logger for production (optional, can be disabled for max speed)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, latency=${latency_human}\n",
	}))

	// Security Middleware
	e.Use(middleware.Secure())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:csrf",
	}))

	// Static Files from Embed
	publicFS, _ := fs.Sub(staticFS, "public")
	e.StaticFS("/js", echo.MustSubFS(publicFS, "js"))
	e.FileFS("/images/tenelife-logo.png", "images/tenelife-logo.png", publicFS)

	// We still need to serve local images for the webcam and other files
	e.Static("/images", "public/images")
	e.Static("/files", "public/files")

	// Dynamic Webcam Image serving
	webcamPath := os.Getenv("WEBCAM_IMAGE_PATH")
	if webcamPath == "" {
		webcamPath = "public/images/tenelife.jpg"
	}
	e.File("/images/tenelife.jpg", webcamPath)

	// Template Renderer using Embed
	renderer := &web.TemplateRenderer{
		Templates: template.Must(template.ParseFS(viewsFS, "views/*.html", "views/statistics/*.html")),
	}

	e.Renderer = renderer

	// Routes
	e.GET("/", handler.IndexHandler)
	e.GET("/webcam/big", handler.WebcamBigHandler)
	e.GET("/api/weather/hourly", handler.GetHourlyDataHandler)

	// Statistics
	e.GET("/statistics", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/statistics/daily")
	})
	e.GET("/statistics/daily", handler.DailyStatisticsHandler)
	e.GET("/statistics/weekly", handler.WeeklyStatisticsHandler)
	e.GET("/statistics/monthly", handler.MonthlyStatisticsHandler)
	e.GET("/statistics/annual", handler.AnnualStatisticsHandler)

	// API Statistics
	e.GET("/api/weather/daily", handler.GetDailyDataHandler)
	e.GET("/api/weather/weekly", handler.GetWeeklyDataHandler)
	e.GET("/api/weather/monthly", handler.GetMonthlyDataHandler)
	e.GET("/api/weather/annual", handler.GetAnnualDataHandler)

	// Start Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
