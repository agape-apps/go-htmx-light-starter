package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration values.
type Config struct {
	Port string // Server port
	// Add other configuration fields here as needed
}

// LoadConfig loads configuration from environment variables.
// It first attempts to load a .env file if present.
func LoadConfig() (*Config, error) {
	// Attempt to load .env file. Errors are ignored if the file doesn't exist.
	err := godotenv.Load()
	if err != nil {
		// Log only if it's not a "file not found" error, otherwise it's expected
		if !os.IsNotExist(err) {
			log.Printf("Warning: Error loading .env file: %v", err)
		}
	}

	cfg := Config{}

	// Load Port, defaulting to "8080" if not set
	cfg.Port = os.Getenv("PORT")
	if cfg.Port == "" {
		cfg.Port = "8080" // Default port
		log.Println("PORT environment variable not set, using default 8080")
	}

	// Load other variables here using os.Getenv("VAR_NAME")
	// Alternatively add github.com/kelseyhightower/envconfig for more complex configuration.
	//
	// Example:
	// cfg.Host = os.Getenv("HOST")
	// if cfg.Host == "" {
	// 	 cfg.Host = "localhost" // Default host
	// }

	return &cfg, nil // No error returned in this simplified version
}
