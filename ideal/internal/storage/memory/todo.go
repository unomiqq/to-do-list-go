package memory

import (
	"errors"
	"ideal-todo/internal/domain"
	"ideal-todo/internal/storage"
	"sync"
	"time"
)

type TodoModel struct {
	ID          uint
	Title       string
	Description *string
	Done        bool
	Deadline    *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TodoRepo struct {
	todos  []*TodoModel
	nextID uint
	mu     sync.RWMutex
}

func NewTodoRepo() storage.TodoRepo {
	tr := TodoRepo{
		todos:  make([]*TodoModel, 0),
		nextID: 1,
		mu:     sync.RWMutex{},
	}

	return &tr
}

func (r *TodoRepo) FindByID(id uint) (domain.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, t := range r.todos {
		if t.ID == id {
			return t.modelToDomain(), nil
		}
	}

	return domain.Todo{}, errors.New("todo not found")
}

func (r *TodoRepo) FindAll() ([]domain.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := make([]domain.Todo, 0, len(r.todos))
	for _, t := range r.todos {
		res = append(res, t.modelToDomain())
	}

	return res, nil
}

func (r *TodoRepo) Create(t domain.Todo) (domain.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextID
	t.ID = &id
	r.nextID++

	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	m := &TodoModel{}
	m.domainToModel(t)

	r.todos = append(r.todos, m)

	return m.modelToDomain(), nil
}

func (r *TodoRepo) Update(t domain.Todo) (domain.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if t.ID == nil {
		return domain.Todo{}, errors.New("todo ID is required")
	}

	var existing *TodoModel
	for _, todo := range r.todos {
		if todo.ID == *t.ID {
			existing = todo
			break
		}
	}

	if existing == nil {
		return domain.Todo{}, errors.New("todo not found")
	}

	t.CreatedAt = existing.CreatedAt
	t.UpdatedAt = time.Now()

	existing.domainToModel(t)

	return existing.modelToDomain(), nil
}

func (r *TodoRepo) Delete(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	index := -1
	for i, todo := range r.todos {
		if todo.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("todo not found")
	}

	r.todos = append(r.todos[:index], r.todos[index+1:]...)

	return nil
}

func (m *TodoModel) domainToModel(d domain.Todo) {
	m.ID = *d.ID
	m.Title = d.Title
	m.Description = d.Description
	m.Done = d.Done
	m.Deadline = d.Deadline
	m.CreatedAt = d.CreatedAt
	m.UpdatedAt = d.UpdatedAt
}

func (m *TodoModel) modelToDomain() domain.Todo {
	return domain.Todo{
		ID:          &m.ID,
		Title:       m.Title,
		Description: m.Description,
		Done:        m.Done,
		Deadline:    m.Deadline,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
