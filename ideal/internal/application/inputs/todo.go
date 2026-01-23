package inputs

import "time"

type CreateTodoInput struct {
	Title       string
	Description *string
	Deadline    *time.Time
}

type UpdateTodoInput struct {
	Title       *string
	Description *string
	Done        *bool
	Deadline    *time.Time
}
