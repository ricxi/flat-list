package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	conf, err := setupConfig()
	if err != nil {
		log.Fatal(err)
	}

	m := NewMailer(conf.username, conf.password, conf.host, conf.port)

	emailService := &EmailService{
		mailer: m,
	}

	mux := http.NewServeMux()

	// sends activation email
	mux.HandleFunc("/v1/mailer/activate", handleSendEmail(emailService))

	// receives activation email call
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// if r.Method != http.MethodPut {
		// 	w.WriteHeader(http.StatusMethodNotAllowed)
		// 	w.Write([]byte("{\"error\":\"invalid method\"}"))
		// 	return
		// }
		fmt.Println(r.URL.Path)
		fmt.Fprint(w, r.URL.Path)
	})

	srv := &http.Server{
		Handler: mux,
		Addr:    ":5000",
	}

	log.Println("starting server on port", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
