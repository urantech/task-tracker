package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DbUrl       string
	ConsumerCfg ConsumerConfig
	OutboxCfg   OutboxConfig
	SmtpCfg     SmtpConfig
}

type ConsumerConfig struct {
	Brokers        []string
	WelcomeTopic   string
	WelcomeGroupID string
	ReportTopic    string
	ReportGroupID  string
}

type OutboxConfig struct {
	WorkersCount int
	IntervalSec  int
	BatchSize    int
}

type SmtpConfig struct {
	Host       string
	Port       string
	Username   string
	Password   string
	From       string
	TimeoutSec int
}

func MustLoad() *Config {
	consumerCfg := ConsumerConfig{
		Brokers:        strings.Split(getEnv("KAFKA_BROKERS"), ","),
		WelcomeTopic:   getEnv("USERS_REGISTRATION_TOPIC"),
		WelcomeGroupID: getEnv("USERS_REGISTRATION_GROUP_ID"),
		ReportTopic:    getEnv("DAILY_REPORT_TOPIC"),
		ReportGroupID:  getEnv("DAILY_REPORT_GROUP_ID"),
	}

	smtpCfg := SmtpConfig{
		Host:       getEnv("SMTP_HOST"),
		Port:       getEnv("SMTP_PORT"),
		Username:   os.Getenv("SMTP_USERNAME"),
		Password:   os.Getenv("SMTP_PASSWORD"),
		From:       getEnv("SMTP_FROM"),
		TimeoutSec: atoi(getEnv("SMTP_TIMEOUT_SEC")),
	}

	outboxCfg := OutboxConfig{
		WorkersCount: atoi(getEnv("OUTBOX_WORKERS_COUNT")),
		BatchSize:    atoi(getEnv("OUTBOX_BATCH_SIZE")),
		IntervalSec:  atoi(getEnv("OUTBOX_INTERVAL_SEC")),
	}

	return &Config{
		DbUrl:       getEnv("EMAIL_SENDER_DB_URL"),
		ConsumerCfg: consumerCfg,
		OutboxCfg:   outboxCfg,
		SmtpCfg:     smtpCfg,
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)

	if val == "" {
		log.Fatalf("Required env variable %s is missing", key)
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
