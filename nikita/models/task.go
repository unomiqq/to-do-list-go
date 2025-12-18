package models

import "time"

type Task struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	IsDone      bool       `json:"is_done"`
	Deadline    *time.Time `json:"deadline,omitempty"`
}
