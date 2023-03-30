package task

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
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
	r, err := http.NewRequest(http.MethodPost, "v1/task/create", &body)
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

func TestHandleGetTask(t *testing.T) {
	t.Run("GetTaskSuccess", func(t *testing.T) {
		expectedTask := createTaskForHTTPTests()
		h := httpHandler{
			s: &mockService{
				task: &expectedTask,
				err:  nil,
			},
		}
		assert := assert.New(t)

		w := httptest.NewRecorder()

		// url doesn't really matter
		r, err := http.NewRequest(http.MethodGet, "v1/task", nil)
		if err != nil {
			t.Fatal(err)
		}
		// manually add url params to request context to avoid missing url param error
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", expectedTask.ID)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		h.handleGetTask(w, r)

		result := w.Result()

		assert.Equal(http.StatusOK, result.StatusCode)

		resBody := struct {
			Success bool `json:"success"`
			Task    Task `json:"task"`
		}{}
		if err := json.NewDecoder(result.Body).Decode(&resBody); err != nil {
			t.Fatal(err)
		}
		defer result.Body.Close()

		assert.Equal(true, resBody.Success)

		actualTask := resBody.Task
		if assert.NotNil(&actualTask) && assert.NotEmpty(actualTask) {
			assert.Equal(expectedTask.ID, actualTask.ID)
			assert.Equal(expectedTask.Name, actualTask.Name)
			assert.Equal(expectedTask.UserID, actualTask.UserID)
			assert.Equal(expectedTask.Details, actualTask.Details)
			assert.Equal(expectedTask.Category, actualTask.Category)
			assert.Equal(expectedTask.Priority, actualTask.Priority)
			assert.WithinDuration(*expectedTask.CreatedAt, *actualTask.CreatedAt, time.Second)
			assert.WithinDuration(*expectedTask.UpdatedAt, *actualTask.UpdatedAt, time.Second)
		}
	})

	t.Run("GetTaskFailMissingUrlParams", func(t *testing.T) {
		h := httpHandler{
			s: &mockService{
				// I didn't include task because it shouldn't get that far
				err: nil,
			},
		}
		expectedMessage := "missing url param id"
		assert := assert.New(t)

		w := httptest.NewRecorder()

		// url doesn't really matter
		r, err := http.NewRequest(http.MethodGet, "v1/task", nil)
		if err != nil {
			t.Fatal(err)
		}
		h.handleGetTask(w, r)

		result := w.Result()

		assert.Equal(http.StatusBadRequest, result.StatusCode)

		resBody := struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{}
		if err := json.NewDecoder(result.Body).Decode(&resBody); err != nil {
			t.Fatal(err)
		}
		defer result.Body.Close()

		assert.Equal(false, resBody.Success)

		actualMessage := resBody.Message
		if assert.NotEmpty(actualMessage) {
			assert.Equal(expectedMessage, actualMessage)
		}
	})

	t.Run("GetTaskFailMissingUrlParams", func(t *testing.T) {
		h := httpHandler{
			s: &mockService{
				// I didn't include task because it shouldn't get that far
				err: nil,
			},
		}
		expectedMessage := "missing url param id"
		assert := assert.New(t)

		w := httptest.NewRecorder()

		// url doesn't really matter
		r, err := http.NewRequest(http.MethodGet, "v1/task", nil)
		if err != nil {
			t.Fatal(err)
		}
		h.handleGetTask(w, r)

		result := w.Result()

		assert.Equal(http.StatusBadRequest, result.StatusCode)

		resBody := struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{}
		if err := json.NewDecoder(result.Body).Decode(&resBody); err != nil {
			t.Fatal(err)
		}
		defer result.Body.Close()

		assert.Equal(false, resBody.Success)

		actualMessage := resBody.Message
		if assert.NotEmpty(actualMessage) {
			assert.Equal(expectedMessage, actualMessage)
		}
	})
}
