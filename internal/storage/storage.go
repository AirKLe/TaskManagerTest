package storage

import (
	"TaskManager/internal/models"
	"errors"
	"fmt"
)

type TaskStorage interface {
	GetById(id int) (*models.Task, error)
	GetAll() ([]*models.Task, error)
	Create(t *models.Task) error
	Update(t *models.Task) error
	Delete(id int) error
}

var ErrNotFound = errors.New("not found")

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
		return nil, fmt.Errorf("task %v :%w", id, ErrNotFound)
	}

	copy := *task
	return &copy, nil
}

func (s *inMemoryTaskStorage) GetAll() ([]*models.Task, error) {
	tasks := make([]*models.Task, 0, len(s.data))
	for _, t := range s.data {
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *inMemoryTaskStorage) Create(t *models.Task) error {
	t.Id = s.nextId
	s.data[t.Id] = t
	s.nextId++
	return nil
}

func (s *inMemoryTaskStorage) Update(t *models.Task) error {
	if _, ok := s.data[t.Id]; !ok {
		return fmt.Errorf("task %v: %w", t.Id, ErrNotFound)
	}

	s.data[t.Id] = t
	return nil
}

func (s *inMemoryTaskStorage) Delete(id int) error {
	if _, ok := s.data[id]; ok {
		delete(s.data, id)
		return nil
	}

	return fmt.Errorf("task %v: %w", id, ErrNotFound)
}
