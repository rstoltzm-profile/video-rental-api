package config

import "os"

type Config struct {
	DatabaseURL string
	Port        string
}

func LoadConfig() Config {
	return Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        getEnvOrDefault("PORT", "8080"),
	}
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
