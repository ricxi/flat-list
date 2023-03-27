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

func TestCreateOne(t *testing.T) {
	r, teardown := setupRepo(t)
	defer teardown(t)

	task := createOneTask()
	taskID, err := r.CreateOne(context.Background(), &task)

	assert := assert.New(t)
	assert.NoError(err)
	if assert.NotEmpty(taskID) {
		if !primitive.IsValidObjectID(taskID) {
			t.Errorf("expected a hex that can be converted into a primitive.ObjectID")
		}
	}
}

func TestGetOne(t *testing.T) {
	r, teardown := setupRepo(t)
	defer teardown(t)

	t.Run("SuccessGetOneTask", func(t *testing.T) {
		task := createOneTask()
		taskID, err := r.CreateOne(context.Background(), &task)
		require.NoError(t, err)
		require.NotEmpty(t, taskID)

		actualTask, err := r.GetOne(context.Background(), taskID)

		assert := assert.New(t)
		assert.NoError(err)
		assert.NotNil(actualTask)
		if assert.NotEmpty(*actualTask) {
			assert.Equal(task.UserID, actualTask.UserID)
			assert.Equal(task.Name, actualTask.Name)
			assert.Equal(task.Details, actualTask.Details)
			assert.Equal(task.Priority, actualTask.Priority)
			assert.Equal(task.Category, actualTask.Category)
			assert.WithinDuration(*task.CreatedAt, *actualTask.CreatedAt, time.Second)
			assert.WithinDuration(*task.UpdatedAt, *actualTask.UpdatedAt, time.Second)
		}
	})

	t.Run("FailGetOneTask", func(t *testing.T) {
		taskID := primitive.NewObjectID().Hex()

		task, err := r.GetOne(context.Background(), taskID)
		assert.Error(t, err)
		assert.EqualError(t, err, mongo.ErrNoDocuments.Error())
		assert.Nil(t, task)
	})
}
