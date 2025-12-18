package storage

import "ideal-todo/internal/domain"

type TodoRepo interface {
	FindByID(id uint) (domain.Todo, error)
	FindAll() ([]domain.Todo, error)
	Create(todo domain.Todo) (domain.Todo, error)
	Update(todo domain.Todo) (domain.Todo, error)
	Delete(id uint) error
}
