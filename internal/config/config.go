package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	envAppPort = "APP_PORT"
)

type Config struct {
	AppPort string
}

func NewConfig(path string) *Config {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		AppPort: getEnv(envAppPort, "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
