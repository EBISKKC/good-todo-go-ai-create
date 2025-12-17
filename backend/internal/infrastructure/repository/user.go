package repository

import (
	"context"
	"fmt"

	"good-todo-go/internal/domain/model"
	"good-todo-go/internal/domain/repository"
	"good-todo-go/internal/ent/generated"
	"good-todo-go/internal/ent/generated/user"
)

type UserRepository struct {
	client *generated.Client
}

func NewUserRepository(client *generated.Client) repository.IUserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) Create(ctx context.Context, u *model.User) (*model.User, error) {
	builder := r.client.User.Create().
		SetID(u.ID).
		SetTenantID(u.TenantID).
		SetEmail(u.Email).
		SetPasswordHash(u.PasswordHash).
		SetName(u.Name).
		SetRole(user.Role(u.Role)).
		SetEmailVerified(u.EmailVerified)

	if u.VerificationToken != nil {
		builder.SetVerificationToken(*u.VerificationToken)
	}
	if u.VerificationTokenExpiresAt != nil {
		builder.SetVerificationTokenExpiresAt(*u.VerificationTokenExpiresAt)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return toModelUser(created), nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	u, err := r.client.User.Get(ctx, id)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}
	return toModelUser(u), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, tenantID, email string) (*model.User, error) {
	u, err := r.client.User.Query().
		Where(
			user.TenantIDEQ(tenantID),
			user.EmailEQ(email),
		).
		Only(ctx)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return toModelUser(u), nil
}

func (r *UserRepository) FindByVerificationToken(ctx context.Context, token string) (*model.User, error) {
	u, err := r.client.User.Query().
		Where(user.VerificationTokenEQ(token)).
		Only(ctx)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by verification token: %w", err)
	}
	return toModelUser(u), nil
}

func (r *UserRepository) Update(ctx context.Context, u *model.User) (*model.User, error) {
	builder := r.client.User.UpdateOneID(u.ID).
		SetEmail(u.Email).
		SetName(u.Name).
		SetRole(user.Role(u.Role)).
		SetEmailVerified(u.EmailVerified)

	if u.VerificationToken != nil {
		builder.SetVerificationToken(*u.VerificationToken)
	} else {
		builder.ClearVerificationToken()
	}
	if u.VerificationTokenExpiresAt != nil {
		builder.SetVerificationTokenExpiresAt(*u.VerificationTokenExpiresAt)
	} else {
		builder.ClearVerificationTokenExpiresAt()
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return toModelUser(updated), nil
}

func toModelUser(u *generated.User) *model.User {
	return &model.User{
		ID:                         u.ID,
		TenantID:                   u.TenantID,
		Email:                      u.Email,
		PasswordHash:               u.PasswordHash,
		Name:                       u.Name,
		Role:                       model.UserRole(u.Role),
		EmailVerified:              u.EmailVerified,
		VerificationToken:          u.VerificationToken,
		VerificationTokenExpiresAt: u.VerificationTokenExpiresAt,
		CreatedAt:                  u.CreatedAt,
		UpdatedAt:                  u.UpdatedAt,
	}
}
