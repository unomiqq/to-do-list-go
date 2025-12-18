package repository

import "github.com/unomiqq/to-do-list-go/models"

type TaskRepository interface {
	Create(task *models.Task) error
	GetAll() ([]models.Task, error)
	GetByID(id int) (*models.Task, error)
	Update(task *models.Task) error
	MarkAsDone(id int) error
	Delete(id int) error
}
