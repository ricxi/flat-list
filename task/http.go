package task

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	req "github.com/ricxi/flat-list/shared/request"
	res "github.com/ricxi/flat-list/shared/response"
)

func NewHTTPHandler(service Service, middlewares ...func(http.Handler) http.Handler) http.Handler {
	h := &httpHandler{
		s: service,
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
	s Service
}

func (h *httpHandler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask NewTask
	userID, err := getUserIDFromCtx(r.Context())
	if err != nil {
		res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTask.UserID = userID
	if err := req.ParseJSON(r, &newTask); err != nil {
		res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	taskID, err := h.s.CreateTask(r.Context(), &newTask)
	if err != nil {
		res.SendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	body := map[string]any{
		"success": true,
		"taskId":  taskID,
	}

	res.SendJSON(w, &body, http.StatusCreated, nil)
}

func (h *httpHandler) handleGetTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		writeErrorToResponse(w, "missing url param id", http.StatusBadRequest)
		return
	}

	task, err := h.s.GetTaskByID(r.Context(), taskID)
	if err != nil {
		// check for ErrTaskNotFound and return a status not found
		writeErrorToResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	body := map[string]any{
		"success": true,
		"task":    task,
	}
	res.SendJSON(w, &body, http.StatusOK, nil)
}

func (h *httpHandler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var task Task

	userID, err := getUserIDFromCtx(r.Context())
	if err != nil {
		writeErrorToResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	task.UserID = userID

	if err := req.ParseJSON(r, &task); err != nil {
		writeErrorToResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedTask, err := h.s.UpdateTask(r.Context(), &task)
	if err != nil {
		writeErrorToResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	body := map[string]any{
		"success": true,
		"task":    updatedTask,
	}
	res.SendJSON(w, &body, http.StatusOK, nil)
}

func (h *httpHandler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		writeErrorToResponse(w, "missing url param id", http.StatusBadRequest)
		return
	}

	if err := h.s.DeleteTask(r.Context(), taskID); err != nil {
		writeErrorToResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	body := map[string]any{
		"success": true,
	}
	res.SendJSON(w, &body, http.StatusOK, nil)
}

func writeErrorToResponse(w http.ResponseWriter, message string, statusCode int) {
	errRes := map[string]any{
		"success": false,
		"message": message,
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(&errRes); err != nil {
		panic(err)
	}
}
