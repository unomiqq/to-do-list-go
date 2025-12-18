package dto

import "ideal-todo/internal/domain"

type CreateTodoDTO struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
}

type TodoReturnDTO struct {
	ID          uint    `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Done        bool    `json:"done"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func ToDTO(t domain.Todo) (todo TodoReturnDTO) {
	return TodoReturnDTO{
		ID:          *t.ID,
		Title:       t.Title,
		Description: t.Description,
		Done:        t.Done,
		CreatedAt:   t.CreatedAt.String(),
		UpdatedAt:   t.UpdatedAt.String(),
	}
}

func ToDTOs(t []domain.Todo) (todos []TodoReturnDTO) {
	for _, v := range t {
		todos = append(todos, ToDTO(v))
	}

	return todos
}
