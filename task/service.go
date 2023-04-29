package task

import (
	"context"
	"fmt"
	"time"
)

type Service interface {
	createTask(ctx context.Context, task *NewTask) (string, error)
	getTaskByID(ctx context.Context, id string) (*Task, error)
	updateTask(ctx context.Context, task *Task) (*Task, error)
	deleteTask(ctx context.Context, id string) error
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
func (s *service) createTask(ctx context.Context, task *NewTask) (string, error) {
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
	return s.repository.createTask(ctx, task)
}

func (s *service) getTaskByID(ctx context.Context, id string) (*Task, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: taskId", ErrMissingField)
	}

	return s.repository.getTaskByID(ctx, id)
}

func (s *service) updateTask(ctx context.Context, task *Task) (*Task, error) {
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

	return s.repository.updateTask(ctx, task)
}

func (s *service) deleteTask(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("%w: taskId", ErrMissingField)
	}

	return s.repository.deleteTaskByID(ctx, id)
}
