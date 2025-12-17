package model

import "time"

type Todo struct {
	ID          string
	TenantID    string
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
