package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	envAppPort       = "APP_PORT"
	envVarDbHost     = "DB_HOST"
	envVarDbPort     = "DB_PORT"
	envVarDbUser     = "DB_USER"
	envVarDbPassword = "DB_Password"
	envVarDbName     = "DB_NAME"
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
	err := godotenv.Load(path)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var db Postgresql = Postgresql{
		Host:     getEnv(envVarDbHost, "localhost"),
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
