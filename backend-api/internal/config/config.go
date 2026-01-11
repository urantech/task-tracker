package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	HttpPort  string
	GrpcPort  string
	DbUrl     string
	JwtSecret string
	Brokers   []string
}

func MustLoad() *Config {
	inDocker := os.Getenv("DOCKER") == "true"

	httpPort := getEnvOrFallback("BACKEND_API_HTTP_PORT", "8080", inDocker)
	grpcPort := getEnvOrFallback("BACKEND_API_GRPC_PORT", "50051", inDocker)
	dbUrl := getEnvOrFallback("BACKEND_API_DB_URL", "postgres://task-tracker-user:task-tracker-password@localhost:5432/task-tracker-db?sslmode=disable", inDocker)
	jwtSecret := getEnvOrFallback("JWT_SECRET", "K7gNU3sdo+OL0wNhqoVWhr3g6s1xYv72ol/pe/Unolp=", inDocker)
	brokers := strings.Split(getEnvOrFallback("KAFKA_BROKERS", "localhost:9092", inDocker), ",")

	return &Config{
		HttpPort:  httpPort,
		GrpcPort:  grpcPort,
		DbUrl:     dbUrl,
		JwtSecret: jwtSecret,
		Brokers:   brokers,
	}
}

func getEnvOrFallback(key, fallback string, inDocker bool) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}

	if inDocker {
		log.Fatalf("Required env variable %s is missing", key)
	}

	return fallback
}
