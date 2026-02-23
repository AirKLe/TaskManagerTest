package service

import (
	"TaskManager/internal/models"
	"TaskManager/internal/storage"
	"fmt"
)

type ValidationError struct {
	Field   string
	Value   any
	Message string
}

func (e *ValidationError) Error() string {
	if e.Value == nil {
		return fmt.Sprintf("validation failed: field = %s message = %s", e.Field, e.Message)
	}

	return fmt.Sprintf("validation failed: field = %s  value = %v message = %s", e.Field, e.Value, e.Message)
}

func NewValidationError(field string, value any, message string) error {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

type TaskService struct {
	storage storage.TaskStorage
}

func NewTaskService(storage storage.TaskStorage) *TaskService {
	return &TaskService{storage: storage}
}

func (s *TaskService) GetTask(id int) (*models.Task, error) {
	if id <= 0 {
		return nil, NewValidationError("id", id, "invalid")
	}

	return s.storage.GetById(id)
}

func (s *TaskService) ListTasks() ([]*models.Task, error) {
	return s.storage.GetAll()
}

func (s *TaskService) CreateTask(t *models.Task) error {
	if t == nil {
		return NewValidationError("task", nil, "is nil")
	}

	if t.Title == "" {
		return NewValidationError("title", nil, "invalid")
	}

	return s.storage.Create(t)
}

func (s *TaskService) UpdateTask(t *models.Task) error {
	if t == nil {
		return NewValidationError("task", nil, "is nil")
	}

	if t.Id <= 0 {
		return NewValidationError("id", t.Id, "invalid")
	}

	if t.Title == "" {
		return NewValidationError("title", nil, "invalid")
	}

	return s.storage.Update(t)
}

func (s *TaskService) DeleteTask(id int) error {
	if id <= 0 {
		return NewValidationError("id", id, "invalid")
	}

	return s.storage.Delete(id)
}
