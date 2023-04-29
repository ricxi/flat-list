package mailer

import (
	"encoding/json"
	"net/http"

	res "github.com/ricxi/flat-list/shared/response"
)

func HandleSendActivationEmail(mailerService *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data ActivationEmailData

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			res.SendInternalServerErrorAsJSON(w, err.Error())
			return
		}

		if err := mailerService.sendActivationEmail(data); err != nil {
			res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.SendJSON(w, map[string]any{"success": true}, http.StatusOK, nil)
	}
}
