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

// GetAllTasks godoc
// @Summary Get all tasks
// @Description Get all tasks by status
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param   status  query  string  false  "Status Filter"  Enum(active,done)
// @Success 200 {array} models.Task
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/tasks [get]
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

// GetTask godoc
// @Summary Get task by ID
// @Description Get a task by ID
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param   id   path  string  true  "Task ID"
// @Success 200 {object} models.Task
// @Failure 404 {string} string "Not Found"
// @Router /api/tasks/{id} [get]
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

// PostTask godoc
// @Summary Create a new task
// @Description Create a new task
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param   task  body  models.TaskRequest  true  "Task"
// @Success 201 {object} models.Task
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/tasks [post]
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

// PutTask godoc
// @Summary Update a task
// @Description Update a task by ID
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param   id    path  string       true  "Task ID"
// @Param   task  body  models.TaskRequest  true  "Task"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Router /api/tasks/{id} [put]
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

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete a task by ID
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "Task ID"
// @Success 204 {string} string "No Content"
// @Failure 404 {string} string "Not Found"
// @Router /api/tasks/{id} [delete]
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := service.DeleteTask(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DoneTask godoc
// @Summary Mark task as done
// @Description Update a task status by ID
// @Tags tasks
// @Param   id    path  string       true  "Task ID"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Router /api/tasks/{id}/done [put]
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
