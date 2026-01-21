package storage

import (
	"TaskManager/iternal/models"
	"fmt"
)

type TaskStorage interface {
	GetById(id int) (*models.Task, error)
	GetAll() (map[int]*models.Task, error)
	Create(t *models.Task) error
	Update(t *models.Task) error
	Delete(id int) error
}

type inMemoryTaskStorage struct {
	data   map[int]*models.Task
	nextId int
}

func NewInMemoryTaskStorage() *inMemoryTaskStorage {
	return &inMemoryTaskStorage{
		data:   make(map[int]*models.Task),
		nextId: 1,
	}
}

func (s *inMemoryTaskStorage) GetById(id int) (*models.Task, error) {
	task, ok := s.data[id]
	if !ok {
		return nil, fmt.Errorf("There's no %d id", id)
	}

	return task, nil
}

func (s *inMemoryTaskStorage) GetAll() (map[int]*models.Task, error) {
	if len(s.data) == 0 {
		return nil, fmt.Errorf("Task list is empty")
	}

	return s.data, nil
}

func (s *inMemoryTaskStorage) Create(t *models.Task) error {
	t.Id = s.nextId
	s.data[t.Id] = t
	s.nextId++
	return nil
}

func (s *inMemoryTaskStorage) Update(t *models.Task) error {
	if t.Id <= 0 {
		return fmt.Errorf("Некорректный id %d", t.Id)
	}

	s.data[t.Id] = t
	return nil
}

func (s *inMemoryTaskStorage) Delete(id int) error {
	if _, ok := s.data[id]; ok {
		delete(s.data, id)
		return nil
	}

	return fmt.Errorf("нет таски с таким id %d", id)
}
