package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DbUrl     string
	JwtSecret string
	Brokers   []string
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "../.env"
	}

	if err := godotenv.Load(configPath); err != nil {
		log.Printf("Note: .env file not loaded from %s", configPath)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		dbUrl = "postgres://task-tracker-user:task-tracker-password@localhost:5432/task-tracker-db?sslmode=disable"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("jwt secret not found")
	}

	rawBrokers := os.Getenv("KAFKA_BROKERS")
	if rawBrokers == "" {
		rawBrokers = "localhost:9092"
	}

	brokers := strings.Split(rawBrokers, ",")

	return &Config{
		Port:      port,
		DbUrl:     dbUrl,
		JwtSecret: jwtSecret,
		Brokers:   brokers,
	}
}
