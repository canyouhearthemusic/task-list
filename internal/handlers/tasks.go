package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/canyouhearthemusic/todo-list/internal/models"
	"github.com/canyouhearthemusic/todo-list/internal/repositories"
	"github.com/canyouhearthemusic/todo-list/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var (
	service *services.TaskService = services.New(repositories.NewSyncMapTaskRepo())
)

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	status := r.URL.Query().Get("status")
	switch status {
		case "", "active":
			status = "active"
		case "done":
			status = "done"
	}

	tasks, err := service.GetAllTasks(ctx, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	task, err := service.GetTask(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func PostTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var task models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := task.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task.ID = uuid.New().String()

	if err := service.PostTask(ctx, &task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, task)
}

func PutTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var updatedTask models.Task

	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := updatedTask.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := service.PutTask(ctx, id, &updatedTask); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := service.DeleteTask(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DoneTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := service.DoneTask(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func respondWithJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
