package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbUrl string
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "../.env"
	}

	if err := godotenv.Load(configPath); err != nil {
		log.Printf("Note: .env file not loaded from %s", configPath)
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		dbUrl = "postgres://cron-service-user:cron-service-password@localhost:5433/cron-service-db?sslmode=disable"
	}

	return &Config{
		DbUrl: dbUrl,
	}
}
