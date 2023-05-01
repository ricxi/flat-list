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
	envs, err := config.LoadEnvs(
		"PORT",
		"MONGODB_URI",
		"MONGODB_NAME",
		"MONGODB_TIMEOUT",
		"MAILER_GRPC_PORT",
		"TOKEN_GRPC_PORT",
	)
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

	mr := user.NewRepository(client, envs["MONGODB_NAME"])
	v := user.NewValidator()
	pm := user.NewPasswordManager(bcrypt.MinCost)

	mc, err := user.NewGRPCMailerClient(envs["MAILER_GRPC_PORT"])
	if err != nil {
		log.Fatal(err)
	}

	tc, err := user.NewTokenClient(envs["TOKEN_GRPC_PORT"])
	if err != nil {
		log.Fatalln(err)
	}

	service := user.NewService(
		mr,
		user.WithValidator(v),
		user.WithPasswordManager(pm),
		user.WithMailerClient(mc),
		user.WithTokenClient(tc),
	)

	handler := user.NewHTTPHandler(service)
	server := user.NewServer(handler, envs["PORT"])

	server.Run()
}
