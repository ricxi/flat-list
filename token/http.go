package token

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func NewHTTPHandler(repository Repository) http.Handler {
	h := httpHandler{
		repository: repository,
	}

	r := httprouter.New()

	r.POST("/v1/token/activation/:userId", h.handleCreateToken)
	r.GET("/v1/token/:userId", h.handleValidateToken)

	return r
}

type httpHandler struct {
	repository Repository
}

// create a token for a given user
func (h *httpHandler) handleCreateToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	activationToken, err := generateActivationToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := ps.ByName("userId")
	if userID == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	info := ActivationTokenInfo{
		Token:  activationToken,
		UserID: userID,
	}

	if err := h.repository.InsertActivationToken(r.Context(), &info); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"token": activationToken}); err != nil {
		panic(err)
	}
}

// Return all the activation tokens that are listed under
// the same user id
func (h *httpHandler) handleValidateToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	activationToken := ps.ByName("token")
	if activationToken == "" {
		http.Error(w, "token required", http.StatusBadRequest)
		return
	}

	userID, err := h.repository.GetUserID(r.Context(), activationToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(struct {
		UserID string `json:"userId"`
	}{
		UserID: userID,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
