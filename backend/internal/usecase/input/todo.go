package input

import "time"

type CreateTodoInput struct {
	Title       string
	Description string
	IsPublic    bool
	DueDate     *time.Time
}

type UpdateTodoInput struct {
	ID          string
	Title       string
	Description string
	Completed   bool
	IsPublic    bool
	DueDate     *time.Time
}
