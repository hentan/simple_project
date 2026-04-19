package config

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	envAppPort             = "APP_PORT"
	envVarPostgresHost     = "POSTGRES_HOST"
	envVarPostgresPort     = "POSTGRES_PORT"
	envVarPostgresUser     = "POSTGRES_USER"
	envVarPostgresPassword = "POSTGRES_PASSWORD"
	envVarPostgresName     = "POSTGRES_DB"
)

type Postgres struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type Config struct {
	AppPort  string
	Postgres Postgres
}

func NewConfig(path string) *Config {
	_ = godotenv.Load()

	db := Postgres{
		Host:     getEnv(envVarPostgresHost, "postgres"),
		Port:     getEnv(envVarPostgresPort, "5432"),
		Username: getEnv(envVarPostgresUser, "postgres"),
		Password: getEnv(envVarPostgresPassword, "postgres"),
		Database: getEnv(envVarPostgresName, "postgres"),
	}

	return &Config{
		AppPort:  getEnv(envAppPort, "8080"),
		Postgres: db,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
