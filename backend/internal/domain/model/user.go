package model

import "time"

type UserRole string

const (
	UserRoleAdmin  UserRole = "admin"
	UserRoleMember UserRole = "member"
)

type User struct {
	ID                         string
	TenantID                   string
	Email                      string
	PasswordHash               string
	Name                       string
	Role                       UserRole
	EmailVerified              bool
	VerificationToken          *string
	VerificationTokenExpiresAt *time.Time
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}
