package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	HttpPort     string
	GrpcPort     string
	DbUrl        string
	JwtSecret    string
	Brokers      []string
	WelcomeTopic string
	ReportTopic  string
}

func MustLoad() *Config {
	httpPort := getEnv("BACKEND_API_HTTP_PORT")
	grpcPort := getEnv("BACKEND_API_GRPC_PORT")
	dbUrl := getEnv("BACKEND_API_DB_URL")
	jwtSecret := getEnv("JWT_SECRET")
	brokers := strings.Split(getEnv("KAFKA_BROKERS"), ",")
	welcomeTopic := getEnv("USERS_REGISTRATION_TOPIC")
	reportTopic := getEnv("DAILY_REPORT_TOPIC")

	return &Config{
		HttpPort:     httpPort,
		GrpcPort:     grpcPort,
		DbUrl:        dbUrl,
		JwtSecret:    jwtSecret,
		Brokers:      brokers,
		WelcomeTopic: welcomeTopic,
		ReportTopic:  reportTopic,
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)

	if val == "" {
		log.Fatalf("%s env is not set", key)
	}

	return val
}
