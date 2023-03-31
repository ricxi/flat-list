package task

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHandleCreateTask(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		taskID := primitive.NewObjectID().Hex()
		expected := `{"success":true,"taskId":"` + taskID + `"}`
		h := NewHTTPHandler(
			&mockService{
				taskID: taskID,
				err:    nil,
			},
		)

		nt := NewTask{
			UserID:   primitive.NewObjectID().Hex(),
			Name:     "Laundry",
			Details:  "tumble low and dry",
			Priority: "low",
			Category: "chores",
		}

		rr := httptest.NewRecorder()

		body := toJSON(t, &nt)
		r, err := http.NewRequest(http.MethodPost, "/v1/task", body)
		require.NoError(err)

		h.ServeHTTP(rr, r)

		result := rr.Result()
		assert.Equal(http.StatusCreated, result.StatusCode)

		actual := strings.TrimSpace(rr.Body.String())
		if assert.NotEmpty(actual) {
			assert.Equal(expected, actual)
		}
	})
}

func TestHandleGetTask(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)
		expectedTask := createTaskForHTTPTests()
		h := NewHTTPHandler(
			&mockService{
				task: &expectedTask,
				err:  nil,
			},
		)

		rr := httptest.NewRecorder()

		r, err := http.NewRequest(http.MethodGet, "/v1/task/"+expectedTask.ID, nil)
		require.NoError(err)

		h.ServeHTTP(rr, r)

		result := rr.Result()
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

	t.Run("FailMissingUrlParams", func(t *testing.T) {
		assert := assert.New(t)
		h := httpHandler{
			s: &mockService{
				err: nil,
			},
		}

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
			assert.Equal("missing url param id", actualMessage)
		}
	})

	t.Run("FailTaskNotFound", func(t *testing.T) {
		assert := assert.New(t)
		h := httpHandler{
			s: &mockService{
				err: ErrTaskNotFound,
			},
		}

		w := httptest.NewRecorder()

		// url doesn't really matter
		r, err := http.NewRequest(http.MethodGet, "v1/task", nil)
		if err != nil {
			t.Fatal(err)
		}
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", primitive.NewObjectID().Hex())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
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
			assert.Equal(ErrTaskNotFound.Error(), actualMessage)
		}
	})
}

func TestHandleUpdateTask(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)
		// task does not matter
		task := createTaskForHTTPTests()
		h := httpHandler{
			s: &mockService{
				task: &task,
				err:  nil,
			},
		}

		w := httptest.NewRecorder()

		// nothing matters here except for the task
		body := toJSON(t, &task)
		r, err := http.NewRequest(http.MethodPut, "/v1/task", body)
		require.NoError(err)

		h.handleUpdateTask(w, r)

		result := w.Result()
		assert.Equal(http.StatusOK, result.StatusCode)
	})
}
