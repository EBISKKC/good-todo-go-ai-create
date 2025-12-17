package output

import "time"

type TodoOutput struct {
	ID          string
	UserID      string
	Title       string
	Description string
	Completed   bool
	IsPublic    bool
	DueDate     *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
