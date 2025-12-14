package config

import (
	"flag"
	"fmt"
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

// Load loads configuration from command-line flags, .env file, or environment variables
// Priority: command-line flags > environment variables > .env file > defaults
func Load() *Config {
	// Define command-line flags
	piholeURL := flag.String("pihole-url", "", "Pi-hole API URL (e.g., http://192.168.1.100/admin/api.php)")
	piholePassword := flag.String("pihole-password", "", "Pi-hole API password (required)")
	port := flag.String("port", "", "MCP server port (default: 8081)")
	// Custom usage message
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options]\n\n", os.Args[0])
		fmt.Println("Options:")
		fmt.Println("  --pihole-url string")
		fmt.Println("    \tPi-hole API URL (e.g., http://192.168.1.100/api)")
		fmt.Println("  --pihole-password string")
		fmt.Println("    \tPi-hole API password (required)")
		fmt.Println("  --port string")
		fmt.Println("    \tMCP server port (default: 8081)")
		fmt.Println("\nConfiguration priority: command-line flags > environment variables > .env file > defaults")
		fmt.Println("\nEnvironment variables:")
		fmt.Println("  PIHOLE_URL         Pi-hole API URL")
		fmt.Println("  PIHOLE_PASSWORD    Pi-hole API password (required)")
		fmt.Println("  PORT               MCP server port")
	}

	flag.Parse()

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables or flags")
	}

	cfg := &Config{
		PiHoleURL:      getConfigValue(*piholeURL, "PIHOLE_URL", "http://localhost:8080/api"),
		PiHolePassword: getConfigValue(*piholePassword, "PIHOLE_PASSWORD", ""),
		Port:           getConfigValue(*port, "PORT", "8081"),
	}

	// Validate required fields
	if cfg.PiHolePassword == "" {
		log.Fatal("PIHOLE_PASSWORD is required (set via --pihole-password flag or PIHOLE_PASSWORD environment variable)")
	}

	return cfg
}

// getConfigValue gets a configuration value with priority: flag > env var > default
func getConfigValue(flagValue, envKey, defaultValue string) string {
	// Priority 1: Command-line flag
	if flagValue != "" {
		return flagValue
	}
	// Priority 2: Environment variable
	if envValue := os.Getenv(envKey); envValue != "" {
		return envValue
	}
	// Priority 3: Default value
	return defaultValue
}
