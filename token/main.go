package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/julienschmidt/httprouter"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalln("unable to connect to postgres", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalln("unable to ping postgres", err)
	}

	repo := repository{
		db: db,
	}

	r := httprouter.New()
	r.POST("/v1/token/activation/:userID", handlerCreateToken(&repo))
	r.GET("/v1/token/:userID", handleGetTokens(&repo))

	httpPort := os.Getenv("HTTP_PORT")

	srv := &http.Server{
		Handler: r,
		Addr:    ":" + httpPort,
	}

	log.Println("starting token service on port", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
