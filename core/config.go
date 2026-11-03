package core

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration values loaded from the environment.
type Config struct {
	Port      string
	Env       string
	JWTSecret string
	DBName    string
}

// LoadConfig reads the .env file (if present) and populates a Config struct
// with sensible defaults for any missing values.
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables / defaults")
	}

	cfg := &Config{
		Port:      getEnv("APP_PORT", "8080"),
		Env:       getEnv("APP_ENV", "debug"),
		JWTSecret: getEnv("JWT_SECRET", "default-secret-change-me"),
		DBName:    getEnv("DB_NAME", "app.db"),
	}

	return cfg
}

// IsDebug returns true when the application runs in debug mode.
func (c *Config) IsDebug() bool {
	return c.Env == "debug"
}

// getEnv reads a key from the environment with a fallback default.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// getEnvInt reads an integer key from the environment with a fallback default.
func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
