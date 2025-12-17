package usecase

import (
	"context"
	"testing"
	"time"

	"good-todo-go/internal/domain/model"
	"good-todo-go/internal/domain/repository/mock"
	"good-todo-go/internal/usecase/input"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestUserInteractor_GetMe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockIUserRepository(ctrl)

	interactor := NewUserInteractor(mockUserRepo)

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"

		expectedUser := &model.User{
			ID:            userID,
			TenantID:      "tenant-123",
			Email:         "test@example.com",
			Name:          "Test User",
			Role:          model.UserRoleAdmin,
			EmailVerified: true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		mockUserRepo.EXPECT().
			FindByID(ctx, userID).
			Return(expectedUser, nil)

		result, err := interactor.GetMe(ctx, userID)

		require.NoError(t, err)
		assert.Equal(t, userID, result.ID)
		assert.Equal(t, "test@example.com", result.Email)
		assert.Equal(t, "Test User", result.Name)
		assert.Equal(t, "admin", result.Role)
		assert.True(t, result.EmailVerified)
	})

	t.Run("user not found", func(t *testing.T) {
		ctx := context.Background()
		userID := "non-existent"

		mockUserRepo.EXPECT().
			FindByID(ctx, userID).
			Return(nil, nil)

		_, err := interactor.GetMe(ctx, userID)

		assert.Equal(t, ErrUserNotFound, err)
	})
}

func TestUserInteractor_UpdateMe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockIUserRepository(ctrl)

	interactor := NewUserInteractor(mockUserRepo)

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		inp := &input.UpdateUserInput{
			Name: "Updated Name",
		}

		existingUser := &model.User{
			ID:            userID,
			TenantID:      "tenant-123",
			Email:         "test@example.com",
			Name:          "Original Name",
			Role:          model.UserRoleMember,
			EmailVerified: true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		updatedUser := &model.User{
			ID:            userID,
			TenantID:      "tenant-123",
			Email:         "test@example.com",
			Name:          "Updated Name",
			Role:          model.UserRoleMember,
			EmailVerified: true,
			CreatedAt:     existingUser.CreatedAt,
			UpdatedAt:     time.Now(),
		}

		mockUserRepo.EXPECT().
			FindByID(ctx, userID).
			Return(existingUser, nil)

		mockUserRepo.EXPECT().
			Update(ctx, gomock.Any()).
			Return(updatedUser, nil)

		result, err := interactor.UpdateMe(ctx, userID, inp)

		require.NoError(t, err)
		assert.Equal(t, "Updated Name", result.Name)
	})

	t.Run("user not found", func(t *testing.T) {
		ctx := context.Background()
		userID := "non-existent"
		inp := &input.UpdateUserInput{
			Name: "Updated Name",
		}

		mockUserRepo.EXPECT().
			FindByID(ctx, userID).
			Return(nil, nil)

		_, err := interactor.UpdateMe(ctx, userID, inp)

		assert.Equal(t, ErrUserNotFound, err)
	})
}
