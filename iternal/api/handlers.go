package api

import (
	"TaskManager/iternal/models"
	"TaskManager/iternal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	strId := r.URL.Query().Get("id")
	if strId == "" {
		switch r.Method {
		case http.MethodGet:
			h.handleGetAll(w, r)
		case http.MethodPost:
			h.handlePost(w, r)
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
		h.handleGetById(w, r, id)
	case http.MethodDelete:
		h.handleDelete(w, r, id)
	case http.MethodPut:
		h.handlePut(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	tasks, _ := h.service.GetAll()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) handleGetById(w http.ResponseWriter, r *http.Request, id int) {
	task, _ := h.service.GetById(id)
	if task == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	t := &models.Task{
		Id:          0,
		Title:       body.Title,
		Description: body.Description,
	}

	h.service.Create(t)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func (h *TaskHandler) handlePut(w http.ResponseWriter, r *http.Request, id int) {
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	t := &models.Task{
		Id:          id,
		Title:       body.Title,
		Description: body.Description,
	}

	h.service.Update(t)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(t)
}

func (h *TaskHandler) handleDelete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.Delete(id); err != nil {
		http.Error(w, "not deleted", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
