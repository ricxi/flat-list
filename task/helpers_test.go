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

// from json is used to convert json responses
// into a native go type. Ensure that a pointer
// is passed for out
func fromJSON(t testing.TB, r io.Reader, out any) {
	t.Helper()

	if err := json.NewDecoder(r).Decode(&out); err != nil {
		t.Fatal(err)
	}
}

// createNewTaskForRepo is a helper function used
// to create a new task for tests in the repository layer
func createNewTaskForRepo() NewTask {
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

// createExpectedTaskFromNew creates an expected task
// from a new task
func createExpectedTaskFromNew(id string, nt NewTask) Task {
	return Task{
		ID:        id,
		UserID:    nt.UserID,
		Name:      nt.Name,
		Details:   nt.Details,
		Priority:  nt.Priority,
		Category:  nt.Category,
		CreatedAt: nt.CreatedAt,
		UpdatedAt: nt.UpdatedAt,
	}
}

// createExpectedTask creates an expected task for tests
func createExpectedTask() Task {
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
