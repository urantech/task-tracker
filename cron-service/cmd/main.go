package main

import (
	"context"
	"cron-service/internal/config"
	"cron-service/internal/job"
	"cron-service/migrations"
	"cron-service/pkg/postgres"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	taskv1 "api-proto/gen/go/task/v1"
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

	conn, err := grpc.NewClient(cfg.GrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}

	client := taskv1.NewTaskAnalyticsServiceClient(conn)
	grpcClient := job.NewGrpcClient(client)

	storage := job.NewStorage(db)
	service := job.NewService(storage, grpcClient)
	scheduler := gocron.NewScheduler(location)

	rootCtx, cancel := context.WithCancel(context.Background())

	err = service.SetupScheduler(rootCtx, scheduler)
	if err != nil {
		log.Fatalf("Failed to setup scheduler: %v", err)
	}

	defer cancel()

	scheduler.StartAsync()

	shutdownCtx, stop := signal.NotifyContext(rootCtx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-shutdownCtx.Done()

	shutdown(scheduler, conn)
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

func shutdown(scheduler *gocron.Scheduler, conn *grpc.ClientConn) {
	scheduler.Stop()

	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			log.Printf("connection close error: %v", closeErr)
		}
	}()

	log.Print("App shutdown gracefully")
}
