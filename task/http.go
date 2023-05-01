package task

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	req "github.com/ricxi/flat-list/shared/request"
	res "github.com/ricxi/flat-list/shared/response"
)

func NewHTTPHandler(service Service, middlewares ...func(http.Handler) http.Handler) http.Handler {
	h := &httpHandler{
		service: service,
	}

	r := chi.NewMux()
	r.Use(middlewares...)

	r.Route("/v1/task", func(r chi.Router) {
		r.Post("/", h.handleCreateTask)
		r.Get("/{id}", h.handleGetTask)
		r.Put("/", h.handleUpdateTask)
		r.Delete("/{id}", h.handleDeleteTask)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"success":"false", "message":"method not allowed"}`))
	})

	return r
}

type httpHandler struct {
	service Service
}

func (h *httpHandler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask NewTask // it doesn't need its date fields yet, but should I really create an entirely new data type for this?

	if err := req.ParseJSON(r, &newTask); err != nil {
		res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromCtx(r.Context())
	if err != nil {
		res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	newTask.UserID = userID

	taskID, err := h.service.createTask(r.Context(), &newTask)
	if err != nil {
		res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	res.SendSuccessJSON(w, res.Payload{"taskId": taskID}, http.StatusCreated, nil)
}

func (h *httpHandler) handleGetTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		res.SendErrorJSON(w, "missing url param id", http.StatusBadRequest)
		return
	}

	task, err := h.service.getTaskByID(r.Context(), taskID)
	if err != nil {
		// check for ErrTaskNotFound and return a status not found
		res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	res.SendSuccessJSON(w, res.Payload{"task": task}, http.StatusOK, nil)
}

func (h *httpHandler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var task Task

	userID, err := getUserIDFromCtx(r.Context())
	if err != nil {
		res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	task.UserID = userID

	if err := req.ParseJSON(r, &task); err != nil {
		res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedTask, err := h.service.updateTask(r.Context(), &task)
	if err != nil {
		res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	res.SendSuccessJSON(w, res.Payload{"task": updatedTask}, http.StatusOK, nil)
}

func (h *httpHandler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		res.SendErrorJSON(w, "missing url param id", http.StatusBadRequest)
		return
	}

	if err := h.service.deleteTask(r.Context(), taskID); err != nil {
		res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	body := map[string]any{
		"success": true,
	}

	res.SendJSON(w, &body, http.StatusOK, nil)
}
