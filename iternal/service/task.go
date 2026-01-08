package service

import (
	"TaskManager/iternal/models"
	"TaskManager/iternal/storage"
)

type TaskService struct {
	storage storage.TaskStorage
}

func NewTaskService(storage storage.TaskStorage) *TaskService {
	return &TaskService{storage: storage}
}

func (s *TaskService) GetById(id int) (*models.Task, error) {
	return s.storage.GetById(id)
}

func (s *TaskService) GetAll() (map[int]*models.Task, error) {
	return s.storage.GetAll()
}

func (s *TaskService) Create(t *models.Task) error {
	return s.storage.Create(t)
}

func (s *TaskService) Update(t *models.Task) error {
	return s.storage.Update(t)
}

func (s *TaskService) Delete(id int) error {
	return s.storage.Delete(id)
}
