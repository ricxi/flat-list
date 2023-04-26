package mailer

import (
	"encoding/json"
	"net/http"
)

func HandleSendActivationEmail(mailerService *MailerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data EmailActivationData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := mailerService.sendActivationEmail(data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write([]byte("success!"))
	}
}
