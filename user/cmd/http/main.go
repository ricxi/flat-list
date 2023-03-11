package main

import (
	"log"

	"github.com/ricxi/flat-list/user"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg, err := getEnvs()
	if err != nil {
		log.Fatal(err)
	}

	client, err := newMongoClient(cfg.mongoURI, cfg.mongoTimeout)
	if err != nil {
		log.Fatal(err)
	}

	mongoRepository := user.NewMongoRepository(client, cfg.mongoDBName, cfg.mongoTimeout)
	passwordService := user.NewPasswordService(bcrypt.MinCost)
	service := user.NewService(mongoRepository, passwordService)
	vService := &user.ValidationService{Service: service}

	handler := user.NewHandler(vService)
	server := user.NewServer(handler, cfg.port)
	server.Run()
}
