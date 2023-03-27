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

func (m *mockRepository) CreateOne(ctx context.Context, task *NewTask) (string, error) {
	return m.taskID, m.err
}

func (m *mockRepository) GetOne(ctx context.Context, id string) (*Task, error) {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return nil, err
	}
	return m.task, m.err
}

func (m *mockRepository) UpdateOne(ctx context.Context, task *Task) (*Task, error) {
	return m.task, m.err
}
