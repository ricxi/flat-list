package task

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// newRequestWithHeaders creates an *http.Request from the httptest package, and adds some headers to it.
func newRequestWithHeaders(method, target string, body io.Reader, headers map[string]string) *http.Request {
	r := httptest.NewRequest(method, target, body)
	for key, value := range headers {
		r.Header.Set(key, value)
	}
	return r
}

// newRequestWithJSONHeader creates an *http.Request from the httptest package,
// and adds a 'Content-Type: application/json' header to it.
func newRequestWithJSONHeader(method, target string, body io.Reader) *http.Request {
	return newRequestWithHeaders(
		method,
		target,
		body,
		map[string]string{"Content-Type": "application/json"},
	)
}

func TestHandleCreateTask(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		type expected struct {
			statusCode int
			hasBody    bool
			bodyStr    string
		}
		testCases := []struct {
			name           string
			mockService    Service
			testServer     *httptest.Server
			initMiddleware func() ([]func(http.Handler) http.Handler, func())
			request        *http.Request
			expected       expected
		}{
			{
				name: "Success",
				mockService: &mockService{
					taskID: "6067f0c53c56d02bf8e8dc74",
				},
				testServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(`{"userId":"507f191e810c19729de860ea"}`))
				})),
				initMiddleware: func() ([]func(http.Handler) http.Handler, func()) {
					// not sure if this is really awful and complicated, but I'll keep it here just in case
					ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte(`{"userId":"507f191e810c19729de860ea"}`))
					}))

					auth := (&Middleware{AuthEndpoint: ts.URL}).Authenticate

					return []func(http.Handler) http.Handler{auth}, func() {
						ts.Close()
					}

				},
				request: newRequestWithHeaders(
					http.MethodPost,
					"/v1/task",
					strings.NewReader(`
					{
						"name"     :"laundry",
						"details"  :"quickly",
						"priority" :"high",
						"category" :"chores"
					}
					`),
					map[string]string{
						"Content-Type":  "application/json",
						"Authorization": "Bearer jwtsignedtokengoeshere",
					}),
				expected: expected{
					statusCode: 201,
					bodyStr:    `{"success": true, "taskId": "6067f0c53c56d02bf8e8dc74"}`,
				},
			},
			{
				name: "MissingFieldName",
				mockService: &mockService{
					err: fmt.Errorf("%w: name", ErrMissingField),
				},
				testServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(`{"userId":"507f191e810c19729de860ea"}`))
				})),
				request: newRequestWithHeaders(
					http.MethodPost,
					"/v1/task",
					strings.NewReader(`
					{
						"details"  :"quickly",
						"priority" :"high",
						"category" :"chores"
					}
					`),
					map[string]string{
						"Content-Type":  "application/json",
						"Authorization": "Bearer jwtsignedtokengoeshere",
					}),
				expected: expected{
					statusCode: 400,
					bodyStr:    `{"success": false, "error": "missing field is required: name"}`,
				},
			},
		}
		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				assert := assert.New(t)
				defer tt.testServer.Close()

				// middlewares, cleanup := tt.initMiddleware()
				// defer cleanup()

				h := NewHTTPHandler(
					tt.mockService,
					(&Middleware{AuthEndpoint: tt.testServer.URL}).Authenticate,
					// middlewares...,
				)

				rr := httptest.NewRecorder()

				h.ServeHTTP(rr, tt.request)

				assert.Equal(tt.expected.statusCode, rr.Code)
				if assert.NotEmpty(rr.Body) {
					assert.JSONEq(tt.expected.bodyStr, rr.Body.String())
				}
			})
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
			(&Middleware{AuthEndpoint: ts.URL}).Authenticate,
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

		expected := `{"error":"task not found","success":false}`

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"userId":"507f191e810c19729de860ea"}`))
		}))
		defer ts.Close()

		h := NewHTTPHandler(
			&mockService{
				err: ErrTaskNotFound,
			},
			(&Middleware{AuthEndpoint: ts.URL}).Authenticate,
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
			(&Middleware{AuthEndpoint: ts.URL}).Authenticate,
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

		expResBody := `{"error":"missing field is required: taskId","success":false}`

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"userId":"507f191e810c19729de860ea"}`))
		}))
		defer ts.Close()

		h := NewHTTPHandler(
			&mockService{
				err: fmt.Errorf("%w: taskId", ErrMissingField),
			},
			(&Middleware{AuthEndpoint: ts.URL}).Authenticate,
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

		expected := `{"error":"missing field is required: userId","success":false}`

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"userId":"507f191e810c19729de860ea"}`))
		}))
		defer ts.Close()

		h := NewHTTPHandler(
			&mockService{
				err: fmt.Errorf("%w: userId", ErrMissingField),
			},
			(&Middleware{AuthEndpoint: ts.URL}).Authenticate,
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

		expected := `{"error":"missing field is required: name","success":false}`
		h := NewHTTPHandler(
			&mockService{
				err: fmt.Errorf("%w: name", ErrMissingField),
			},
			(&Middleware{AuthEndpoint: ts.URL}).Authenticate,
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

		expected := `{"success":true}`

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"userId":"507f191e810c19729de860ea"}`))
		}))
		defer ts.Close()

		h := NewHTTPHandler(
			&mockService{
				err: nil,
			},
			(&Middleware{AuthEndpoint: ts.URL}).Authenticate,
		)

		rr := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodDelete, "/v1/task/"+primitive.NewObjectID().Hex(), nil)
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", "Bearer jsonwebtokengoeshere")

		h.ServeHTTP(rr, r)

		assert.Equal(http.StatusOK, rr.Code)

		if assert.NotEmpty(rr.Body) {
			assert.JSONEq(expected, rr.Body.String())
		}
	})
}
