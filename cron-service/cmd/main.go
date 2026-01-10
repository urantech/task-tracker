package main

import (
	"cron-service/internal/config"
	"cron-service/migrations"
	"cron-service/pkg/postgres"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.MustLoad()

	db, err := postgres.NewConnection(cfg.DbUrl)
	if err != nil {
		log.Fatalf("Failed connect to database: %v", err)
	}

	runMigrations(db)
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
