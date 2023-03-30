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
	t.Run("GetTaskByIDFail", func(t *testing.T) {
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

	t.Run("GetTaskByIDFail", func(t *testing.T) {
		assert := assert.New(t)
		s := &service{
			r: &mockRepository{
				err:  ErrTaskNotFound,
				task: nil,
			},
		}
		taskID := primitive.NewObjectID().Hex()

		actualTask, err := s.GetTaskByID(context.Background(), taskID)
		require.Nil(t, actualTask)
		if assert.Error(err) {
			assert.EqualError(err, ErrTaskNotFound.Error())
		}
	})
}

// TestUpdateTask mainly tests the validation
// carried out by the service layer
func TestUpdateTask(t *testing.T) {
	t.Run("UpdateTaskNameSuccess", func(t *testing.T) {
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
			r: &mockRepository{
				err:  nil,
				task: &expectedTask,
			},
		}

		updatePayload := Task{
			ID:     expectedTask.ID,
			UserID: expectedTask.UserID,
			Name:   "Repair the laundry machine",
		}

		actualTask, err := s.UpdateTask(context.Background(), &updatePayload)
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

	t.Run("UpdateTaskFailMissingID", func(t *testing.T) {
		assert := assert.New(t)
		s := &service{
			r: &mockRepository{
				err:  nil,
				task: nil,
			},
		}

		updatePayload := Task{
			UserID: primitive.NewObjectID().Hex(),
			Name:   "Repair the laundry machine",
		}

		actualTask, err := s.UpdateTask(context.Background(), &updatePayload)
		assert.Nil(actualTask)
		if assert.Error(err) {
			assert.EqualError(err, ErrMissingField.Error()+": taskId")
		}
	})

	t.Run("UpdateTaskFailMissingUserID", func(t *testing.T) {
		assert := assert.New(t)
		s := &service{
			r: &mockRepository{
				err:  nil,
				task: nil,
			},
		}

		updatePayload := Task{
			ID:   primitive.NewObjectID().Hex(),
			Name: "Repair the laundry machine",
		}

		actualTask, err := s.UpdateTask(context.Background(), &updatePayload)
		assert.Nil(actualTask)
		if assert.Error(err) {
			assert.EqualError(err, ErrMissingField.Error()+": userId")
		}
	})

	t.Run("UpdateTaskFailMissingTaskName", func(t *testing.T) {
		assert := assert.New(t)
		s := &service{
			r: &mockRepository{
				err:  nil,
				task: nil,
			},
		}

		updatePayload := Task{
			ID:     primitive.NewObjectID().Hex(),
			UserID: primitive.NewObjectID().Hex(),
		}

		actualTask, err := s.UpdateTask(context.Background(), &updatePayload)
		assert.Nil(actualTask)
		if assert.Error(err) {
			assert.EqualError(err, ErrMissingField.Error()+": name")
		}
	})
}

func TestDeleteTask(t *testing.T) {
	t.Run("DeleteTaskSuccess", func(t *testing.T) {
		s := service{
			r: &mockRepository{
				err: nil,
			},
		}

		taskID := primitive.NewObjectID().Hex()
		err := s.DeleteTask(context.Background(), taskID)
		assert.NoError(t, err)
	})

	t.Run("DeleteTaskFailMissingFieldTaskID", func(t *testing.T) {
		// I didn't include an id because it never gets that far
		s := service{
			r: &mockRepository{
				err: nil,
			},
		}

		taskID := ""
		err := s.DeleteTask(context.Background(), taskID)
		foundErr := assert.Error(t, err)
		if foundErr {
			assert.EqualError(t, err, ErrMissingField.Error()+": taskId")
		}
	})

	t.Run("DeleteTaskFailTaskNotFound", func(t *testing.T) {
		assert := assert.New(t)
		s := service{
			r: &mockRepository{
				err: ErrTaskNotFound,
			},
		}

		taskID := primitive.NewObjectID().Hex()
		err := s.DeleteTask(context.Background(), taskID)
		foundErr := assert.Error(err)
		if foundErr {
			assert.EqualError(err, ErrTaskNotFound.Error())
		}
	})
}
