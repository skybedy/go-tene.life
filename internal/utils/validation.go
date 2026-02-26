package utils

import (
	"log"
	"os"
	"strings"
)

// ValidateEnv validates environment variables and ensures they meet security requirements
func ValidateEnv() {
	// Validate database credentials
	if os.Getenv("DB_PASSWORD") == "" {
		log.Fatal("DB_PASSWORD is not set")
	}
	
	if os.Getenv("DB_USERNAME") == "" {
		log.Fatal("DB_USERNAME is not set")
	}
	
	if os.Getenv("DB_HOST") == "" {
		log.Fatal("DB_HOST is not set")
	}
	
	if os.Getenv("DB_DATABASE") == "" {
		log.Fatal("DB_DATABASE is not set")
	}
	
	// Validate file paths to prevent path traversal attacks
	weatherPath := os.Getenv("WEATHER_JSON_PATH")
	if weatherPath != "" && !isSafePath(weatherPath) {
		log.Fatal("Invalid WEATHER_JSON_PATH: potential path traversal attack")
	}
	
	webcamPath := os.Getenv("WEBCAM_IMAGE_PATH")
	if webcamPath != "" && !isSafePath(webcamPath) {
		log.Fatal("Invalid WEBCAM_IMAGE_PATH: potential path traversal attack")
	}
}

// IsSafePath checks if a path is safe and doesn't contain traversal sequences
func IsSafePath(path string) bool {
	// Check for path traversal sequences
	if strings.Contains(path, "../") || strings.Contains(path, "..\\") {
		return false
	}
	
	// Check if path starts with expected directories
	if strings.HasPrefix(path, "/var/www/") || 
	   strings.HasPrefix(path, "/public/") ||
	   strings.HasPrefix(path, "public/") {
		return true
	}
	
	return false
}
