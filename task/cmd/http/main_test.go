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

const createTaskPayload = `
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
	body := strings.NewReader(createTaskPayload)
	response, err := ts.Client().Post(createEndpoint, "application/json", body)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(http.StatusCreated, response.StatusCode)

	actualBody := struct {
		Success bool   `json:"success"`
		TaskID  string `json:"taskId"`
	}{}
	defer response.Body.Close()
	fromJSON(t, response.Body, &actualBody)

	if assert.NotEmpty(actualBody) {
		assert.True(primitive.IsValidObjectID(actualBody.TaskID))
		assert.True(actualBody.Success)
	}
}

func TestCreateThenGetTask(t *testing.T) {
	require := require.New(t)
	// assert := assert.New(t)
	h := task.NewHTTPHandler(service)
	ts := httptest.NewTLSServer(h)
	defer ts.Close()

	createEndpoint := ts.URL + "/v1/task"
	body := strings.NewReader(createTaskPayload)
	response, err := ts.Client().Post(createEndpoint, "application/json", body)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(http.StatusCreated, response.StatusCode)

	respBody := struct {
		Success bool   `json:"success"`
		TaskID  string `json:"taskId"`
	}{}
	defer response.Body.Close()
	fromJSON(t, response.Body, &respBody)

	getEndpoint := ts.URL + "/v1/task/" + respBody.TaskID
	response, err = ts.Client().Get(getEndpoint)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(http.StatusOK, response.StatusCode)
}

func fromJSON(t testing.TB, r io.Reader, out any) {
	t.Helper()
	if err := json.NewDecoder(r).Decode(&out); err != nil {
		t.Fatal(err)
	}
}
