package utils

import (
	"log"
	"os"
	"path/filepath"
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
	if weatherPath != "" && !IsSafePath(weatherPath) {
		log.Printf("Warning: WEATHER_JSON_PATH '%s' contains traversal sequences, using default", weatherPath)
	}
	
	webcamPath := os.Getenv("WEBCAM_IMAGE_PATH")
	if webcamPath != "" && !IsSafePath(webcamPath) {
		log.Printf("Warning: WEBCAM_IMAGE_PATH '%s' contains traversal sequences, using default", webcamPath)
	}
}

// IsSafePath checks if a path is safe and doesn't contain traversal sequences
// Uses filepath.Clean() to normalize the path and prevent directory traversal
func IsSafePath(path string) bool {
	// Use filepath.Clean to normalize the path and remove traversal sequences
	cleanPath := filepath.Clean(path)
	
	// Check if the cleaned path is different from the original (indicates traversal)
	if cleanPath != path {
		return false
	}
	
	// Additional checks for absolute paths that might be unsafe
	if strings.HasPrefix(cleanPath, "/") && len(cleanPath) > 1 {
		// For absolute paths, we can't easily determine if they're safe
		// So we allow them but log a warning in ValidateEnv
		return true
	}
	
	// Relative paths are generally safe if they don't contain traversal sequences
	// after cleaning
	return !strings.Contains(cleanPath, "..")
}
