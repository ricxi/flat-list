package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

func main() {
  conf, err := setupConfig()
  if err != nil {
    log.Fatal(err)
  }

  m := NewMailer(conf.username, conf.password, conf.host, conf.port)
  
  mux := http.NewServeMux()

  mux.HandleFunc("/v1/send", handleSendEmail(m))

  srv := &http.Server{
    Handler: mux,
    Addr: ":5000",
  }

  log.Println("starting server on port", srv.Addr)
  if err := srv.ListenAndServe(); err != nil {
    log.Fatal(err)
  }
  // from := "team@flatlist.com"
  // to := "dwightschrute@dundermifflin.com"
  // subject := "Please activate your account"
  // body := "Hello Dwight,</br>Follow the instructions to activate your account."
  // m.Send(from, to, subject, body)
}


func handleSendEmail(m *Mailer) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {

    type EmailData struct {
      From string `json:"from"`
      User struct {
        FirstName string `json:"firstName"`
        LastName string `json:"lastName"`
        Email string `json:"email"`
      } `json:"user"`
    }

    var data EmailData
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    t, err := template.ParseFiles("./useractivation.html")
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    body := new(bytes.Buffer)
    t.Execute(body, struct{FirstName string}{FirstName: data.User.FirstName})


    subject := "email test"
    if err := m.Send(data.From, data.User.Email, subject, body.String()); err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
    } else {
      w.Write([]byte("success!"))
    }
  }
}

