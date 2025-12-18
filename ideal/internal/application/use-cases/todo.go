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

//func (uc *TodoUseCase) Get() {}
//
//func (uc *TodoUseCase) Update() {}
//
//func (uc *TodoUseCase) Delete() {}
//
//func (uc *TodoUseCase) MarkCompleted() {}
