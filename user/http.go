package user

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type httpHandler struct {
	service Service
}

func NewHandler(service Service) http.Handler {
	h := httpHandler{service: service}

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
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
	if err := readFromRequest(w, r, &u); err != nil {
		writeErrorToResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// remove this and use the request context
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
	if err := readFromRequest(w, r, &u); err != nil {
		writeErrorToResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	uInfo, err := h.service.LoginUser(r.Context(), &u)
	if err != nil {
		writeErrorToResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeToResponse(w, Response{"user": uInfo}, http.StatusOK)
}

func readFromRequest(w http.ResponseWriter, r *http.Request, dest any) error {
	err := json.NewDecoder(r.Body).Decode(dest)
	if err != nil {
		return err
	}

	return nil
}

type Response map[string]any

// implements WriterTo interface
type errorResponseLogger struct {
	bytes.Buffer
}

// A proxy that logs an error if one occurs while writing to the response
func (e errorResponseLogger) WriteTo(w io.Writer) (n int64, err error) {
	n, err = e.Buffer.WriteTo(w)
	if err != nil {
		log.Println("writing to response failed: ", err.Error())
	}
	return
}

func writeToResponse(w http.ResponseWriter, res Response, statusCode int) {
	var bodyBuffer bytes.Buffer
	err := json.NewEncoder(&bodyBuffer).Encode(res)
	if err != nil {
		// Should this panic?
		writeErrorToResponse(w, "json encoder: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)

	b := errorResponseLogger{bodyBuffer}
	if _, err := b.WriteTo(w); err != nil {
		// panic and log error to avoid superfluous header warning? Best way to set the error code?
		panic(err)
	}
}

func writeErrorToResponse(w http.ResponseWriter, message string, statusCode int) {
	if message == "" {
		message = "unknown"
	}

	var bodyBuffer bytes.Buffer
	if err := json.NewEncoder(&bodyBuffer).Encode(Response{"error": message}); err != nil {
		// Should I panic instead?
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	b := errorResponseLogger{bodyBuffer}
	if _, err := b.WriteTo(w); err != nil {
		panic(err)
	}
}
