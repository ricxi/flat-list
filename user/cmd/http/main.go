package main

import (
	"context"
	"log"
	"strconv"

	"github.com/ricxi/flat-list/shared/config"
	"github.com/ricxi/flat-list/user"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	envs, err := config.LoadEnvs("PORT", "MONGODB_URI", "MONGODB_NAME", "MONGODB_TIMEOUT")
	if err != nil {
		log.Fatal(err)
	}
	mongoTimeout, err := strconv.Atoi(envs["MONGODB_TIMEOUT"])
	if err != nil {
		log.Fatal(err)
	}

	client, err := user.NewMongoClient(envs["MONGODB_URI"], mongoTimeout)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	mongoRepository := user.NewRepository(client, envs["MONGODB_NAME"], mongoTimeout)
	service, err := buildService(mongoRepository)
	if err != nil {
		log.Fatal(err)
	}

	handler := user.NewHandler(service)
	server := user.NewServer(handler, envs["PORT"])
	server.Run()
}

// build the service with its peripheral dependencies
func buildService(repository user.Repository) (user.Service, error) {
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
