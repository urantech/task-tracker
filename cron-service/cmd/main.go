package main

import (
	"context"
	"cron-service/internal/config"
	"cron-service/internal/job"
	"cron-service/migrations"
	"cron-service/pkg/postgres"
	"database/sql"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.MustLoad()

	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}

	db, err := postgres.NewConnection(cfg.DbUrl)
	if err != nil {
		log.Fatalf("Failed connect to database: %v", err)
	}

	runMigrations(db)

	storage := job.NewStorage(db)
	service := job.NewService(storage)
	scheduler := gocron.NewScheduler(location)

	err = service.SetupScheduler(context.Background(), scheduler)
	if err != nil {
		log.Fatalf("Failed to setup scheduler: %v", err)
	}

	scheduler.StartAsync()
}

func runMigrations(db *sql.DB) {
	goose.SetBaseFS(migrations.EmbedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	log.Println("Running migrations...")

	if err := goose.Up(db, "."); err != nil {
		log.Fatal(err)
	}

	log.Println("Migrations completed successfully")
}
