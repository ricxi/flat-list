package task

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestCreateOne(t *testing.T) {
	r, teardown := setupRepo(t)
	defer teardown(t)

	createdAt := time.Date(2023, time.March, 1, 2, 3, 4, 0, time.UTC)
	task := NewTask{
		UserID:    primitive.NewObjectID().Hex(),
		Name:      "Laundry",
		Details:   "tumble low and dry",
		Priority:  "low",
		Category:  "chores",
		CreatedAt: &createdAt,
		UpdatedAt: &createdAt,
	}

	taskID, err := r.CreateOne(context.Background(), &task)
	require.NoError(t, err)
	require.NotEmpty(t, taskID)

	if _, err := primitive.ObjectIDFromHex(taskID); err != nil {
		t.Errorf("expected a hex that can be converted into a primitive.ObjectID")
	}
}

func TestGetOne(t *testing.T) {
	r, teardown := setupRepo(t)
	defer teardown(t)

	createdAt := time.Date(2023, time.March, 1, 2, 3, 4, 0, time.UTC)
	task := NewTask{
		UserID:    primitive.NewObjectID().Hex(),
		Name:      "Laundry",
		Details:   "tumble low and dry",
		Priority:  "low",
		Category:  "chores",
		CreatedAt: &createdAt,
		UpdatedAt: &createdAt,
	}

	taskID, err := r.CreateOne(context.Background(), &task)
	require.NoError(t, err)
	require.NotEmpty(t, taskID)

	actualTask, err := r.GetOne(context.Background(), taskID)
	require.NoError(t, err)

	require.Equal(t, task.UserID, actualTask.UserID)
	require.Equal(t, task.Name, actualTask.Name)
	require.Equal(t, task.Details, actualTask.Details)
	require.Equal(t, task.Priority, actualTask.Priority)
	require.Equal(t, task.Category, actualTask.Category)
	require.Equal(t, *task.CreatedAt, *actualTask.CreatedAt)
	require.Equal(t, *task.UpdatedAt, *actualTask.UpdatedAt)
}
