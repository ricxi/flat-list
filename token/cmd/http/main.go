package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ricxi/flat-list/shared/config"
	"github.com/ricxi/flat-list/token"
)

func main() {
	envs, err := config.LoadEnvs("DATABASE_URL", "HTTP_PORT")
	if err != nil {
		log.Fatalln("problem loading configuation: ", err)
	}

	db, err := token.Connect(envs["DATABASE_URL"])
	if err != nil {
		log.Fatalln("problem connecting to postgres: ", err)
	}
	defer db.Close()

	repo := token.NewRepository(db)
	h := token.NewHTTPHandler(repo)

	srv := &http.Server{
		Handler:      h,
		Addr:         ":" + envs["HTTP_PORT"],
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Println("starting token service on port", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
