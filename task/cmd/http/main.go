package main

import (
	"log"

	"github.com/ricxi/flat-list/shared/config"
	"github.com/ricxi/flat-list/task"
)

func main() {
	envs, err := config.LoadEnvs("PORT", "MONGODB_URI", "MONGODB_NAME")
	if err != nil {
		log.Fatal(err)
	}

	client, err := task.NewMongoClient(envs["MONGODB_URI"], 15)
	if err != nil {
		log.Fatalln("unable to connect to db", err)
	}

	r := task.NewRepository(client, envs["MONGODB_NAME"])
	s := task.NewService(r)
	h := task.NewHTTPHandler(s)

	srv := task.Server{
		Port:    envs["PORT"],
		Handler: h,
	}

	if err := srv.Run(); err != nil {
		log.Println(err)
	}
}
