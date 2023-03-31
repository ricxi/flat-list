package task

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func toJSON(t testing.TB, in any) io.Reader {
	t.Helper()

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(in); err != nil {
		t.Fatal(err)
	}

	return &body
}

// createOneTask is a helper function used to make a task
// for testing
func createNewTask() NewTask {
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

func createTaskForHTTPTests() Task {
	createdAt := time.Now().UTC()
	return Task{
		ID:        primitive.NewObjectID().Hex(),
		UserID:    primitive.NewObjectID().Hex(),
		Name:      "Repair the laundry machine",
		Details:   "tumble low and dry",
		Priority:  "low",
		Category:  "chores",
		CreatedAt: &createdAt,
		UpdatedAt: &createdAt,
	}
}
