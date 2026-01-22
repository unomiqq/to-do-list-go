package app

import (
	"log"
	"net/http"

	use_cases "ideal-todo/internal/application/use-cases"
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

	mux.HandleFunc("/todo/", func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) > 7 && r.URL.Path[len(r.URL.Path)-5:] == "/done" {
			if r.Method == http.MethodPut {
				todoHandler.MarkCompleted(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else {
			switch r.Method {
			case http.MethodGet:
				todoHandler.Get(w, r)
			case http.MethodPut:
				todoHandler.Update(w, r)
			case http.MethodDelete:
				todoHandler.Delete(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})

	log.Printf("REST server running on %s\n", a.addr)
	log.Fatal(http.ListenAndServe(a.addr, mux))
}
