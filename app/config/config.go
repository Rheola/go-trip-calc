package config

import (
	"os"
)

type Config struct {
	DbUrl  string
	ApiKey string
	Port   string
}

func New() *Config {
	return &Config{
		DbUrl:  getEnv("DB_URL", ""),
		ApiKey: getEnv("API_KEY", ""),
		Port:   getEnv("PORT", "3000"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
