package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// create a token for a given user
func handlerCreateToken(repo *repository) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		activationToken, err := generateActivationToken()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userID := ps.ByName("userID")
		info := ActivationTokenInfo{
			Token:  activationToken,
			UserID: userID,
		}

		if err := repo.InsertToken(r.Context(), &info); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write([]byte("{\"status\":\"success\"}"))
	}
}

// Return all the activation tokens that are listed under
// the same user id
func handleGetTokens(repo *repository) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		userID := ps.ByName("userID")

		tokens, err := repo.GetTokens(r.Context(), userID)
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
