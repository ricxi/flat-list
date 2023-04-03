package task

import (
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
		expResBody := fmt.Sprintf(`{"success":true,"taskId":"%s"}`, taskID)

		h := NewHTTPHandler(
			&mockService{
				taskID: taskID,
				err:    nil,
			},
		)

		rr := httptest.NewRecorder()

		newTask := `
		{
			"userId"   :"507f1f77bcf86cd799439011",
			"name"     :"laundry",
			"details"  :"quickly",
			"priority" :"high",
			"category" :"chores"
		}`
		reqBody := strings.NewReader(newTask)

		r, err := http.NewRequest(http.MethodPost, "/v1/task", reqBody)
		require.NoError(err)

		h.ServeHTTP(rr, r)

		require.Equal(http.StatusCreated, rr.Code)

		actResBody := rr.Body.String()
		if assert.NotEmpty(actResBody) {
			assert.JSONEq(expResBody, actResBody)
		}
	})
}

type ResponseBody struct {
	Success bool `json:"success"`
	Task    Task `json:"task"`
}

func TestHandleGetTask(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		expTask := createExpectedTask()
		h := NewHTTPHandler(
			&mockService{
				task: &expTask,
				err:  nil,
			},
		)

		rr := httptest.NewRecorder()

		r, err := http.NewRequest(http.MethodGet, "/v1/task/"+expTask.ID, nil)
		require.NoError(err)

		h.ServeHTTP(rr, r)

		res := rr.Result()
		assert.Equal(http.StatusOK, res.StatusCode)

		var resBody ResponseBody
		defer res.Body.Close()
		fromJSON(t, res.Body, &resBody)

		assert.Equal(true, resBody.Success)

		actualTask := resBody.Task
		if assert.NotNil(&actualTask) && assert.NotEmpty(actualTask) {
			assert.Equal(expTask.ID, actualTask.ID)
			assert.Equal(expTask.Name, actualTask.Name)
			assert.Equal(expTask.UserID, actualTask.UserID)
			assert.Equal(expTask.Details, actualTask.Details)
			assert.Equal(expTask.Category, actualTask.Category)
			assert.Equal(expTask.Priority, actualTask.Priority)
			assert.WithinDuration(*expTask.CreatedAt, *actualTask.CreatedAt, time.Second)
			assert.WithinDuration(*expTask.UpdatedAt, *actualTask.UpdatedAt, time.Second)
		}
	})

	t.Run("FailMissingUrlParams", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		expResBody := `{"message":"missing url param id","success":false}`

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

		actRespBody := rr.Body.String()
		if assert.NotEmpty(actRespBody) {
			assert.JSONEq(expResBody, actRespBody)
		}
	})

	t.Run("FailTaskNotFound", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		expResBody := `{"message":"task not found","success":false}`
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

		actResBody := rr.Body.String()
		if assert.NotEmpty(actResBody) {
			assert.JSONEq(expResBody, actResBody)
		}
	})
}

func TestHandleUpdateTask(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		expTask := createExpectedTask()

		h := NewHTTPHandler(
			&mockService{
				task: &expTask,
				err:  nil,
			},
		)

		rr := httptest.NewRecorder()

		// nothing matters here except for the task
		reqBody := toJSON(t, &expTask)
		r, err := http.NewRequest(http.MethodPut, "/v1/task", reqBody)
		require.NoError(err)

		h.ServeHTTP(rr, r)

		res := rr.Result()
		require.Equal(http.StatusOK, res.StatusCode)

		var resBody ResponseBody
		defer res.Body.Close()
		fromJSON(t, res.Body, &resBody)

		assert.Equal(true, resBody.Success)

		actualTask := resBody.Task
		if assert.NotEmpty(resBody) {
			assert.Equal(expTask.ID, actualTask.ID)
			assert.Equal(expTask.Name, actualTask.Name)
			assert.Equal(expTask.UserID, actualTask.UserID)
			assert.Equal(expTask.Details, actualTask.Details)
			assert.Equal(expTask.Category, actualTask.Category)
			assert.Equal(expTask.Priority, actualTask.Priority)
			assert.WithinDuration(*expTask.CreatedAt, *actualTask.CreatedAt, time.Second)
			assert.WithinDuration(*expTask.UpdatedAt, *actualTask.UpdatedAt, time.Second)
		}
	})

	t.Run("FailMissingIDField", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		expResBody := `{"message":"missing field is required: taskId","success":false}`

		h := NewHTTPHandler(
			&mockService{
				err: fmt.Errorf("%w: taskId", ErrMissingField),
			},
		)

		rr := httptest.NewRecorder()

		body := toJSON(t, struct{}{})
		r, err := http.NewRequest(http.MethodPut, "/v1/task", body)
		require.NoError(err)

		h.ServeHTTP(rr, r)

		require.Equal(http.StatusBadRequest, rr.Code)

		actResBody := rr.Body.String()
		if assert.NotEmpty(actResBody) {
			assert.JSONEq(expResBody, actResBody)
		}
	})

	t.Run("FailMissingUserIDField", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		expResBody := `{"message":"missing field is required: userId","success":false}`
		h := NewHTTPHandler(
			&mockService{
				err: fmt.Errorf("%w: userId", ErrMissingField),
			},
		)

		rr := httptest.NewRecorder()

		placeholderBody := toJSON(t, struct{}{})
		r, err := http.NewRequest(http.MethodPut, "/v1/task", placeholderBody)
		require.NoError(err)

		h.ServeHTTP(rr, r)

		require.Equal(http.StatusBadRequest, rr.Code)

		actResBody := rr.Body.String()
		if assert.NotEmpty(actResBody) {
			assert.JSONEq(expResBody, actResBody)
		}
	})

	t.Run("FailMissingNameField", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		expResBody := `{"message":"missing field is required: name","success":false}`
		h := NewHTTPHandler(
			&mockService{
				err: fmt.Errorf("%w: name", ErrMissingField),
			},
		)

		rr := httptest.NewRecorder()

		placeholderBody := toJSON(t, struct{}{})
		r, err := http.NewRequest(http.MethodPut, "/v1/task", placeholderBody)
		require.NoError(err)

		h.ServeHTTP(rr, r)

		require.Equal(http.StatusBadRequest, rr.Code)

		actualResBody := rr.Body.String()
		if assert.NotEmpty(actualResBody) {
			assert.JSONEq(expResBody, actualResBody)
		}
	})
}

func TestHandleDeleteTask(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		expResBody := `{"success":true}`
		h := NewHTTPHandler(
			&mockService{
				err: nil,
			},
		)

		rr := httptest.NewRecorder()

		r, err := http.NewRequest(http.MethodDelete, "/v1/task/"+primitive.NewObjectID().Hex(), nil)
		require.NoError(err)

		h.ServeHTTP(rr, r)

		assert.Equal(http.StatusOK, rr.Code)

		actResBody := rr.Body.String()
		if assert.NotEmpty(actResBody) {
			assert.JSONEq(expResBody, actResBody)
		}
	})
}
