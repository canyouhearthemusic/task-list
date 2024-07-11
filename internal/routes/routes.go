package routes

import (
	"github.com/canyouhearthemusic/todo-list/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.URLFormat)

	loadRoutes(r)

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
