package task

import (
	"context"
	"fmt"
	"time"
)

type Service interface {
	CreateTask(ctx context.Context, task *NewTask) (string, error)
	GetTaskByID(ctx context.Context, id string) (*Task, error)
	UpdateTask(ctx context.Context, task *Task) (*Task, error)
	DeleteTask(ctx context.Context, id string) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{
		repository: repository,
	}
}

// CreateTask returns the task's id if a new task is successfully
// created and inserted into the database; otherwise it returns an error
func (s *service) CreateTask(ctx context.Context, task *NewTask) (string, error) {
	if task.Name == "" {
		return "", fmt.Errorf("%w: name", ErrMissingField)
	}

	if task.UserID == "" {
		return "", fmt.Errorf("%w: userId", ErrMissingField)
	}

	createdAt := time.Now().UTC()
	task.CreatedAt = &createdAt
	task.UpdatedAt = &createdAt

	// include a logger
	return s.repository.CreateTask(ctx, task)
}

func (s *service) GetTaskByID(ctx context.Context, id string) (*Task, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: taskId", ErrMissingField)
	}

	return s.repository.GetTaskByID(ctx, id)
}

func (s *service) UpdateTask(ctx context.Context, task *Task) (*Task, error) {
	if task.ID == "" {
		return nil, fmt.Errorf("%w: taskId", ErrMissingField)
	}

	if task.UserID == "" {
		return nil, fmt.Errorf("%w: userId", ErrMissingField)
	}

	if task.Name == "" {
		return nil, fmt.Errorf("%w: name", ErrMissingField)
	}

	updatedAt := time.Now().UTC()
	task.UpdatedAt = &updatedAt

	return s.repository.UpdateTask(ctx, task)
}

func (s *service) DeleteTask(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("%w: taskId", ErrMissingField)
	}

	return s.repository.DeleteTaskByID(ctx, id)
}
