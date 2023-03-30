package task

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type httpHandler struct {
	s Service
}

func NewHTTPHandler(s Service) *httpHandler {
	return &httpHandler{
		s: s,
	}
}

func (h *httpHandler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var nt NewTask
	if err := json.NewDecoder(r.Body).Decode(&nt); err != nil {
		writeErrorToResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	taskID, err := h.s.CreateTask(r.Context(), &nt)
	if err != nil {
		writeErrorToResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := map[string]any{
		"success": true,
		"taskId":  taskID,
	}
	if err := writeToResponse(w, res, http.StatusCreated); err != nil {
		writeErrorToResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *httpHandler) handleGetTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		writeErrorToResponse(w, "missing url param id", http.StatusBadRequest)
		return
	}

	t, err := h.s.GetTaskByID(r.Context(), taskID)
	if err != nil {
		writeErrorToResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := map[string]any{
		"success": true,
		"task":    t,
	}
	if err := writeToResponse(w, res, http.StatusOK); err != nil {
		writeErrorToResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// func (h *httpHandler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {}
// func (h *httpHandler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {}

func writeToResponse(w http.ResponseWriter, res map[string]any, statusCode int) error {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		return err
	}

	return nil
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
