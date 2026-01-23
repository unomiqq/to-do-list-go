package use_cases

import (
	"ideal-todo/internal/application/inputs"
	use_case_dependencies "ideal-todo/internal/application/use-cases/dependencies"
	"ideal-todo/internal/domain"
)

type TodoUseCase struct {
	todoRepo use_case_dependencies.TodoRepo
}

func NewTodoUseCase(repo use_case_dependencies.TodoRepo) *TodoUseCase {
	return &TodoUseCase{
		todoRepo: repo,
	}
}

func (uc *TodoUseCase) Create(input inputs.CreateTodoInput) (todo domain.Todo, err error) {
	todo = domain.Todo{
		Title:       input.Title,
		Description: input.Description,
		Deadline:    input.Deadline,
		Done:        false,
	}

	todo, err = uc.todoRepo.Create(todo)
	if err != nil {
		return domain.Todo{}, err
	}

	return todo, nil
}

func (uc *TodoUseCase) List() ([]domain.Todo, error) {
	res, err := uc.todoRepo.FindAll()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (uc *TodoUseCase) Get(id uint) (domain.Todo, error) {
	todo, err := uc.todoRepo.FindByID(id)
	if err != nil {
		return domain.Todo{}, err
	}

	return todo, nil
}

func (uc *TodoUseCase) Update(id uint, input inputs.UpdateTodoInput) (domain.Todo, error) {
	existing, err := uc.todoRepo.FindByID(id)
	if err != nil {
		return domain.Todo{}, err
	}

	if input.Title != nil {
		existing.Title = *input.Title
	}

	if input.Description != nil {
		existing.Description = input.Description
	}

	if input.Done != nil {
		existing.Done = *input.Done
	}

	if input.Deadline != nil {
		existing.Deadline = input.Deadline
	}

	updated, err := uc.todoRepo.Update(existing)
	if err != nil {
		return domain.Todo{}, err
	}

	return updated, nil
}

func (uc *TodoUseCase) Delete(id uint) error {
	err := uc.todoRepo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *TodoUseCase) MarkCompleted(id uint) (domain.Todo, error) {
	existing, err := uc.todoRepo.FindByID(id)
	if err != nil {
		return domain.Todo{}, err
	}

	existing.Done = true

	updated, err := uc.todoRepo.Update(existing)
	if err != nil {
		return domain.Todo{}, err
	}

	return updated, nil
}
