package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr          string
	ExternalBase  string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Addr:         getEnv("APP_ADDR", ":8085"),
		ExternalBase: getEnv("EXTERNAL_BASE_URL", "https://example.com"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
