package task

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateTask(t *testing.T) {
	t.Run("SucccessCreateTask", func(t *testing.T) {
		assert := assert.New(t)
		s := &service{
			r: &mockRepository{
				taskID: primitive.NewObjectID().Hex(),
				err:    nil,
			},
		}

		task := NewTask{
			UserID:   primitive.NewObjectID().Hex(),
			Name:     "Laundry",
			Details:  "tumble low and dry",
			Priority: "low",
			Category: "chores",
		}

		gotTaskID, err := s.CreateTask(context.Background(), &task)
		assert.NoError(err)
		if assert.NotEmpty(gotTaskID) {
			if !primitive.IsValidObjectID(gotTaskID) {
				t.Errorf("expected a hex that can be converted into a primitive.ObjectID, but did not get one")
			}
		}
	})

	t.Run("FailCreateTaskMissingNameField", func(t *testing.T) {
		assert := assert.New(t)
		s := &service{
			r: &mockRepository{
				taskID: "",
				err:    ErrMissingField,
			},
		}

		task := NewTask{
			UserID:   primitive.NewObjectID().Hex(),
			Details:  "tumble low and dry",
			Priority: "low",
			Category: "chores",
		}

		gotTaskID, err := s.CreateTask(context.Background(), &task)
		assert.Empty(gotTaskID)
		if assert.Error(err) {
			assert.EqualError(err, fmt.Errorf("%w: name", ErrMissingField).Error())
		}
	})

	t.Run("FailCreateTaskMissingUserIDField", func(t *testing.T) {
		assert := assert.New(t)
		s := &service{
			r: &mockRepository{
				taskID: "",
				err:    ErrMissingField,
			},
		}

		task := NewTask{
			Name:     "Laundry",
			Details:  "tumble low and dry",
			Priority: "low",
			Category: "chores",
		}

		gotTaskID, err := s.CreateTask(context.Background(), &task)
		assert.Empty(gotTaskID)
		if assert.Error(err) {
			assert.EqualError(err, fmt.Errorf("%w: userId", ErrMissingField).Error())
		}
	})
}

func TestGetTaskByID(t *testing.T) {
	t.Run("SuccessGetTask", func(t *testing.T) {
		assert := assert.New(t)

		task := Task{
			ID:       primitive.NewObjectID().Hex(),
			UserID:   primitive.NewObjectID().Hex(),
			Name:     "Laundry",
			Details:  "tumble low and dry",
			Priority: "low",
			Category: "chores",
		}

		s := &service{
			r: &mockRepository{
				err:  nil,
				task: &task,
			},
		}

		actualTask, err := s.GetTaskByID(context.Background(), task.ID)
		assert.NoError(err)
		if assert.NotNil(actualTask) && assert.NotEmpty(*actualTask) {
			assert.Equal(task.ID, actualTask.ID)
			assert.Equal(task.UserID, actualTask.UserID)
			assert.Equal(task.Name, actualTask.Name)
			assert.Equal(task.Details, actualTask.Details)
			assert.Equal(task.Priority, actualTask.Priority)
			assert.Equal(task.Category, actualTask.Category)
			assert.WithinDuration(*task.CreatedAt, *actualTask.CreatedAt, time.Second)
			assert.WithinDuration(*task.UpdatedAt, *actualTask.UpdatedAt, time.Second)
		}
	})

}
