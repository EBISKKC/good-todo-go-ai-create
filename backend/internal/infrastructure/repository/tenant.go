package repository

import (
	"context"
	"fmt"

	"good-todo-go/internal/domain/model"
	"good-todo-go/internal/domain/repository"
	"good-todo-go/internal/ent/generated"
	"good-todo-go/internal/ent/generated/tenant"
)

type TenantRepository struct {
	client *generated.Client
}

func NewTenantRepository(client *generated.Client) repository.ITenantRepository {
	return &TenantRepository{client: client}
}

func (r *TenantRepository) Create(ctx context.Context, t *model.Tenant) (*model.Tenant, error) {
	created, err := r.client.Tenant.Create().
		SetID(t.ID).
		SetName(t.Name).
		SetSlug(t.Slug).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}
	return toModelTenant(created), nil
}

func (r *TenantRepository) FindByID(ctx context.Context, id string) (*model.Tenant, error) {
	t, err := r.client.Tenant.Get(ctx, id)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find tenant by id: %w", err)
	}
	return toModelTenant(t), nil
}

func (r *TenantRepository) FindBySlug(ctx context.Context, slug string) (*model.Tenant, error) {
	t, err := r.client.Tenant.Query().
		Where(tenant.SlugEQ(slug)).
		Only(ctx)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find tenant by slug: %w", err)
	}
	return toModelTenant(t), nil
}

func toModelTenant(t *generated.Tenant) *model.Tenant {
	return &model.Tenant{
		ID:        t.ID,
		Name:      t.Name,
		Slug:      t.Slug,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
