package app

import (
	"log"
	"net/http"

	"ideal-todo/internal/application/use-cases"
	"ideal-todo/internal/storage/memory"
	"ideal-todo/internal/transport/rest"
)

type RestApp struct {
	addr string
}

func NewRestApp() *RestApp {
	return &RestApp{
		addr: ":8080",
	}
}

func (a *RestApp) Run() {
	repo := memory.NewTodoRepo()

	todoUC := use_cases.NewTodoUseCase(repo)

	todoHandler := rest.NewTodoHandler(todoUC)

	mux := http.NewServeMux()
	mux.HandleFunc("/todo", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			todoHandler.List(w, r)
		case http.MethodPost:
			todoHandler.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Printf("REST server running on %s\n", a.addr)
	log.Fatal(http.ListenAndServe(a.addr, mux))
}
