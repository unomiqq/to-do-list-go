package rest

import (
	"encoding/json"
	"errors"
	"ideal-todo/internal/application/inputs"
	"ideal-todo/internal/transport/rest/dto"
	"net/http"
	"strconv"
	"strings"
	"time"

	use_cases "ideal-todo/internal/application/use-cases"
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

	var deadline *time.Time
	if data.Deadline != nil && *data.Deadline != "" {
		parsed, err := time.Parse(time.RFC3339, *data.Deadline)
		if err != nil {
			http.Error(w, "invalid deadline format, use RFC3339 (e.g. 2026-01-22T21:30:00Z)", http.StatusBadRequest)
			return
		}
		deadline = &parsed
	}

	input := inputs.CreateTodoInput{
		Title:       data.Title,
		Description: data.Description,
		Deadline:    deadline,
	}

	todo, err := h.useCase.Create(input)
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

func (h *TodoHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := parseIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid todo ID", http.StatusBadRequest)
		return
	}

	todo, err := h.useCase.Get(id)
	if err != nil {
		if err.Error() == "todo not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(dto.ToDTO(todo))
}

func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := parseIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid todo ID", http.StatusBadRequest)
		return
	}

	var data dto.UpdateTodoDTO
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	var deadline *time.Time
	if data.Deadline != nil && *data.Deadline != "" {
		parsed, err := time.Parse(time.RFC3339, *data.Deadline)
		if err != nil {
			http.Error(w, "invalid deadline format, use RFC3339 (e.g. 2026-01-22T21:30:00Z)", http.StatusBadRequest)
			return
		}
		deadline = &parsed
	} else if data.Deadline != nil && *data.Deadline == "" {
		deadline = nil
	}

	input := inputs.UpdateTodoInput{
		Title:       data.Title,
		Description: data.Description,
		Done:        data.Done,
		Deadline:    deadline,
	}

	todo, err := h.useCase.Update(id, input)
	if err != nil {
		if err.Error() == "todo not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(dto.ToDTO(todo))
}

func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := parseIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid todo ID", http.StatusBadRequest)
		return
	}

	err = h.useCase.Delete(id)
	if err != nil {
		if err.Error() == "todo not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TodoHandler) MarkCompleted(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := parseIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid todo ID", http.StatusBadRequest)
		return
	}

	todo, err := h.useCase.MarkCompleted(id)
	if err != nil {
		if err.Error() == "todo not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(dto.ToDTO(todo))
}

func parseIDFromPath(path string) (uint, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return 0, errors.New("invalid path format")
	}

	if parts[0] != "todo" {
		return 0, errors.New("invalid path format")
	}

	id, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
