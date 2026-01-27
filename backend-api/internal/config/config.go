package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ApiConfig      ApiConfig
	DbUrl          string
	JwtSecret      string
	ProducerConfig ProducerConfig
}

type ApiConfig struct {
	HttpPort string
	GrpcPort string
}

type ProducerConfig struct {
	Brokers           []string
	WelcomeTopic      string
	ReportTopic       string
	NumPartitions     int
	ReplicationFactor int
}

func MustLoad() *Config {
	apiCfg := ApiConfig{
		HttpPort: getEnv("BACKEND_API_HTTP_PORT"),
		GrpcPort: getEnv("BACKEND_API_GRPC_PORT"),
	}

	producerCfg := ProducerConfig{
		Brokers:           strings.Split(getEnv("KAFKA_BROKERS"), ","),
		WelcomeTopic:      getEnv("USERS_REGISTRATION_TOPIC"),
		ReportTopic:       getEnv("DAILY_REPORT_TOPIC"),
		NumPartitions:     atoi(getEnv("NUM_PARTITIONS")),
		ReplicationFactor: atoi(getEnv("REPLICATION_FACTOR")),
	}

	return &Config{
		ApiConfig:      apiCfg,
		DbUrl:          getEnv("BACKEND_API_DB_URL"),
		JwtSecret:      getEnv("JWT_SECRET"),
		ProducerConfig: producerCfg,
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)

	if val == "" {
		log.Fatalf("%s env is not set", key)
	}

	return val
}

func atoi(key string) int {
	val, err := strconv.Atoi(key)
	if err != nil {
		log.Fatalf("Invalid data type for env variable %s", key)
	}

	return val
}
