package memory

import (
	"api/internal/domain/models"
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"
)

type Storage struct {
	mu    sync.RWMutex
	tasks map[string]*models.Task
}

func New() *Storage {
	return &Storage{
		tasks: make(map[string]*models.Task),
	}
}

func (s *Storage) CreateTask(ctx context.Context, task *models.Task) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := generateID()
	task.ID = id
	s.tasks[id] = task
	return task.ID, nil
}

func (s *Storage) Task(ctx context.Context, id string) (models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return models.Task{}, errors.New("task not found")
	}
	return *task, nil
}

func (s *Storage) UpdateTask(ctx context.Context, task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[task.ID]; !exists {
		return errors.New("task not found")
	}
	s.tasks[task.ID] = task
	return nil
}

func (s *Storage) DeleteTask(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return errors.New("task not found")
	}
	delete(s.tasks, id)
	return nil
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randString(6)
}

func randString(n int) string {
	letters := "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
