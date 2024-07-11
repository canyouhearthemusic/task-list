package routes

import (
	_ "github.com/canyouhearthemusic/todo-list/docs"
	"github.com/canyouhearthemusic/todo-list/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Todo List API
// @version 1.0
// @description This is a simple Todo List API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath

func New() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.URLFormat)

	loadRoutes(r)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	return r
}

func loadRoutes(r *chi.Mux) {
	r.Route("/api", func(api chi.Router) {
		api.Route("/tasks", func(tasks chi.Router) {
			tasks.Get("/", handlers.GetAllTasks)
			tasks.Post("/", handlers.PostTask)
			tasks.Get("/{id}", handlers.GetTask)
			tasks.Put("/{id}", handlers.PutTask)
			tasks.Put("/{id}/done", handlers.DoneTask)
			tasks.Delete("/{id}", handlers.DeleteTask)
		})
	})
}
