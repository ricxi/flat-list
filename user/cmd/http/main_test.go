package main_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/ricxi/flat-list/shared/config"
	"github.com/ricxi/flat-list/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

var service user.Service

func TestMain(m *testing.M) {
	envs, err := config.LoadEnvs("MONGODB_URI")
	if err != nil {
		log.Fatal(err)
	}

	client, err := user.NewMongoClient(envs["MONGODB_URI"], 15)
	if err != nil {
		log.Fatalln("unable to connect to db", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Println("error occurred when disconnecting mongo client", err)
		}
	}()

	dbname := uuid.New().String()
	r := user.NewRepository(client, dbname)

	service, err = setupService(r)
	if err != nil {
		log.Fatal(err)
	}

	exitCode := m.Run()
	// I think this is stopping the database from being cleaned up
	// because it exits the program before the deferred functions are called.
	os.Exit(exitCode)
}

func setupService(repository user.Repository) (user.Service, error) {
	passwordManager := user.NewPasswordManager(bcrypt.MinCost)
	validator := user.NewValidator()
	grpcClient, err := user.NewMailerClient("grpc", "5001")
	if err != nil {
		return nil, err
	}

	tokenClient, err := user.NewTokenClient("5003")
	if err != nil {
		log.Fatalln(err)
	}

	service := user.
		NewServiceBuilder().
		Repository(repository).
		PasswordManager(passwordManager).
		MailerClient(grpcClient).
		TokenClient(tokenClient).
		Validator(validator).
		Build()

	return service, nil
}

const registerUserPayload string = `
{
    "firstName": "Michael",
    "lastName": "Scott",
    "email": "michaelscott@dundermifflin.com",
    "password": "1234"
}
`

func TestRegisterUser(t *testing.T) {
	t.Skip()
	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)
		h := user.NewHTTPHandler(service)
		ts := httptest.NewTLSServer(h)
		defer ts.Close()

		endpoint := ts.URL + "/v1/user/register"
		body := strings.NewReader(registerUserPayload)
		resp, err := ts.Client().Post(endpoint, "application/json", body)
		require.NoError(err)

		defer resp.Body.Close()
		var u user.UserInfo
		fromJSON(t, resp.Body, &u)

		if assert.NotEmpty(u) {
			fmt.Println(u)
		}

	})

}

func fromJSON(t testing.TB, r io.Reader, out any) {
	t.Helper()
	if err := json.NewDecoder(r).Decode(&out); err != nil {
		t.Fatal(err)
	}
}
