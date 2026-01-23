package domain

import "time"

type Todo struct {
	ID          *uint
	Title       string
	Description *string
	Done        bool
	Deadline    *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
