package repository

import (
	"context"
)

//go:generate go run go.uber.org/mock/mockgen -source=auth.go -destination=mock/auth.go -package=mock

type IAuthRepository interface {
	SendVerificationEmail(ctx context.Context, email, token string) error
}
