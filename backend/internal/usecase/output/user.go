package output

import "time"

type UserOutput struct {
	ID            string
	TenantID      string
	Email         string
	Name          string
	Role          string
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
