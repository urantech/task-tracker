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
	dbUrl := getEnv("CRON_SERVICE_DB_URL")
	grpcAddr := getEnv("BACKEND_API_GRPC_ADDR")

	return &Config{
		DbUrl:    dbUrl,
		GrpcAddr: grpcAddr,
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)

	if val == "" {
		log.Fatalf("Required env variable %s is missing", key)
	}

	return val
}
