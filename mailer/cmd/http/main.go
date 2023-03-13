package main

import (
	"log"
	"net/http"

	"github.com/ricxi/flat-list/mailer"
)

func main() {
	conf, err := mailer.SetupConfig()
	if err != nil {
		log.Fatal(err)
	}

	m := mailer.NewMailer(conf.Username, conf.Password, conf.Host, conf.Port)

	es := mailer.NewEmailService(m)

	mux := http.NewServeMux()

	// sends activation email
	mux.HandleFunc("/v1/mailer/activate", mailer.HandleSendEmail(es))

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + conf.HttpPort,
	}

	log.Println("starting http server on port", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
