package main

import (
	"context"
	"database/sql"
	"email-sender/internal/config"
	"email-sender/internal/email"
	"email-sender/internal/kafka"
	"email-sender/internal/outbox"
	"email-sender/migrations"
	"email-sender/pkg/postgres"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

var runMigrationsFlag bool

func init() {
	flag.BoolVar(&runMigrationsFlag, "migrate", true, "Run database migrations on startup")
}

func main() {
	flag.Parse()
	
	cfg := config.MustLoad()

	db, err := postgres.NewConnection(cfg.DbUrl)
	if err != nil {
		log.Fatalf("Failed connect to database: %v", err)
	}

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runMigrations(db)

	outboxStorage := outbox.NewStorage(db)
	welcomeOutboxHandler := outbox.NewHandler(outboxStorage, "WELCOME")
	reportOutboxHandler := outbox.NewHandler(outboxStorage, "REPORT")

	welcomeReader := kafka.NewReader(
		cfg.ConsumerCfg.Brokers,
		cfg.ConsumerCfg.WelcomeTopic,
		cfg.ConsumerCfg.WelcomeGroupID,
		welcomeOutboxHandler,
	)

	reportReader := kafka.NewReader(
		cfg.ConsumerCfg.Brokers,
		cfg.ConsumerCfg.ReportTopic,
		cfg.ConsumerCfg.ReportGroupID,
		reportOutboxHandler,
	)

	sender := email.NewSender(cfg)
	processor := outbox.NewProcessor(outboxStorage, sender, cfg)

	shutdownCtx, stop := signal.NotifyContext(rootCtx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	errCh := make(chan error, 2)

	go func() {
		if err := welcomeReader.Run(shutdownCtx); err != nil {
			errCh <- err
		}
	}()

	go func() {
		if err := reportReader.Run(shutdownCtx); err != nil {
			errCh <- err
		}
	}()

	go func() {
		if err := processor.Run(shutdownCtx); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-shutdownCtx.Done():
		log.Println("shutdown signal received")
	case err := <-errCh:
		log.Printf("background process failed: %v", err)
		cancel()
	}

	log.Println("Shutdown successfully")
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
