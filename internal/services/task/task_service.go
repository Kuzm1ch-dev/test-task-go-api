package task

import (
	"api/internal/domain/models"
	"context"
	"math/rand"
	"sync"
	"time"
)

type Storage interface {
	CreateTask(ctx context.Context, task *models.Task) (string, error)
	Task(ctx context.Context, id string) (models.Task, error)
	DeleteTask(ctx context.Context, id string) error
	UpdateTask(ctx context.Context, task *models.Task) error
}

type Service struct {
	storage   Storage
	cancelers map[string]context.CancelFunc
	mu        sync.RWMutex
}

func New(storage Storage) *Service {
	return &Service{
		storage:   storage,
		cancelers: make(map[string]context.CancelFunc),
	}
}

func (s *Service) CreateTask(ctx context.Context) (string, error) {
	processingTime := rand.Intn(3*60) + 3*60
	now := time.Now().Unix()

	task := &models.Task{
		ProcessingTime: int64(processingTime),
		Status:         "processing",
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	id, err := s.storage.CreateTask(ctx, task)
	if err != nil {
		return "", err
	}

	taskCtx, cancel := context.WithCancel(context.Background())
	s.mu.Lock()
	s.cancelers[id] = cancel
	s.mu.Unlock()

	go s.processTask(taskCtx, id, processingTime)

	return id, nil
}

func (s *Service) GetTask(ctx context.Context, id string) (models.Task, error) {
	return s.storage.Task(ctx, id)
}

func (s *Service) DeleteTask(ctx context.Context, id string) error {
	s.mu.Lock()
	if cancel, exists := s.cancelers[id]; exists {
		cancel()
		delete(s.cancelers, id)
	}
	s.mu.Unlock()

	return s.storage.DeleteTask(ctx, id)
}

func (s *Service) processTask(ctx context.Context, taskID string, processingTime int) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	startTime := time.Now()
	totalDuration := time.Duration(processingTime) * time.Second

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			elapsed := time.Since(startTime)
			progress := elapsed.Seconds()
			
			if elapsed >= totalDuration {
				s.completeTask(taskID)
				return
			}

			s.updateTaskProgress(taskID, progress)
		}
	}
}

func (s *Service) updateTaskProgress(taskID string, progress float64) {
	task, err := s.storage.Task(context.Background(), taskID)
	if err != nil {
		return
	}

	task.Result = map[string]interface{}{
		"progress": progress,
		"message":  "Processing...",
	}
	task.UpdatedAt = time.Now().Unix()

	s.storage.UpdateTask(context.Background(), &task)
}

func (s *Service) completeTask(taskID string) {
	task, err := s.storage.Task(context.Background(), taskID)
	if err != nil {
		return
	}

	task.Status = "completed"
	task.Result = map[string]interface{}{
		"progress": 100.0,
		"message":  "Task completed successfully",
	}
	task.UpdatedAt = time.Now().Unix()

	s.storage.UpdateTask(context.Background(), &task)

	s.mu.Lock()
	delete(s.cancelers, taskID)
	s.mu.Unlock()
}
