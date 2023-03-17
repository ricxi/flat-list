package main

import (
	"context"
	"log"

	"github.com/ricxi/flat-list/user"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg, err := getEnvs()
	if err != nil {
		log.Fatal(err)
	}

	client, err := user.NewMongoClient(cfg.mongoURI, cfg.mongoTimeout)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	mongoRepository := user.NewMongoRepository(client, cfg.mongoDBName, cfg.mongoTimeout)
	service, err := buildService(mongoRepository)
	if err != nil {
		log.Fatal(err)
	}

	handler := user.NewHandler(service)
	server := user.NewServer(handler, cfg.port)
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

	return user.NewService(
		repository,
		grpcClient,
		passwordManager,
		validator,
	), nil
}
