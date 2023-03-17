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
	vService := &user.ValidationService{Service: service}

	handler := user.NewHandler(vService)
	server := user.NewServer(handler, cfg.port)
	server.Run()
}

// build the service with its peripheral dependencies
func buildService(repository user.Repository) (user.Service, error) {
	passwordService := user.NewPasswordService(bcrypt.MinCost)
	grpcClient, err := user.NewMailerClient("grpc", "5000")
	if err != nil {
		return nil, err
	}

	return user.NewService(
		repository,
		passwordService,
		grpcClient,
	), nil
}
