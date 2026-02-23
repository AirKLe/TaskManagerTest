package api

import (
	"TaskManager/internal/models"
	"TaskManager/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	service *service.TaskService
}

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
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
		http.Error(w, "method not allowed", http.StatusBadRequest)
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
	tasks, _ := h.service.ListTasks()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) handleGetTask(w http.ResponseWriter, r *http.Request, id int) {
	task, _ := h.service.GetTask(id)
	if task == nil {
		http.NotFound(w, r)
		return
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
		Id:          0,
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
		http.Error(w, "not deleted", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
