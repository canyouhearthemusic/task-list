package models

import (
	"errors"
	"time"
)

type Task struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	ActiveAt string `json:"activeAt"`
	Status   string `json:"status"`
}

type TaskRequest struct {
	Title    string `json:"title"`
	ActiveAt string `json:"activeAt"`
}

func (t *Task) Validate() error {
	if t.ID != "" {
		return errors.New("id mustn't present in request")
	}

	if len(t.Title) > 200 {
		return errors.New("title exceeds 200 characters")
	}

	if _, err := time.Parse("2006-01-02", t.ActiveAt); err != nil {
		return errors.New("invalid activeAt format")
	}

	return nil
}
