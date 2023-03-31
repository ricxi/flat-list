package task

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

		assert.Equal(http.StatusCreated, rr.Code)

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
		expectedTask := createExpectedTask()
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
		require := require.New(t)

		expected := `{"message":"missing url param id","success":false}`
		h := NewHTTPHandler(
			&mockService{
				err: nil,
			},
		)

		rr := httptest.NewRecorder()

		r, err := http.NewRequest(http.MethodGet, "/v1/task", nil)
		require.NoError(err)

		h.ServeHTTP(rr, r)

		require.Equal(http.StatusBadRequest, rr.Code)

		actual := strings.TrimSpace(rr.Body.String())
		if assert.NotEmpty(actual) {
			assert.Equal(expected, actual)
		}
	})

	t.Run("FailTaskNotFound", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		expected := `{"message":"task not found","success":false}`
		h := NewHTTPHandler(
			&mockService{
				err: ErrTaskNotFound,
			},
		)

		rr := httptest.NewRecorder()

		r, err := http.NewRequest(http.MethodGet, "/v1/task/"+primitive.NewObjectID().Hex(), nil)
		require.NoError(err)

		h.ServeHTTP(rr, r)

		assert.Equal(http.StatusBadRequest, rr.Code)

		actual := strings.TrimSpace(rr.Body.String())
		if assert.NotEmpty(actual) {
			assert.Equal(expected, actual)
		}
	})
}

func TestHandleUpdateTask(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		expectedTask := createExpectedTask()
		h := NewHTTPHandler(
			&mockService{
				task: &expectedTask,
				err:  nil,
			},
		)

		rr := httptest.NewRecorder()

		// nothing matters here except for the task
		body := toJSON(t, &expectedTask)
		r, err := http.NewRequest(http.MethodPut, "/v1/task", body)
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
		if assert.NotEmpty(resBody) {
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

	t.Run("FailMissingIDField", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		expected := `{"message":"missing field is required: taskId","success":false}`
		h := NewHTTPHandler(
			&mockService{
				err: fmt.Errorf("%w: taskId", ErrMissingField),
			},
		)

		rr := httptest.NewRecorder()

		// task is just placeholder to avoid nil pointer dereference by decoder?
		body := toJSON(t, struct{}{})
		r, err := http.NewRequest(http.MethodPut, "/v1/task", body)
		require.NoError(err)

		h.ServeHTTP(rr, r)

		assert.Equal(http.StatusBadRequest, rr.Code)

		actual := strings.TrimSpace(rr.Body.String())
		if assert.NotEmpty(actual) {
			assert.Equal(expected, actual)
		}
	})

	t.Run("FailMissingUserIDField", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		expected := `{"message":"missing field is required: userId","success":false}`
		h := NewHTTPHandler(
			&mockService{
				err: fmt.Errorf("%w: userId", ErrMissingField),
			},
		)

		rr := httptest.NewRecorder()

		body := toJSON(t, struct{}{})
		r, err := http.NewRequest(http.MethodPut, "/v1/task", body)
		require.NoError(err)

		h.ServeHTTP(rr, r)

		assert.Equal(http.StatusBadRequest, rr.Code)

		actual := strings.TrimSpace(rr.Body.String())
		if assert.NotEmpty(actual) {
			assert.Equal(expected, actual)
		}
	})

	t.Run("FailMissingNameField", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		expected := `{"message":"missing field is required: name","success":false}`
		h := NewHTTPHandler(
			&mockService{
				err: fmt.Errorf("%w: name", ErrMissingField),
			},
		)

		rr := httptest.NewRecorder()

		body := toJSON(t, struct{}{})
		r, err := http.NewRequest(http.MethodPut, "/v1/task", body)
		require.NoError(err)

		h.ServeHTTP(rr, r)

		assert.Equal(http.StatusBadRequest, rr.Code)

		actual := strings.TrimSpace(rr.Body.String())
		if assert.NotEmpty(actual) {
			assert.Equal(expected, actual)
		}
	})
}
