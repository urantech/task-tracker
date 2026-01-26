package router

import (
	"backend-api/internal/auth"
	"backend-api/internal/task"
	"backend-api/internal/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	userHandler *user.Handler,
	authHandler *auth.Handler,
	authService *auth.Service,
	taskHandler *task.Handler,
) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	setupPublicRoutes(router, userHandler, authHandler)
	setupProtectedRoutes(router, authService, userHandler, taskHandler)

	return router
}

func setupPublicRoutes(
	router *chi.Mux,
	userHandler *user.Handler,
	authHandler *auth.Handler,
) {
	router.Group(func(r chi.Router) {
		r.Post("/users/register", userHandler.Register)
		r.Post("/auth/login", authHandler.Login)
	})
}

func setupProtectedRoutes(
	router *chi.Mux,
	authService *auth.Service,
	userHandler *user.Handler,
	taskHandler *task.Handler,
) {
	router.Group(func(r chi.Router) {
		r.Use(authService.AuthMiddleware)
		r.Get("/users/user", userHandler.GetCurrentUser)
		r.Post("/tasks", taskHandler.CreateTask)
		r.Get("/tasks", taskHandler.List)
		r.Patch("/tasks/{id}", taskHandler.UpdateTask)
	})
}
