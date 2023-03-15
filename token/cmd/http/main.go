package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/julienschmidt/httprouter"
	"githug.com/ricxi/flat-list/token"
)

func main() {
	config, err := token.LoadConfig()
	if err != nil {
		log.Fatalln("problem loading configuation: ", err)
	}

	db, err := token.Connect(config.DatabaseURL)
	if err != nil {
		log.Fatalln("problem connecting to postgres: ", err)
	}
	defer db.Close()

	repo := token.NewRepository(db)

	r := httprouter.New()
	r.POST("/v1/token/activation/:userID", token.HandlerCreateToken(&repo))
	r.GET("/v1/token/:userID", token.HandleGetTokens(&repo))

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + config.HttpPort,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Println("starting token service on port", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
