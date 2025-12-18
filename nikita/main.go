package main

import (
	"log"
	"net/http"

	"github.com/unomiqq/to-do-list-go/handler"
	"github.com/unomiqq/to-do-list-go/repository"
)

func main() {
	// Выбираем тип хранилища
	// repo := repository.NewMemoryRepository()
	repo, err := repository.NewJSONRepository("tasks.json")
	if err != nil {
		log.Fatal(err)
	}

	taskHandler := handler.NewTaskHandler(repo)

	// Роутинг
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Todo API is running"))
		if err != nil {
			return
		}
	})

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetTasks(w, r)
		case http.MethodPost:
			taskHandler.CreateTask(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, это /tasks/{id}/done или просто /tasks/{id}
		if len(r.URL.Path) > 7 && r.URL.Path[len(r.URL.Path)-5:] == "/done" {
			taskHandler.MarkTaskDone(w, r)
		} else {
			switch r.Method {
			case http.MethodGet:
				taskHandler.GetTask(w, r)
			case http.MethodPut:
				taskHandler.UpdateTask(w, r)
			case http.MethodDelete:
				taskHandler.DeleteTask(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
