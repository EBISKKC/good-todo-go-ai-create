package repository

import (
	"context"

	"good-todo-go/internal/domain/model"
)

//go:generate go run go.uber.org/mock/mockgen -source=tenant.go -destination=mock/tenant.go -package=mock

type ITenantRepository interface {
	Create(ctx context.Context, tenant *model.Tenant) (*model.Tenant, error)
	FindByID(ctx context.Context, id string) (*model.Tenant, error)
	FindBySlug(ctx context.Context, slug string) (*model.Tenant, error)
}
