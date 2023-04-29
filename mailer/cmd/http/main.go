package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ricxi/flat-list/mailer"
	"github.com/ricxi/flat-list/shared/config"
)

func main() {
	envs, err := config.LoadEnvs("HOST", "PORT", "USERNAME", "PASSWORD", "EMAIL_TEMPLATES", "HTTP_PORT")
	if err != nil {
		log.Fatal(err)
	}
	smtpPort, err := strconv.Atoi(envs["PORT"])
	if err != nil {
		log.Fatal(err)
	}

	m := mailer.NewMailer(envs["USERNAME"], envs["PASSWORD"], envs["HOST"], smtpPort)
	mailerService := mailer.NewService(m, envs["EMAIL_TEMPLATES"])

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/mailer/activate", mailer.HandleSendActivationEmail(mailerService))

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + envs["HTTP_PORT"],
	}

	log.Println("starting http mailer server on port", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
