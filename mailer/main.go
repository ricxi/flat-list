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
  
  from := "team@flatlist.com"
  to := "dwightschrute@dundermifflin.com"
  subject := "Please activate your account"
  body := "Hello Dwight,</br>Follow the instructions to activate your account."

  m.Send(from, to, subject, body)
}


func handleSendEmail(m Mailer) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {


  }
}

