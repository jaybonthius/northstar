package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host string
	Port string
}

func Load() *Config {
	godotenv.Load()

	getEnv := func(key, fallback string) string {
		if val, ok := os.LookupEnv(key); ok {
			return val
		}
		return fallback
	}

	return &Config{
		Host: getEnv("HOST", "0.0.0.0"),
		Port: getEnv("PORT", "8080"),
	}
}