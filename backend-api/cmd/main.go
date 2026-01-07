package main

import (
	"backend-api/internal/config"
	"backend-api/internal/user"
	"backend-api/migrations"
	"backend-api/pkg/postgres"
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.MustLoad()

	log.Printf("Starting server on port: %s", cfg.Port)

	db, err := postgres.NewConnection(cfg.DbUrl)
	if err != nil {
		log.Fatalf("Failed connect to database: %v", err)
	}

	goose.SetBaseFS(migrations.EmbedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	log.Println("Running migrations...")

	if err := goose.Up(db, "."); err != nil {
		log.Fatal(err)
	}

	log.Println("Migrations completed successfully")

	userStorage := user.NewStorage(db)
	userService := user.NewService(userStorage)
	userHandler := user.NewHandler(userService)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Mount("/users", userHandler.Routes())

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	startSever(srv)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	shutdownServer(srv)
}

func startSever(srv *http.Server) {
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()
}

func shutdownServer(srv *http.Server) {
	log.Println("Shutting down server...")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("Server shutdown failed: %v", err)
		panic(err)
	}

	log.Println("Server shutdown gracefully.")
}
