package config

import (
"log"
"os"

"github.com/joho/godotenv"
)

// Config holds all configuration values for the application
type Config struct {
	PiHoleURL      string
	PiHolePassword string
	Port           string
}

// Load loads configuration from .env file and environment variables
func Load() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{
		PiHoleURL:      getEnv("PIHOLE_URL", "http://localhost:8080/api"),
		PiHolePassword: getEnv("PIHOLE_PASSWORD", ""),
		Port:           getEnv("PORT", "8081"),
	}

	// Validate required fields
	if cfg.PiHolePassword == "" {
		log.Fatal("PIHOLE_PASSWORD environment variable is required")
	}

	return cfg
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
