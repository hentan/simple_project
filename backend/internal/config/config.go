package config

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	envAppPort       = "APP_PORT"
	envVarDbHost     = "POSTGRES_HOST"
	envVarDbPort     = "POSTGRES_PORT"
	envVarDbUser     = "POSTGRES_USER"
	envVarDbPassword = "POSTGRES_PASSWORD"
	envVarDbName     = "POSTGRES_DB"
)

type Postgresql struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type Config struct {
	AppPort    string
	Postgresql Postgresql
}

func NewConfig(path string) *Config {
	_ = godotenv.Load()

	var db Postgresql = Postgresql{
		Host:     getEnv(envVarDbHost, "postgres"),
		Port:     getEnv(envVarDbPort, "5432"),
		Username: getEnv(envVarDbUser, "postgres"),
		Password: getEnv(envVarDbPassword, "postgres"),
		Database: getEnv(envVarDbName, "postgres"),
	}

	return &Config{
		AppPort:    getEnv(envAppPort, "8080"),
		Postgresql: db,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
