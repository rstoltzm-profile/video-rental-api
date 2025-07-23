package config

import "os"

type Config struct {
	DatabaseURL string
	Port        string
	APIKey      string
}

func LoadConfig() Config {
	return Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        getEnvOrDefault("PORT", "8080"),
		APIKey:      getEnvOrDefault("API_KEY", "default-dev-key-123"),
	}
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
