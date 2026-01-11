package config

import (
	"log"
	"os"
)

type Config struct {
	DbUrl    string
	GrpcAddr string
}

func MustLoad() *Config {
	inDocker := os.Getenv("DOCKER") == "true"

	dbUrl := getEnvOrFallback("CRON_SERVICE_DB_URL", "postgres://cron-service-user:cron-service-password@localhost:5433/cron-service-db?sslmode=disable", inDocker)
	grpcAddr := getEnvOrFallback("BACKEND_API_GRPC_ADDR", "localhost:50051", inDocker)

	return &Config{
		DbUrl:    dbUrl,
		GrpcAddr: grpcAddr,
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
