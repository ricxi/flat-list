package user

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type httpHandler struct {
	service Service
}

func NewHandler(service Service) http.Handler {
	h := httpHandler{service: service}

	r := chi.NewRouter()
	r.Route("/v1/user", func(r chi.Router) {
		r.Get("/healthcheck", h.handleHealthCheck)
		r.Post("/register", h.handleRegister)
		r.Post("/login", h.handleLogin)
	})

	return r
}

func (h httpHandler) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	writeToResponse(w, Response{"status": "success", "message": "user service is running"}, http.StatusOK)
}

func (h httpHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var u UserRegistrationInfo
	readFromRequest(w, r, &u)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := h.service.RegisterUser(ctx, &u)
	if err != nil {
		writeErrorToResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeToResponse(w, Response{"id": id}, http.StatusCreated)
}

func (h httpHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var u UserLoginInfo
	readFromRequest(w, r, &u)

	uInfo, err := h.service.LoginUser(r.Context(), &u)
	if err != nil {
		writeErrorToResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeToResponse(w, Response{"user": uInfo}, http.StatusOK)
}

type Response map[string]any

func readFromRequest(w http.ResponseWriter, r *http.Request, dest any) {
	err := json.NewDecoder(r.Body).Decode(dest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func writeToResponse(w http.ResponseWriter, res Response, statusCode int) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func writeErrorToResponse(w http.ResponseWriter, message string, statusCode int) {
	writeToResponse(w, Response{"error": message}, statusCode)
}
