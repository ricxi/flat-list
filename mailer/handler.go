package main

import (
	"encoding/json"
	"net/http"
)

func handleSendEmail(es *EmailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data UserActivationData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := es.SendActivationEmail(data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write([]byte("success!"))
	}
}
