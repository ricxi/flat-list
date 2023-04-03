package main_test

import (
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
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var service task.Service

func TestMain(m *testing.M) {
	envs, err := config.LoadEnvs("MONGODB_URI")
	if err != nil {
		log.Fatal(err)
	}

	client, err := task.NewMongoClient(envs["MONGODB_URI"], 15)
	if err != nil {
		log.Fatalln("unable to connect to db", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Println("error occurred when disconnecting mongo client", err)
		}
	}()

	dbname := uuid.New().String()
	r := task.NewRepository(client, dbname)

	service = task.NewService(r)

	exitCode := m.Run()
	os.Exit(exitCode)
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
	"userId"   :"507f1f77bcf86cd799439011",
	"name"     :"laundry",
	"details"  :"quickly",
	"priority" :"high",
	"category" :"chores"
}`

func TestCreateTask(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	h := task.NewHTTPHandler(service)
	ts := httptest.NewTLSServer(h)
	defer ts.Close()

	createEndpoint := ts.URL + "/v1/task"
	body := strings.NewReader(createTaskPayloadStr)
	response, err := ts.Client().Post(createEndpoint, "application/json", body)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(http.StatusCreated, response.StatusCode)

	var actual createTaskResponse
	defer response.Body.Close()
	fromJSON(t, response.Body, &actual)

	if assert.NotEmpty(actual) {
		assert.True(primitive.IsValidObjectID(actual.TaskID))
		assert.True(actual.Success)
	}
}

func TestCreateThenGetTask(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	h := task.NewHTTPHandler(service)
	ts := httptest.NewTLSServer(h)
	defer ts.Close()

	// create a new task
	createEndpoint := ts.URL + "/v1/task"
	body := strings.NewReader(createTaskPayloadStr)
	response, err := ts.Client().Post(createEndpoint, "application/json", body)
	if err != nil {
		t.Fatal(err)
	}

	// if the task was successfully created, get its id
	// and use it to make a call and retreive it
	require.Equal(http.StatusCreated, response.StatusCode)

	var respBody createTaskResponse
	defer response.Body.Close()
	fromJSON(t, response.Body, &respBody)
	require.True(primitive.IsValidObjectID(respBody.TaskID))

	getEndpoint := ts.URL + "/v1/task/" + respBody.TaskID
	response, err = ts.Client().Get(getEndpoint)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(http.StatusOK, response.StatusCode)

	var getTaskResp getTaskResponse
	defer response.Body.Close()
	fromJSON(t, response.Body, &getTaskResp)

	expectedTask := Task{
		ID:       respBody.TaskID,
		UserID:   "507f1f77bcf86cd799439011",
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

func fromJSON(t testing.TB, r io.Reader, out any) {
	t.Helper()
	if err := json.NewDecoder(r).Decode(&out); err != nil {
		t.Fatal(err)
	}
}
