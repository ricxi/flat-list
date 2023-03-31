package task

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type teardownFunc func(t testing.TB)

// setupRepo connects to mongo and creates and returns a
// repository for testing, and returns a teardown function
func setupRepo(t testing.TB) (Repository, teardownFunc) {
	var (
		uri     string
		timeout int64
		client  *mongo.Client
		dbname  string
	)

	uri = os.Getenv("MONGODB_URI")
	if uri == "" {
		t.Fatal("missing uri")
	}

	timeout = 10

	client, err := NewMongoClient(uri, timeout)
	if err != nil {
		t.Fatal(err)
	}

	dbname = uuid.New().String()

	repository := NewRepository(client, dbname)

	return repository, func(t testing.TB) {
		if err := client.Database(dbname).Drop(context.Background()); err != nil {
			log.Println("unable to drop repo", err)
		}
		if err := client.Disconnect(context.Background()); err != nil {
			log.Println("Problem disconnecting from mongo", err)
		}
	}
}

// createOneTask is a helper function used to make a task
// for testing
func createOneTask() NewTask {
	// createdAt := time.Date(2023, time.March, 1, 2, 3, 4, 0, time.UTC)
	createdAt := time.Now().UTC()
	task := NewTask{
		UserID:    primitive.NewObjectID().Hex(),
		Name:      "Laundry",
		Details:   "tumble low and dry",
		Priority:  "low",
		Category:  "chores",
		CreatedAt: &createdAt,
		UpdatedAt: &createdAt,
	}

	return task
}

func TestRepositoryCreateTask(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		r, teardown := setupRepo(t)
		defer teardown(t)
		assert := assert.New(t)

		task := createOneTask()
		gotTaskID, err := r.CreateTask(context.Background(), &task)

		assert.NoError(err)
		if assert.NotEmpty(gotTaskID) {
			if !primitive.IsValidObjectID(gotTaskID) {
				t.Errorf("expected a hex that can be converted into a primitive.ObjectID")
			}
		}
	})
}

func TestRepositoryGetTaskByID(t *testing.T) {
	r, teardown := setupRepo(t)
	defer teardown(t)

	t.Run("SuccessGetTask", func(t *testing.T) {
		assert := assert.New(t)

		task := createOneTask()
		taskID, err := r.CreateTask(context.Background(), &task)
		require.NoError(t, err)
		require.NotEmpty(t, taskID)

		actualTask, err := r.GetTaskByID(context.Background(), taskID)
		assert.NoError(err)

		if assert.NotNil(actualTask) && assert.NotEmpty(*actualTask) {
			assert.Equal(taskID, actualTask.ID)
			assert.Equal(task.UserID, actualTask.UserID)
			assert.Equal(task.Name, actualTask.Name)
			assert.Equal(task.Details, actualTask.Details)
			assert.Equal(task.Priority, actualTask.Priority)
			assert.Equal(task.Category, actualTask.Category)
			assert.WithinDuration(*task.CreatedAt, *actualTask.CreatedAt, time.Second)
			assert.WithinDuration(*task.UpdatedAt, *actualTask.UpdatedAt, time.Second)
		}
	})

	t.Run("FailGetTask", func(t *testing.T) {
		assert := assert.New(t)
		taskID := primitive.NewObjectID().Hex()

		task, err := r.GetTaskByID(context.Background(), taskID)
		assert.Nil(task)
		if assert.Error(err) {
			assert.EqualError(err, ErrTaskNotFound.Error())
		}
	})
}

func TestRepositoryUpdateTask(t *testing.T) {
	r, teardown := setupRepo(t)
	defer teardown(t)
	t.Run("SuccessUpdateTask", func(t *testing.T) {
		assert := assert.New(t)
		newTask := createOneTask()
		taskID, err := r.CreateTask(context.Background(), &newTask)
		require.NoError(t, err)

		updatePayload := Task{
			ID:       taskID,
			Priority: "medium",
		}

		updatedTask, err := r.UpdateTask(context.Background(), &updatePayload)
		assert.NoError(err)

		expectedTask := newTask
		expectedTask.Priority = updatePayload.Priority
		if assert.NotNil(updatedTask) && assert.NotEmpty(*updatedTask) {
			assert.Equal(taskID, updatedTask.ID)
			assert.Equal(expectedTask.UserID, updatedTask.UserID)
			assert.Equal(expectedTask.Name, updatedTask.Name)
			assert.Equal(expectedTask.Details, updatedTask.Details)
			assert.Equal(expectedTask.Priority, updatedTask.Priority)
			assert.Equal(expectedTask.Category, updatedTask.Category)
			assert.WithinDuration(*expectedTask.CreatedAt, *updatedTask.CreatedAt, time.Second)
			assert.WithinDuration(*expectedTask.UpdatedAt, *updatedTask.UpdatedAt, time.Second)
		}
	})

	t.Run("FailUpdateTask", func(t *testing.T) {
		assert := assert.New(t)
		taskID := primitive.NewObjectID().Hex()
		updatePayload := Task{
			ID:       taskID,
			Priority: "medium",
		}

		updatedTask, err := r.UpdateTask(context.Background(), &updatePayload)
		require.Nil(t, updatedTask)
		if assert.Error(err) {
			assert.EqualError(err, ErrTaskNotFound.Error())
		}
	})
}

func TestDeleteTaskByID(t *testing.T) {
	r, teardown := setupRepo(t)
	defer teardown(t)

	t.Run("SuccessDeleteTaskByID", func(t *testing.T) {
		assert := assert.New(t)
		newTask := createOneTask()
		taskID, err := r.CreateTask(context.Background(), &newTask)
		require.NoError(t, err)

		err = r.DeleteTaskByID(context.Background(), taskID)
		assert.NoError(err)
	})

	t.Run("FailDeleteTaskDocumentNotFound", func(t *testing.T) {
		assert := assert.New(t)

		taskID := primitive.NewObjectID().Hex()
		err := r.DeleteTaskByID(context.Background(), taskID)
		if assert.Error(err) {
			assert.EqualError(err, ErrTaskNotFound.Error())
		}
	})
}
