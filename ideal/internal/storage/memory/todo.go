package memory

import (
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
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TodoRepo struct {
	todos []*TodoModel
	mu    sync.RWMutex
}

func NewTodoRepo() storage.TodoRepo {
	tr := TodoRepo{
		todos: make([]*TodoModel, 0),
		mu:    sync.RWMutex{},
	}

	return &tr
}

func (r *TodoRepo) FindByID(id uint) (domain.Todo, error) {
	return domain.Todo{}, nil
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

	id := uint(len(r.todos)) + 1
	t.ID = &id

	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	m := &TodoModel{}
	m.domainToModel(t)

	r.todos = append(r.todos, m)

	return m.modelToDomain(), nil
}

func (r *TodoRepo) Update(t domain.Todo) (domain.Todo, error) {
	return domain.Todo{}, nil
}

func (r *TodoRepo) Delete(id uint) error {
	return nil
}

func (m *TodoModel) domainToModel(d domain.Todo) {
	m.ID = *d.ID
	m.Title = d.Title
	m.Description = d.Description
	m.Done = d.Done
	m.CreatedAt = d.CreatedAt
	m.UpdatedAt = d.UpdatedAt
}

func (m *TodoModel) modelToDomain() domain.Todo {
	return domain.Todo{
		ID:          &m.ID,
		Title:       m.Title,
		Description: m.Description,
		Done:        m.Done,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
