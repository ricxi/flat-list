package task

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ Repository = &mockRepository{}

type mockRepository struct {
	taskID string
	task   *Task
	err    error
}

func (m *mockRepository) CreateTask(ctx context.Context, task *NewTask) (string, error) {
	return m.taskID, m.err
}

func (m *mockRepository) GetTaskByID(ctx context.Context, id string) (*Task, error) {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return nil, err
	}
	return m.task, m.err
}

func (m *mockRepository) UpdateTask(ctx context.Context, task *Task) (*Task, error) {
	return m.task, m.err
}

func (m *mockRepository) DeleteTaskByID(ctx context.Context, id string) error {
	return m.err
}

var _ Service = &mockService{}

// mockService is used by the http handler
type mockService struct {
	taskID string
	task   *Task
	err    error
}

func (m *mockService) CreateTask(ctx context.Context, task *NewTask) (string, error) {
	return m.taskID, m.err
}

func (m *mockService) GetTaskByID(ctx context.Context, id string) (*Task, error) {
	return m.task, m.err
}

func (m *mockService) UpdateTask(ctx context.Context, task *Task) (*Task, error) {
	// m.task will be a mock of the updated task
	return m.task, m.err
}

func (m *mockService) DeleteTask(ctx context.Context, id string) error {
	return m.err
}
