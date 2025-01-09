package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct holds application configuration values
type Config struct {
	Port     string
	LogLevel string
}

// LoadConfig reads configuration from environment variables or a `.env` file
func LoadConfig() (*Config, error) {
	// Load environment variables from .env file (if exists)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration values
	cfg := &Config{
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	return cfg, nil
}

// getEnv fetches environment variables with a fallback default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
