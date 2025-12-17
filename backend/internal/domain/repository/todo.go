package repository

import (
	"context"

	"good-todo-go/internal/domain/model"
)

//go:generate go run go.uber.org/mock/mockgen -source=todo.go -destination=mock/todo.go -package=mock

type ITodoRepository interface {
	Create(ctx context.Context, todo *model.Todo) (*model.Todo, error)
	FindByID(ctx context.Context, id string) (*model.Todo, error)
	FindByUserID(ctx context.Context, userID string) ([]*model.Todo, error)
	FindPublicByTenantID(ctx context.Context, tenantID string) ([]*model.Todo, error)
	Update(ctx context.Context, todo *model.Todo) (*model.Todo, error)
	Delete(ctx context.Context, id string) error
}
