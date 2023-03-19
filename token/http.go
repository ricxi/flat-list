package token

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// create a token for a given user
func HandlerCreateToken(repo *Repository) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		activationToken, err := generateActivationToken()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userID := ps.ByName("userID")
		if userID == "" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		info := ActivationTokenInfo{
			Token:  activationToken,
			UserID: userID,
		}

		if err := repo.InsertActivationToken(r.Context(), &info); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]string{"token": activationToken}); err != nil {
			panic(err)
		}
	}
}

// Return all the activation tokens that are listed under
// the same user id
func HandleGetTokens(repo *Repository) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		userID := ps.ByName("userID")

		tokens, err := repo.GetActivationTokens(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := json.NewEncoder(w).Encode(struct {
			Tokens []string `json:"tokens"`
		}{
			Tokens: tokens,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
