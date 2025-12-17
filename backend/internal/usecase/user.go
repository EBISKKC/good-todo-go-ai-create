package usecase

import (
	"context"
	"errors"

	"good-todo-go/internal/domain/repository"
	"good-todo-go/internal/usecase/input"
	"good-todo-go/internal/usecase/output"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type IUserInteractor interface {
	GetMe(ctx context.Context, userID string) (*output.UserOutput, error)
	UpdateMe(ctx context.Context, userID string, input *input.UpdateUserInput) (*output.UserOutput, error)
}

type UserInteractor struct {
	userRepo repository.IUserRepository
}

func NewUserInteractor(userRepo repository.IUserRepository) IUserInteractor {
	return &UserInteractor{userRepo: userRepo}
}

func (i *UserInteractor) GetMe(ctx context.Context, userID string) (*output.UserOutput, error) {
	user, err := i.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return &output.UserOutput{
		ID:            user.ID,
		TenantID:      user.TenantID,
		Email:         user.Email,
		Name:          user.Name,
		Role:          string(user.Role),
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}, nil
}

func (i *UserInteractor) UpdateMe(ctx context.Context, userID string, inp *input.UpdateUserInput) (*output.UserOutput, error) {
	user, err := i.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	user.Name = inp.Name

	updated, err := i.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return &output.UserOutput{
		ID:            updated.ID,
		TenantID:      updated.TenantID,
		Email:         updated.Email,
		Name:          updated.Name,
		Role:          string(updated.Role),
		EmailVerified: updated.EmailVerified,
		CreatedAt:     updated.CreatedAt,
		UpdatedAt:     updated.UpdatedAt,
	}, nil
}
