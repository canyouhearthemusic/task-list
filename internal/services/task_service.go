package services

import (
	"context"

	"github.com/canyouhearthemusic/todo-list/internal/models"
	"github.com/canyouhearthemusic/todo-list/internal/repositories"
)

type TaskService struct {
	repo repositories.TaskRepo
}

func New(repo repositories.TaskRepo) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (ts *TaskService) GetAllTasks(ctx context.Context) ([]*models.Task, error) {
	return ts.repo.All(ctx)
}

func (ts *TaskService) GetTask(ctx context.Context, id string) (*models.Task, error) {
	return ts.repo.GetByID(ctx, id)
}

func (ts *TaskService) PostTask(ctx context.Context, task *models.Task) error {
	return ts.repo.Post(ctx, task)
}

func (ts *TaskService) PutTask(ctx context.Context, id string, task *models.Task) error {
	return ts.repo.Put(ctx, id, task)
}

func (ts *TaskService) DeleteTask(ctx context.Context, id string) error {
	return ts.repo.Delete(ctx, id)
}
