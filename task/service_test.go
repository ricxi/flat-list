package task

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestServiceCreateTask(t *testing.T) {
	t.Run("SuccessCreateTask", func(t *testing.T) {
		assert := assert.New(t)

		s := &service{
			repository: &mockRepository{
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

		gotTaskID, err := s.createTask(context.Background(), &task)
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
			repository: &mockRepository{
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

		gotTaskID, err := s.createTask(context.Background(), &task)
		assert.Empty(gotTaskID)
		if assert.Error(err) {
			assert.EqualError(err, fmt.Errorf("%w: name", ErrMissingField).Error())
		}
	})

	t.Run("FailCreateTaskMissingUserIDField", func(t *testing.T) {
		assert := assert.New(t)
		s := &service{
			repository: &mockRepository{
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

		gotTaskID, err := s.createTask(context.Background(), &task)
		assert.Empty(gotTaskID)
		if assert.Error(err) {
			assert.EqualError(err, fmt.Errorf("%w: userId", ErrMissingField).Error())
		}
	})
}

func TestServiceGetTaskByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)

		createdAt := time.Now().UTC()
		task := Task{
			ID:        primitive.NewObjectID().Hex(),
			UserID:    primitive.NewObjectID().Hex(),
			Name:      "Laundry",
			Details:   "tumble low and dry",
			Priority:  "low",
			Category:  "chores",
			CreatedAt: &createdAt,
			UpdatedAt: &createdAt,
		}

		s := &service{
			repository: &mockRepository{
				err:  nil,
				task: &task,
			},
		}

		actualTask, err := s.getTaskByID(context.Background(), task.ID)
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

	t.Run("FailTaskNotFound", func(t *testing.T) {
		// Do not create a task
		assert := assert.New(t)
		s := &service{
			repository: &mockRepository{
				err:  ErrTaskNotFound,
				task: nil,
			},
		}
		taskID := primitive.NewObjectID().Hex()

		actualTask, err := s.getTaskByID(context.Background(), taskID)
		require.Nil(t, actualTask)
		if assert.Error(err) {
			assert.EqualError(err, ErrTaskNotFound.Error())
		}
	})
}

// TestUpdateTask mainly tests the validation
// carried out by the service layer
func TestServiceUpdateTask(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)

		createdAt := time.Now().UTC()
		expectedTask := Task{
			ID:        primitive.NewObjectID().Hex(),
			UserID:    primitive.NewObjectID().Hex(),
			Name:      "Repair the laundry machine",
			Details:   "tumble low and dry",
			Priority:  "low",
			Category:  "chores",
			CreatedAt: &createdAt,
			UpdatedAt: &createdAt,
		}

		s := &service{
			repository: &mockRepository{
				err:  nil,
				task: &expectedTask,
			},
		}

		updatePayload := Task{
			ID:     expectedTask.ID,
			UserID: expectedTask.UserID,
			Name:   "Repair the laundry machine",
		}

		actualTask, err := s.updateTask(context.Background(), &updatePayload)
		assert.NoError(err)
		if assert.NotNil(actualTask) && assert.NotEmpty(*actualTask) {
			assert.Equal(expectedTask.ID, actualTask.ID)
			assert.Equal(expectedTask.UserID, actualTask.UserID)
			assert.Equal(expectedTask.Name, actualTask.Name)
			assert.Equal(expectedTask.Details, actualTask.Details)
			assert.Equal(expectedTask.Priority, actualTask.Priority)
			assert.Equal(expectedTask.Category, actualTask.Category)
			assert.WithinDuration(*expectedTask.CreatedAt, *actualTask.CreatedAt, time.Second)
			assert.WithinDuration(*expectedTask.UpdatedAt, *actualTask.UpdatedAt, time.Second)
		}

	})

	t.Run("FailMissingID", func(t *testing.T) {
		assert := assert.New(t)
		s := &service{
			repository: &mockRepository{
				err:  nil,
				task: nil,
			},
		}

		updatePayload := Task{
			UserID: primitive.NewObjectID().Hex(),
			Name:   "Repair the laundry machine",
		}

		actualTask, err := s.updateTask(context.Background(), &updatePayload)
		assert.Nil(actualTask)
		if assert.Error(err) {
			assert.EqualError(err, ErrMissingField.Error()+": taskId")
		}
	})

	t.Run("FailMissingUserID", func(t *testing.T) {
		assert := assert.New(t)
		s := &service{
			repository: &mockRepository{
				err:  nil,
				task: nil,
			},
		}

		updatePayload := Task{
			ID:   primitive.NewObjectID().Hex(),
			Name: "Repair the laundry machine",
		}

		actualTask, err := s.updateTask(context.Background(), &updatePayload)
		assert.Nil(actualTask)
		if assert.Error(err) {
			assert.EqualError(err, ErrMissingField.Error()+": userId")
		}
	})

	t.Run("FailMissingTaskName", func(t *testing.T) {
		assert := assert.New(t)
		s := &service{
			repository: &mockRepository{
				err:  nil,
				task: nil,
			},
		}

		updatePayload := Task{
			ID:     primitive.NewObjectID().Hex(),
			UserID: primitive.NewObjectID().Hex(),
		}

		actualTask, err := s.updateTask(context.Background(), &updatePayload)
		assert.Nil(actualTask)
		if assert.Error(err) {
			assert.EqualError(err, ErrMissingField.Error()+": name")
		}
	})
}

func TestServiceDeleteTask(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		s := service{
			repository: &mockRepository{
				err: nil,
			},
		}

		taskID := primitive.NewObjectID().Hex()
		err := s.deleteTask(context.Background(), taskID)
		assert.NoError(t, err)
	})

	t.Run("FailMissingFieldTaskID", func(t *testing.T) {
		// I didn't include an id because it never gets that far
		s := service{
			repository: &mockRepository{
				err: nil,
			},
		}

		taskID := ""
		err := s.deleteTask(context.Background(), taskID)
		foundErr := assert.Error(t, err)
		if foundErr {
			assert.EqualError(t, err, ErrMissingField.Error()+": taskId")
		}
	})

	t.Run("FailTaskNotFound", func(t *testing.T) {
		assert := assert.New(t)
		s := service{
			repository: &mockRepository{
				err: ErrTaskNotFound,
			},
		}

		taskID := primitive.NewObjectID().Hex()
		err := s.deleteTask(context.Background(), taskID)
		foundErr := assert.Error(err)
		if foundErr {
			assert.EqualError(err, ErrTaskNotFound.Error())
		}
	})
}
