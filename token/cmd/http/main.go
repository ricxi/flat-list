package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/julienschmidt/httprouter"
	"githug.com/ricxi/flat-list/token"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatalln("db connection env cannot be empty")
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		log.Fatalln("http port env cannot be empty")
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalln("unable to connect to postgres", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalln("unable to ping postgres", err)
	}

	repo := token.Repository{
		DB: db,
	}

	r := httprouter.New()

	r.POST("/v1/token/activation/:userID", token.HandlerCreateToken(&repo))
	r.GET("/v1/token/:userID", token.HandleGetTokens(&repo))

	srv := &http.Server{
		Handler: r,
		Addr:    ":" + httpPort,
	}

	log.Println("starting token service on port", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
