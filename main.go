package main

// barbus loh

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

// domain
type Task struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

//storage

type TaskStorage interface {
	GetByID(id int) (*Task, error)
	GetAll() (map[int]*Task, error)
	Post(t *Task) error
	Delete(id int) error
}

type InMemoryTaskStorage struct {
	data map[int]*Task
}

func NewInMemoryTaskStorage() *InMemoryTaskStorage {
	return &InMemoryTaskStorage{
		data: make(map[int]*Task),
	}
}

func (s *InMemoryTaskStorage) GetByID(id int) (*Task, error) {
	t, ok := s.data[id]
	if !ok {
		err := fmt.Errorf("нет таски с таким id %d", id)
		return nil, err
	}

	return t, nil
}

func (s *InMemoryTaskStorage) GetAll() (map[int]*Task, error) {
	return s.data, nil
}

func (s *InMemoryTaskStorage) Post(t *Task) error {
	newId := rand.Intn(10000)
	_, alreadyUse := s.data[newId]
	for alreadyUse == true {
		newId = rand.Intn(10000)
		_, alreadyUse = s.data[newId]
	}
	t.Id = newId
	s.data[t.Id] = t
	return nil
}

func (s *InMemoryTaskStorage) Delete(id int) error {
	if _, ok := s.data[id]; ok {
		delete(s.data, id)
		return nil
	}

	return fmt.Errorf("нет таски с таким id %d", id)
}

//service

type TaskService struct {
	storage TaskStorage
}

func NewTaskService(storage TaskStorage) *TaskService {
	return &TaskService{storage: storage}
}

func (s *TaskService) GetAll() (map[int]*Task, error) {
	return s.storage.GetAll()
}
func (s *TaskService) GetByID(id int) (*Task, error) {
	return s.storage.GetByID(id)
}

func (s *TaskService) Post(t *Task) error {
	return s.storage.Post(t)
}

func (s *TaskService) Delete(id int) error {
	return s.storage.Delete(id)
}

//handle

type TaskHandler struct {
	service *TaskService
}

func NewTaskHandler(service *TaskService) *TaskHandler {
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
	task, _ := h.service.GetByID(id)
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

	t := &Task{
		Id:          0,
		Title:       body.Title,
		Description: body.Description,
	}

	h.service.Post(t)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func (h *TaskHandler) handleDelete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.Delete(id); err != nil {
		http.Error(w, "not deleted", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	storage := NewInMemoryTaskStorage()
	service := NewTaskService(storage)
	handler := NewTaskHandler(service)

	http.Handle("/tasks/", handler)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
