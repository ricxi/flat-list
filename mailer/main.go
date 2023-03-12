package main

import (
	"log"
	"net/http"
)

func main() {
	conf, err := setupConfig()
	if err != nil {
		log.Fatal(err)
	}

	m := NewMailer(conf.username, conf.password, conf.host, conf.port)

	emailService := &Service{
		tmplFilename: "./useractivation.html",
		mailer:       m,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/send", handleSendEmail(emailService))

	srv := &http.Server{
		Handler: mux,
		Addr:    ":5000",
	}

	log.Println("starting server on port", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
