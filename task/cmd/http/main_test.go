// Package main_test contains some e2e tests.
// ! I don't think I'm doing this right.
package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ricxi/flat-list/shared/config"
	"github.com/ricxi/flat-list/task"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var service task.Service

// TestMain needs a live database connection
// and the connection string to that database set
// to the environment variable 'MONGODB_URI', or
// it won't run the tests.
func TestMain(m *testing.M) {
	envs, err := config.LoadEnvs("MONGODB_URI")
	if err != nil {
		log.Fatal(err)
	}

	// create a new database for testing
	dbname := uuid.New().String()

	client, err := task.NewMongoClient(envs["MONGODB_URI"], 15)
	if err != nil {
		log.Fatalln("unable to connect to db", err)
	}
	cleanup := func(exitCode int) int {
		// but what happens if it exits before os.Exit calls it?
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := client.Database(dbname).Drop(ctx); err != nil {
			log.Println("error dropping the database")
		}
		if err := client.Disconnect(ctx); err != nil {
			log.Println("error occurred when disconnecting mongo client", err)
		}
		return exitCode
	}

	repo := task.NewRepository(client, dbname)
	service = task.NewService(repo)

	exitCode := m.Run()
	os.Exit(cleanup(exitCode))
}

type Task struct {
	ID        string     `json:"taskId,omitempty"`
	UserID    string     `json:"userId,omitempty"`
	Name      string     `json:"name"`
	Details   string     `json:"details,omitempty"`
	Priority  string     `json:"priority,omitempty"`
	Category  string     `json:"category,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type createTaskResponse struct {
	Success bool   `json:"success"`
	TaskID  string `json:"taskId,omitempty"`
}

type getTaskResponse struct {
	Success bool `json:"success"`
	Task    Task `json:"task"`
}

const createTaskPayloadStr = `
{
	"name"     :"laundry",
	"details"  :"quickly",
	"priority" :"high",
	"category" :"chores"
}`

// mockUserServer mocks an instance of the user service
func mockUserServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"userId":"507f191e810c19729de860ea"}`))
	}))
}

// newRequestWithAuthHeaders creates a new
// http.Request and sets all the headers
// needed in order to make an api call that
// gets through the authentication middleware
func newRequestWithAuthHeaders(t testing.TB, method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer jsonwebtokengoeshere")

	return req
}

func TestCreateTask(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)

		mockUserService := mockUserServer()
		defer mockUserService.Close()

		h := task.NewHTTPHandler(
			service,
			(&task.Middleware{AuthEndpoint: mockUserService.URL}).Authenticate,
		)

		ts := httptest.NewServer(h)
		defer ts.Close()

		body := strings.NewReader(createTaskPayloadStr)
		req := newRequestWithAuthHeaders(t, http.MethodPost, ts.URL+"/v1/task", body)

		resp, err := ts.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(http.StatusCreated, resp.StatusCode)

		var actual createTaskResponse
		defer resp.Body.Close()
		fromJSON(t, resp.Body, &actual)

		// I think I can write a better test than this
		if assert.NotEmpty(actual) {
			// Do I really need to check something like that? Especially for these tests, which don't actually hit the user service?
			assert.True(primitive.IsValidObjectID(actual.TaskID))
			assert.True(actual.Success)
		}

	})

	t.Run("FailNoAuthHeader", func(t *testing.T) {
		assert := assert.New(t)

		expected := `{"error":"auth header is empty or missing"}`

		mockUserService := mockUserServer()
		defer mockUserService.Close()

		h := task.NewHTTPHandler(
			service,
			(&task.Middleware{AuthEndpoint: mockUserService.URL}).Authenticate,
		)

		ts := httptest.NewServer(h)
		defer ts.Close()

		body := strings.NewReader(createTaskPayloadStr)
		req, err := http.NewRequest(http.MethodPost, ts.URL+"/v1/task", body)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := ts.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(http.StatusUnauthorized, resp.StatusCode)

		defer resp.Body.Close()
		actual, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		if assert.NotEmpty(actual) {
			assert.JSONEq(expected, string(actual))
		}
	})
}

func TestCreateThenGetTask(t *testing.T) {
	assert := assert.New(t)

	mockUserService := mockUserServer()
	defer mockUserService.Close()

	h := task.NewHTTPHandler(
		service,
		(&task.Middleware{AuthEndpoint: mockUserService.URL}).Authenticate,
	)

	ts := httptest.NewServer(h)
	defer ts.Close()

	// create a new task
	body := strings.NewReader(createTaskPayloadStr)
	req := newRequestWithAuthHeaders(t, http.MethodPost, ts.URL+"/v1/task", body)

	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	// if the task was successfully created, get its id
	// and use it to make a call to get the newly created task
	assert.Equal(http.StatusCreated, resp.StatusCode)

	var respBody createTaskResponse
	defer resp.Body.Close()
	fromJSON(t, resp.Body, &respBody)
	assert.True(primitive.IsValidObjectID(respBody.TaskID))

	getEndpoint := ts.URL + "/v1/task/" + respBody.TaskID
	req = newRequestWithAuthHeaders(t, http.MethodGet, getEndpoint, nil)

	resp, err = ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(http.StatusOK, resp.StatusCode)

	var getTaskResp getTaskResponse
	defer resp.Body.Close()
	fromJSON(t, resp.Body, &getTaskResp)

	expectedTask := Task{
		ID:       respBody.TaskID,
		UserID:   "507f191e810c19729de860ea",
		Name:     "laundry",
		Details:  "quickly",
		Priority: "high",
		Category: "chores",
	}

	actualTask := getTaskResp.Task
	if assert.NotEmpty(getTaskResp) {
		assert.Equal(expectedTask.ID, actualTask.ID)
		assert.Equal(expectedTask.UserID, actualTask.UserID)
		assert.Equal(expectedTask.Name, actualTask.Name)
		assert.Equal(expectedTask.Details, actualTask.Details)
		assert.Equal(expectedTask.Priority, actualTask.Priority)
		assert.Equal(expectedTask.Category, actualTask.Category)
	}
}
func TestCreateGetUpdateTask(t *testing.T) {
	assert := assert.New(t)

	mockUserService := mockUserServer()
	defer mockUserService.Close()

	h := task.NewHTTPHandler(
		service,
		(&task.Middleware{AuthEndpoint: mockUserService.URL}).Authenticate,
	)

	ts := httptest.NewServer(h)
	defer ts.Close()

	// create a new task
	body := strings.NewReader(createTaskPayloadStr)
	req := newRequestWithAuthHeaders(t, http.MethodPost, ts.URL+"/v1/task", body)

	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	// if the task was successfully created, get its id
	// and use it to make a call and retreive it
	assert.Equal(http.StatusCreated, resp.StatusCode)

	var respBody createTaskResponse
	defer resp.Body.Close()
	fromJSON(t, resp.Body, &respBody)
	newTaskID := respBody.TaskID

	assert.True(primitive.IsValidObjectID(newTaskID))

	getEndpoint := ts.URL + "/v1/task/" + newTaskID
	req = newRequestWithAuthHeaders(t, http.MethodGet, getEndpoint, nil)

	resp, err = ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(http.StatusOK, resp.StatusCode)

	var getTaskResp getTaskResponse
	defer resp.Body.Close()
	fromJSON(t, resp.Body, &getTaskResp)

	expectedTask := Task{
		ID:       newTaskID,
		UserID:   "507f191e810c19729de860ea",
		Name:     "laundry",
		Details:  "quickly",
		Priority: "high",
		Category: "chores",
	}

	actualTask := getTaskResp.Task
	if assert.NotEmpty(getTaskResp) {
		assert.Equal(expectedTask.ID, actualTask.ID)
		assert.Equal(expectedTask.UserID, actualTask.UserID)
		assert.Equal(expectedTask.Name, actualTask.Name)
		assert.Equal(expectedTask.Details, actualTask.Details)
		assert.Equal(expectedTask.Priority, actualTask.Priority)
		assert.Equal(expectedTask.Category, actualTask.Category)
	}

	expectedTask.Name = "dishes"
	expectedTask.Details = "slowly"
	expectedTask.Priority = "low"

	updatePayload := expectedTask

	var buffer bytes.Buffer
	toJSON(t, &buffer, &updatePayload)
	bodyReader := bytes.NewReader(buffer.Bytes())

	req = newRequestWithAuthHeaders(t, http.MethodPut, ts.URL+"/v1/task/", bodyReader)
	resp, err = ts.Client().Do(req)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(http.StatusOK, resp.StatusCode)
	// ! Not finished: I still have to assert the actual update response body returned
}

// fromJSON is a helper function that decodes a response body
// into a native go type (which must be a pointer)
func fromJSON(t testing.TB, r io.Reader, out any) {
	t.Helper()
	if err := json.NewDecoder(r).Decode(&out); err != nil {
		t.Fatal(err)
	}
}

// toJSON is a helper function
func toJSON(t testing.TB, w io.Writer, in any) {
	t.Helper()
	if err := json.NewEncoder(w).Encode(&in); err != nil {
		t.Fatal(err)
	}
}
