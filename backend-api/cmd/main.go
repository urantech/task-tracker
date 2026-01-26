package main

import (
	taskv1 "api-proto/gen/go/task/v1"
	"backend-api/internal/auth"
	"backend-api/internal/config"
	"backend-api/internal/infra"
	"backend-api/internal/task"
	"backend-api/internal/user"
	"backend-api/migrations"
	"backend-api/pkg/postgres"
	"context"
	"database/sql"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"google.golang.org/grpc"
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
		runMigrations(db)
	} else {
		log.Println("Migration flag is OFF, skipping migrations")
	}

	createTopics(cfg)

	userProducer := user.NewProducer(cfg.Brokers)
	userStorage := user.NewStorage(db)
	userService := user.NewService(userStorage, userProducer)
	userHandler := user.NewHandler(userService)

	authService := auth.NewService(userStorage, cfg.JwtSecret)
	authHandler := auth.NewHandler(authService)

	taskProducer := task.NewProducer(cfg.Brokers)
	taskStorage := task.NewStorage(db)
	taskService := task.NewService(taskStorage, taskProducer)
	taskHandler := task.NewHandler(taskService)
	taskGrpcHandler := task.NewGrpcHandler(taskService)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Group(func(r chi.Router) {
		r.Post("/users/register", userHandler.Register)
		r.Post("/auth/login", authHandler.Login)
	})

	router.Group(func(r chi.Router) {
		r.Use(authService.AuthMiddleware)
		r.Get("/users/user", userHandler.GetCurrentUser)
		r.Post("/tasks", taskHandler.CreateTask)
		r.Get("/tasks", taskHandler.List)
		r.Patch("/tasks/{id}", taskHandler.UpdateTask)
	})

	httpSrv := &http.Server{
		Addr:    ":" + cfg.HttpPort,
		Handler: router,
	}

	grpcSrv := grpc.NewServer()
	taskv1.RegisterTaskAnalyticsServiceServer(grpcSrv, taskGrpcHandler)

	startGrpcServer(grpcSrv, cfg.GrpcPort)
	startHttpServer(httpSrv)

	shutdownCtx, stop := signal.NotifyContext(rootCtx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-shutdownCtx.Done()

	shutdownServers(httpSrv, grpcSrv)
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

func createTopics(cfg *config.Config) {
	topics := []string{cfg.WelcomeTopic, cfg.ReportTopic}
	if err := infra.CreateTopics(cfg.Brokers, topics); err != nil {
		log.Printf("Failed to create topics: %v", err)
	}
}

func startHttpServer(httpSrv *http.Server) {
	go func() {
		log.Printf("Starting HTTP server on port %s", httpSrv.Addr)

		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()
}

func startGrpcServer(grpcSrv *grpc.Server, port string) {
	go func() {
		lis, err := net.Listen("tcp", ":"+port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		log.Println("Starting gRPC server on :" + port)

		if err := grpcSrv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()
}

func shutdownServers(httpSrv *http.Server, grpcSrv *grpc.Server) {
	log.Println("Shutting down servers...")

	grpcSrv.GracefulStop()
	log.Println("gRPC server stopped.")

	shutdownCtx, stop := context.WithTimeout(context.Background(), 10*time.Second)
	defer stop()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown failed: %v", err)
	}

	log.Println("All servers shutdown gracefully.")
}
