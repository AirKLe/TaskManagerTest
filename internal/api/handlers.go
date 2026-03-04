package api

import (
	"TaskManager/internal/models"
	"TaskManager/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) handleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	var ve *service.ValidationError
	if errors.As(err, &ve) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"error":   "validation_error",
			"field":   ve.Field,
			"value":   ve.Value,
			"message": ve.Message,
		})
		return
	}

	var nf *service.NotFoundError
	if errors.As(err, &nf) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"error":    "not_found",
			"resource": nf.Resource,
			"id":       nf.Id,
		})
	}

	w.WriteHeader(http.StatusInternalServerError)

	json.NewEncoder(w).Encode(map[string]any{
		"error": "internal_error",
	})
}

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	strId := r.URL.Query().Get("id")
	if strId == "" {
		switch r.Method {
		case http.MethodGet:
			h.handleListTasks(w, r)
		case http.MethodPost:
			h.handleCreateTask(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}

		return
	}

	id, err := strconv.Atoi(strId)

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.handleGetTask(w, r, id)
	case http.MethodDelete:
		h.handleDeleteTask(w, r, id)
	case http.MethodPut:
		h.handleUpdateTask(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandler) handleListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.ListTasks()
	if err != nil {
		h.handleError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) handleGetTask(w http.ResponseWriter, r *http.Request, id int) {
	task, err := h.service.GetTask(id)
	if err != nil {
		h.handleError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var body createTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	t := &models.Task{
		Title:       body.Title,
		Description: body.Description,
	}

	h.service.CreateTask(t)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func (h *TaskHandler) handleUpdateTask(w http.ResponseWriter, r *http.Request, id int) {
	var body createTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	t := &models.Task{
		Id:          id,
		Title:       body.Title,
		Description: body.Description,
	}

	h.service.UpdateTask(t)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(t)
}

func (h *TaskHandler) handleDeleteTask(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.DeleteTask(id); err != nil {
		h.handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
