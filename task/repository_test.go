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

// teardownFunc is returned by setupRepo
// to cleanup after the tests are finished running
type teardownFunc func(t testing.TB)

// setupRepo connects to mongo and creates and returns a
// repository for testing, and returns a teardown function.
// It requires a mongodb server to be running and the uri to access
// this server must be set with the environment variable 'MONGODB_URI'
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

func TestRepositoryCreateTask(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		r, teardown := setupRepo(t)
		defer teardown(t)
		assert := assert.New(t)

		nt := createNewTaskForRepo()
		gotTaskID, err := r.createTask(context.Background(), &nt)

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

	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		newTask := createNewTaskForRepo()
		taskID, err := r.createTask(context.Background(), &newTask)
		require.NoError(err)
		require.NotEmpty(taskID)
		require.True(primitive.IsValidObjectID(taskID))

		expectedTask := createExpectedTaskFromNew(taskID, newTask)
		actualTask, err := r.getTaskByID(context.Background(), taskID)
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

	t.Run("Fail", func(t *testing.T) {
		assert := assert.New(t)
		taskID := primitive.NewObjectID().Hex()

		task, err := r.getTaskByID(context.Background(), taskID)
		assert.Nil(task)
		if assert.Error(err) {
			assert.EqualError(err, ErrTaskNotFound.Error())
		}
	})
}

func TestRepositoryUpdateTask(t *testing.T) {
	r, teardown := setupRepo(t)
	defer teardown(t)

	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		newTask := createNewTaskForRepo()
		taskID, err := r.createTask(context.Background(), &newTask)
		require.NoError(err)
		if assert.NotEmpty(taskID) {
			assert.True(primitive.IsValidObjectID(taskID))
		}

		updatePayload := Task{
			ID:       taskID,
			Priority: "medium",
		}

		updatedTask, err := r.updateTask(context.Background(), &updatePayload)
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

		updatedTask, err := r.updateTask(context.Background(), &updatePayload)
		require.Nil(t, updatedTask)
		if assert.Error(err) {
			assert.EqualError(err, ErrTaskNotFound.Error())
		}
	})
}

func TestDeleteTaskByID(t *testing.T) {
	r, teardown := setupRepo(t)
	defer teardown(t)

	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		newTask := createNewTaskForRepo()
		taskID, err := r.createTask(context.Background(), &newTask)
		require.NoError(err)
		if assert.NotEmpty(taskID) {
			require.True(primitive.IsValidObjectID(taskID))
		}

		err = r.deleteTaskByID(context.Background(), taskID)
		assert.NoError(err)
	})

	t.Run("FailTaskNotFound", func(t *testing.T) {
		assert := assert.New(t)

		taskID := primitive.NewObjectID().Hex()
		err := r.deleteTaskByID(context.Background(), taskID)
		if assert.Error(err) {
			assert.EqualError(err, ErrTaskNotFound.Error())
		}
	})
}
