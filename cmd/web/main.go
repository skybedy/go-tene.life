package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/skybedy/laravel-tene.life/internal/store"
	"github.com/skybedy/laravel-tene.life/internal/web"
)

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
	log.Println("Connected to database successfully!")

	// Initialize Store
	weatherStore := store.NewWeatherStore(db)

	// Initialize Handlers
	handler := web.NewHandler(weatherStore)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Security Middleware
	e.Use(middleware.Secure())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:csrf",
	}))

	// Static Files
	e.Static("/images", "public/images")
	e.Static("/files", "public/files")
	e.Static("/js", "public/js")

	// Template Renderer
	renderer := &web.TemplateRenderer{
		Templates: template.Must(template.ParseGlob("views/*.html")),
	}
	// Add subdirectories manually or use a better search if needed
	template.Must(renderer.Templates.ParseGlob("views/statistics/*.html"))

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
	e.Logger.Fatal(e.Start(":8080"))
}
