package use_case_dependencies

import "ideal-todo/internal/domain"

type TodoRepo interface {
	FindAll() ([]domain.Todo, error)
	Create(todo domain.Todo) (domain.Todo, error)
}
