package rest

import (
	"encoding/json"
	"ideal-todo/internal/application/inputs"
	"ideal-todo/internal/transport/rest/dto"
	"net/http"

	"ideal-todo/internal/application/use-cases"
)

type TodoHandler struct {
	useCase *use_cases.TodoUseCase
}

func NewTodoHandler(uc *use_cases.TodoUseCase) *TodoHandler {
	return &TodoHandler{
		useCase: uc,
	}
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data dto.CreateTodoDTO
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if data.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	todo, err := h.useCase.Create(inputs.CreateTodoInput(data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(dto.ToDTO(todo))
}

func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	res, err := h.useCase.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(dto.ToDTOs(res))
}
