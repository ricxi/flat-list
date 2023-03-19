package user

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type httpHandler struct {
	service Service
}

func NewHandler(service Service) http.Handler {
	h := httpHandler{service: service}

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/v1/user", func(r chi.Router) {
		r.Get("/healthcheck", h.handleHealthCheck)
		r.Post("/register", h.handleRegister)
		r.Post("/login", h.handleLogin)
		r.Put("/activate/{token}", h.handleActivate)
		r.Post("/reactivate", h.handleReactivate)
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

	id, err := h.service.RegisterUser(r.Context(), &u)
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

// handleActivate is called to activate a newly registered user's account
func (h httpHandler) handleActivate(w http.ResponseWriter, r *http.Request) {
	activationToken := chi.URLParam(r, "token")
	if activationToken == "" {
		log.Println("ERROR: missing token")
		writeErrorToResponse(w, "missing token paramater", http.StatusBadRequest)
		return
	}

	if err := h.service.ActivateUser(r.Context(), activationToken); err != nil {
		writeErrorToResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeToResponse(w, Response{"status": "success"}, http.StatusOK)
}

// handleReactivate is called to generate a new activation token and resend a new activation email to a user
func (h httpHandler) handleReactivate(w http.ResponseWriter, r *http.Request) {

	var u UserLoginInfo
	if err := readFromRequest(w, r, &u); err != nil {
		writeErrorToResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.RestartActivation(r.Context(), &u); err != nil {
		writeErrorToResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeToResponse(w, Response{"status": "success"}, http.StatusOK)
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
