package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	req "github.com/ricxi/flat-list/shared/request"
	res "github.com/ricxi/flat-list/shared/response"
)

type httpHandler struct {
	service Service
}

func NewHTTPHandler(service Service) http.Handler {
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

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		res.SendErrorJSON(w, "resource not found", http.StatusNotFound)
	})

	r.Route("/v1/user", func(r chi.Router) {
		r.Get("/healthcheck", h.handleHealthCheck)
		r.Post("/register", h.handleRegister)
		r.Post("/login", h.handleLogin)
		r.Put("/activate/{token}", h.handleActivate)
		r.Post("/restart/activation", h.handleRestartActivation)
		r.Post("/authenticate", h.handleAuthenticate)
	})

	return r
}

// Response is used to wrap user responses
type Response map[string]any

func (h httpHandler) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	res.SendJSON(w, Response{"success": true, "message": "user service is running"}, http.StatusOK, nil)
}

func (h httpHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var u UserRegistrationInfo
	if err := req.ParseJSON(r, &u); err != nil {
		res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.service.registerUser(r.Context(), u)
	if err != nil {
		// Should I return a 409 status code for a duplicate user?
		res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	res.SendSuccessJSON(w, res.Payload{"id": id}, http.StatusCreated, nil)
}

func (h httpHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var u UserLoginInfo
	if err := req.ParseJSON(r, &u); err != nil {
		res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	uInfo, err := h.service.loginUser(r.Context(), u)
	if err != nil {
		res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	res.SendSuccessJSON(w, res.Payload{"user": uInfo}, http.StatusOK, nil)
}

// handleActivate is called to activate a newly registered user's account
func (h httpHandler) handleActivate(w http.ResponseWriter, r *http.Request) {
	activationToken := chi.URLParam(r, "token")
	if activationToken == "" {
		res.SendErrorJSON(w, "missing activation token parameter", http.StatusBadRequest)
		return
	}

	if err := h.service.activateUser(r.Context(), activationToken); err != nil {
		res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleRestartActivation is called to generate a new activation token and resend a new activation email to a user
// TODO: Test this method
func (h httpHandler) handleRestartActivation(w http.ResponseWriter, r *http.Request) {
	var u UserLoginInfo
	if err := req.ParseJSON(r, &u); err != nil {
		res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.restartActivation(r.Context(), u); err != nil {
		res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h httpHandler) handleAuthenticate(w http.ResponseWriter, r *http.Request) {
	token := make(map[string]string)
	if err := req.ParseJSON(r, &token); err != nil {
		res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	if token["token"] == "" {
		res.SendErrorJSON(w, "no token provided", http.StatusBadRequest)
		return
	}

	userID, err := h.service.authenticate(r.Context(), token["token"])
	if err != nil {
		res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	res.SendJSON(w, map[string]string{"userId": userID}, http.StatusOK, nil)
}
