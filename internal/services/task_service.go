package services

import (
	"context"
	"sort"
	"time"

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

func (ts *TaskService) GetAllTasks(ctx context.Context, status string) ([]*models.Task, error) {
	tasks, err := ts.repo.All(ctx)
	if err != nil {
		return nil, err
	}

	var filteredTasks []*models.Task

	for _, task := range tasks {
		if task.Status == status {
			filteredTasks = append(filteredTasks, task)
		}
	}

	sort.Slice(filteredTasks, func(i, j int) bool {
		return filteredTasks[i].ID < filteredTasks[j].ID
	})

	for _, task := range filteredTasks {
		activeDate, _ := time.Parse("2006-01-02", task.ActiveAt)
		if activeDate.Weekday() == time.Saturday || activeDate.Weekday() == time.Sunday {
			task.Title = "ВЫХОДНОЙ - " + task.Title
		}
	}

	return filteredTasks, nil
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

func (ts *TaskService) DoneTask(ctx context.Context, id string) error {
	return ts.repo.MarkAsDone(ctx, id)
}

func (ts *TaskService) DeleteTask(ctx context.Context, id string) error {
	return ts.repo.Delete(ctx, id)
}
