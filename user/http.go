package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/ricxi/flat-list/shared/request"
	"github.com/ricxi/flat-list/shared/response"
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

// Response is used to wrap user responses
type Response map[string]any

func (h httpHandler) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	response.SendJSON(w, Response{"success": true, "message": "user service is running"}, http.StatusOK, nil)
}

func (h httpHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var u UserRegistrationInfo
	if err := request.ParseJSON(r, &u); err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.service.RegisterUser(r.Context(), &u)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendJSON(w, Response{"id": id}, http.StatusCreated, nil)
}

func (h httpHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var u UserLoginInfo
	if err := request.ParseJSON(r, &u); err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	uInfo, err := h.service.LoginUser(r.Context(), &u)
	if err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendJSON(w, Response{"user": uInfo}, http.StatusOK, nil)
}

// handleActivate is called to activate a newly registered user's account
func (h httpHandler) handleActivate(w http.ResponseWriter, r *http.Request) {
	activationToken := chi.URLParam(r, "token")
	if activationToken == "" {
		response.SendErrorJSON(w, "missing activation token parameter", http.StatusBadRequest)
		return
	}

	if err := h.service.ActivateUser(r.Context(), activationToken); err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendJSON(w, Response{"status": "success"}, http.StatusOK, nil)
}

// handleReactivate is called to generate a new activation token and resend a new activation email to a user
// TODO: Test this method
func (h httpHandler) handleReactivate(w http.ResponseWriter, r *http.Request) {
	var u UserLoginInfo
	if err := request.ParseJSON(r, &u); err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.RestartActivation(r.Context(), &u); err != nil {
		response.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.SendJSON(w, Response{"status": "success"}, http.StatusOK, nil)
}
