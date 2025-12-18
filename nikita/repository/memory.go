package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/unomiqq/to-do-list-go/models"
)

type MemoryRepository struct {
	tasks  map[int]*models.Task
	nextID int
	mu     sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		tasks:  make(map[int]*models.Task),
		nextID: 1,
	}
}

func (r *MemoryRepository) Create(task *models.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = r.nextID
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	r.tasks[task.ID] = task
	r.nextID++

	return nil
}

func (r *MemoryRepository) GetAll() ([]models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]models.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, *task)
	}

	return tasks, nil
}

func (r *MemoryRepository) GetByID(id int) (*models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}

	return task, nil
}

func (r *MemoryRepository) Update(task *models.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.tasks[task.ID]
	if !exists {
		return errors.New("task not found")
	}

	// created_at не должен меняться при обновлении
	task.CreatedAt = existing.CreatedAt
	task.UpdatedAt = time.Now()
	r.tasks[task.ID] = task

	return nil
}

func (r *MemoryRepository) MarkAsDone(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, exists := r.tasks[id]
	if !exists {
		return errors.New("task not found")
	}

	task.IsDone = true
	task.UpdatedAt = time.Now()

	return nil
}

func (r *MemoryRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return errors.New("task not found")
	}

	delete(r.tasks, id)

	return nil

}
