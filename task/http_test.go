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

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"userId":"507f191e810c19729de860ea"}`))
		}))
		defer ts.Close()

		taskID := primitive.NewObjectID().Hex()
		expected := fmt.Sprintf(`{"success":true,"taskId":"%s"}`, taskID)

		h := NewHTTPHandler(
			&mockService{
				taskID: taskID,
				err:    nil,
			},
			(&middleware{authEndpoint: ts.URL}).authenticate,
		)

		rr := httptest.NewRecorder()

		newTask := `
		{
			"name"     :"laundry",
			"details"  :"quickly",
			"priority" :"high",
			"category" :"chores"
		}`
		reqBody := strings.NewReader(newTask)

		r := httptest.NewRequest(http.MethodPost, "/v1/task", reqBody)
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", "Bearer jwtsignedtokengoeshere")

		h.ServeHTTP(rr, r)

		assert.Equal(http.StatusCreated, rr.Code)
		if assert.NotEmpty(rr.Body) {
			assert.JSONEq(expected, rr.Body.String())
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

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"userId":"507f191e810c19729de860ea"}`))
		}))
		defer ts.Close()

		expTask := createExpectedTask()
		h := NewHTTPHandler(
			&mockService{
				task: &expTask,
				err:  nil,
			},
			(&middleware{authEndpoint: ts.URL}).authenticate,
		)

		rr := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/v1/task/"+expTask.ID, nil)
		r.Header.Set("Authorization", "Bearer signedjsonwebtokengoeshere")

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

	t.Run("FailTaskNotFound", func(t *testing.T) {
		assert := assert.New(t)

		expected := `{"message":"task not found","success":false}`

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"userId":"507f191e810c19729de860ea"}`))
		}))
		defer ts.Close()

		h := NewHTTPHandler(
			&mockService{
				err: ErrTaskNotFound,
			},
			(&middleware{authEndpoint: ts.URL}).authenticate,
		)

		rr := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/v1/task/"+primitive.NewObjectID().Hex(), nil)
		r.Header.Set("Authorization", "Bearer signedjwttokengoeshere")

		h.ServeHTTP(rr, r)

		assert.Equal(http.StatusBadRequest, rr.Code)
		if assert.NotEmpty(rr.Body) {
			assert.JSONEq(expected, rr.Body.String())
		}
	})
}

func TestHandleUpdateTask(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)

		newTask := createExpectedTask()
		expTask := newTask

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"userId":"507f191e810c19729de860ea"}`))
		}))
		defer ts.Close()

		h := NewHTTPHandler(
			&mockService{
				task: &newTask,
				err:  nil,
			},
			(&middleware{authEndpoint: ts.URL}).authenticate,
		)

		rr := httptest.NewRecorder()

		// nothing matters here except for the task
		reqBody := toJSON(t, &expTask)
		r := httptest.NewRequest(http.MethodPut, "/v1/task", reqBody)
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", "Bearer signedjwttokengoeshere")

		h.ServeHTTP(rr, r)

		res := rr.Result()
		assert.Equal(http.StatusOK, res.StatusCode)

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

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"userId":"507f191e810c19729de860ea"}`))
		}))
		defer ts.Close()

		h := NewHTTPHandler(
			&mockService{
				err: fmt.Errorf("%w: taskId", ErrMissingField),
			},
			(&middleware{authEndpoint: ts.URL}).authenticate,
		)

		rr := httptest.NewRecorder()

		placeholderBody := toJSON(t, struct{}{})
		r := httptest.NewRequest(http.MethodPut, "/v1/task", placeholderBody)
		r.Header.Set("Authorization", "Bearer signedjsonwebtoken")
		r.Header.Set("Content-Type", "application/json")

		h.ServeHTTP(rr, r)

		require.Equal(http.StatusBadRequest, rr.Code)

		if assert.NotEmpty(rr.Body) {
			assert.JSONEq(expResBody, rr.Body.String())
		}
	})

	t.Run("FailMissingUserIDField", func(t *testing.T) {
		assert := assert.New(t)

		expected := `{"message":"missing field is required: userId","success":false}`

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"userId":"507f191e810c19729de860ea"}`))
		}))
		defer ts.Close()

		h := NewHTTPHandler(
			&mockService{
				err: fmt.Errorf("%w: userId", ErrMissingField),
			},
			(&middleware{authEndpoint: ts.URL}).authenticate,
		)

		rr := httptest.NewRecorder()

		placeholderBody := toJSON(t, struct{}{})
		r := httptest.NewRequest(http.MethodPut, "/v1/task", placeholderBody)
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", "Bearer signedjsonwebtokengoeshere")

		h.ServeHTTP(rr, r)

		assert.Equal(http.StatusBadRequest, rr.Code)

		if assert.NotEmpty(rr.Body) {
			assert.JSONEq(expected, rr.Body.String())
		}
	})

	t.Run("FailMissingNameField", func(t *testing.T) {
		assert := assert.New(t)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"userId":"507f191e810c19729de860ea"}`))
		}))
		defer ts.Close()

		expected := `{"message":"missing field is required: name","success":false}`
		h := NewHTTPHandler(
			&mockService{
				err: fmt.Errorf("%w: name", ErrMissingField),
			},
			(&middleware{authEndpoint: ts.URL}).authenticate,
		)

		rr := httptest.NewRecorder()

		placeholderBody := toJSON(t, struct{}{})
		r := httptest.NewRequest(http.MethodPut, "/v1/task", placeholderBody)
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", "Bearer signedjsonwebtokengoeshere")

		h.ServeHTTP(rr, r)

		assert.Equal(http.StatusBadRequest, rr.Code)

		if assert.NotEmpty(rr.Body) {
			assert.JSONEq(expected, rr.Body.String())
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
