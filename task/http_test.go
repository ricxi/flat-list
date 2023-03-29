package task

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHandleCreateTask(t *testing.T) {
	assert := assert.New(t)
	expectedTaskID := primitive.NewObjectID().Hex()
	s := mockService{
		taskID: expectedTaskID,
		err:    nil,
	}
	h := &httpHandler{
		s: &s,
	}

	nt := NewTask{
		UserID:   primitive.NewObjectID().Hex(),
		Name:     "Laundry",
		Details:  "tumble low and dry",
		Priority: "low",
		Category: "chores",
	}

	w := httptest.NewRecorder()

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&nt); err != nil {
		t.Fatal(err)
	}

	// method and endpoint do not matter
	r, err := http.NewRequest(http.MethodPost, "/task/create", &body)
	require.NoError(t, err)

	h.handleCreateTask(w, r)

	result := w.Result()
	assert.Equal(result.StatusCode, http.StatusCreated)

	rBody := struct {
		Success bool   `json:"success"`
		TaskID  string `json:"taskId"`
	}{}
	if err := json.NewDecoder(result.Body).Decode(&rBody); err != nil {
		t.Fatal(err)
	}
	defer result.Body.Close()

	if assert.NotEmpty(rBody) {
		assert.Equal(expectedTaskID, rBody.TaskID)
		assert.Equal(true, rBody.Success)
	}
}
