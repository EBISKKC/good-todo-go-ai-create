package repository

import (
	"context"

	"good-todo-go/internal/domain/model"
)

//go:generate go run go.uber.org/mock/mockgen -source=user.go -destination=mock/user.go -package=mock

type IUserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByEmail(ctx context.Context, tenantID, email string) (*model.User, error)
	FindByVerificationToken(ctx context.Context, token string) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
}
