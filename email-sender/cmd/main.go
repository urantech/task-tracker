package main

import (
	"context"
	"email-sender/internal/config"
	"email-sender/internal/email"
	"email-sender/internal/kafka"
	"email-sender/internal/outbox"
	"email-sender/pkg/postgres"
	"email-sender/pkg/postgres/migrations"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
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

	if runMigrationsFlag {
		log.Println("Migration flag is ON, running migrations...")

		m := migrations.NewMigrator(db)
		m.Run()
	} else {
		log.Println("Migration flag is OFF, skipping migrations")
	}

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
