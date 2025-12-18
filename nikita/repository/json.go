package repository

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/unomiqq/to-do-list-go/models"
)

type JSONRepository struct {
	filePath string
	tasks    map[int]*models.Task
	nextID   int
	mu       sync.RWMutex
}

func NewJSONRepository(filePath string) (*JSONRepository, error) {
	repo := &JSONRepository{
		filePath: filePath,
		tasks:    make(map[int]*models.Task),
		nextID:   1,
	}

	if err := repo.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return repo, nil
}

func (r *JSONRepository) load() error {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	var tasks []models.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return err
	}

	for i := range tasks {
		r.tasks[tasks[i].ID] = &tasks[i]
		if tasks[i].ID >= r.nextID {
			r.nextID = tasks[i].ID + 1
		}
	}

	return nil
}

func (r *JSONRepository) save() error {
	tasks := make([]models.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, *task)
	}

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, data, 0644)
}

func (r *JSONRepository) Create(task *models.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = r.nextID
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	r.tasks[task.ID] = task
	r.nextID++

	return r.save()
}

func (r *JSONRepository) GetAll() ([]models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]models.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, *task)
	}

	return tasks, nil
}

func (r *JSONRepository) GetByID(id int) (*models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}

	return task, nil
}

func (r *JSONRepository) Update(task *models.Task) error {
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

	return r.save()
}

func (r *JSONRepository) MarkAsDone(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, exists := r.tasks[id]
	if !exists {
		return errors.New("task not found")
	}

	task.IsDone = true
	task.UpdatedAt = time.Now()

	return r.save()
}

func (r *JSONRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, exists := r.tasks[id]
	if !exists {
		return errors.New("task not found")
	}

	delete(r.tasks, task.ID)

	return r.save()
}
