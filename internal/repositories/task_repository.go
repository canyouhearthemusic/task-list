package repositories

import (
	"context"
	"errors"
	"sync"

	"github.com/canyouhearthemusic/todo-list/internal/models"
)

type TaskRepo interface {
	GetByID(ctx context.Context, id string) (*models.Task, error)
	All(ctx context.Context) ([]*models.Task, error)
	Post(ctx context.Context, task *models.Task) error
	Put(ctx context.Context, id string, task *models.Task) error
	Delete(ctx context.Context, id string) error
}

type SyncMapTaskRepo struct {
	db sync.Map
}

func NewSyncMapTaskRepo() *SyncMapTaskRepo {
	return &SyncMapTaskRepo{}
}

func (repo *SyncMapTaskRepo) GetByID(ctx context.Context, id string) (*models.Task, error) {
	value, ok := repo.db.Load(id)
	if !ok {
		return nil, errors.New("Task not found")
	}

	task, ok := value.(*models.Task)
	if !ok {
		return nil, errors.New("Type assertion failed")
	}

	return task, nil
}

func (repo *SyncMapTaskRepo) All(ctx context.Context) ([]*models.Task, error) {
	var tasks []*models.Task

	repo.db.Range(func(key, value interface{}) bool {
		task, ok := value.(*models.Task)
		if ok {
			tasks = append(tasks, task)
		}

		return true
	})

	return tasks, nil
}

func (repo *SyncMapTaskRepo) Post(ctx context.Context, task *models.Task) error {
	var found bool
	repo.db.Range(func(key, value interface{}) bool {
		existingTask, ok := value.(*models.Task)

		if ok && existingTask.Title == task.Title {
			found = true
			return false
		}

		return true
	})

	if found {
		return errors.New("task with the same title already exists")
	}

	repo.db.Store(task.ID, task)

	return nil
}

func (repo *SyncMapTaskRepo) Put(ctx context.Context, id string, updatedTask *models.Task) error {
	oldTask, err := repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	updatedTask.ID = oldTask.ID

	repo.db.Store(id, updatedTask)

	return nil
}

func (repo *SyncMapTaskRepo) Delete(ctx context.Context, id string) error {
	_, err := repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	repo.db.Delete(id)

	return nil
}
