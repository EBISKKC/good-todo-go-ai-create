package usecase

import (
	"context"
	"testing"
	"time"

	"good-todo-go/internal/domain/model"
	"good-todo-go/internal/domain/repository/mock"
	mocku "good-todo-go/internal/pkg/mock"
	"good-todo-go/internal/usecase/input"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestTodoInteractor_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTodoRepo := mock.NewMockITodoRepository(ctrl)
	mockUUID := mocku.NewMockUUIDGenerator("test-todo-id")

	interactor := NewTodoInteractor(mockTodoRepo, mockUUID)

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		tenantID := "tenant-123"
		inp := &input.CreateTodoInput{
			Title:       "Test Todo",
			Description: "Test Description",
			IsPublic:    false,
		}

		expectedTodo := &model.Todo{
			ID:          "test-todo-id",
			TenantID:    tenantID,
			UserID:      userID,
			Title:       inp.Title,
			Description: inp.Description,
			Completed:   false,
			IsPublic:    false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockTodoRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Return(expectedTodo, nil)

		result, err := interactor.Create(ctx, userID, tenantID, inp)

		require.NoError(t, err)
		assert.Equal(t, "test-todo-id", result.ID)
		assert.Equal(t, inp.Title, result.Title)
		assert.Equal(t, inp.Description, result.Description)
		assert.False(t, result.Completed)
	})
}

func TestTodoInteractor_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTodoRepo := mock.NewMockITodoRepository(ctrl)
	mockUUID := mocku.NewMockUUIDGenerator()

	interactor := NewTodoInteractor(mockTodoRepo, mockUUID)

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"

		expectedTodos := []*model.Todo{
			{
				ID:        "todo-1",
				UserID:    userID,
				TenantID:  "tenant-123",
				Title:     "Todo 1",
				Completed: false,
			},
			{
				ID:        "todo-2",
				UserID:    userID,
				TenantID:  "tenant-123",
				Title:     "Todo 2",
				Completed: true,
			},
		}

		mockTodoRepo.EXPECT().
			FindByUserID(ctx, userID).
			Return(expectedTodos, nil)

		result, err := interactor.List(ctx, userID)

		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Todo 1", result[0].Title)
		assert.Equal(t, "Todo 2", result[1].Title)
	})

	t.Run("empty list", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-456"

		mockTodoRepo.EXPECT().
			FindByUserID(ctx, userID).
			Return([]*model.Todo{}, nil)

		result, err := interactor.List(ctx, userID)

		require.NoError(t, err)
		assert.Empty(t, result)
	})
}

func TestTodoInteractor_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTodoRepo := mock.NewMockITodoRepository(ctrl)
	mockUUID := mocku.NewMockUUIDGenerator()

	interactor := NewTodoInteractor(mockTodoRepo, mockUUID)

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		inp := &input.UpdateTodoInput{
			ID:          "todo-1",
			Title:       "Updated Todo",
			Description: "Updated Description",
			Completed:   true,
			IsPublic:    true,
		}

		existingTodo := &model.Todo{
			ID:        "todo-1",
			UserID:    userID,
			TenantID:  "tenant-123",
			Title:     "Original Todo",
			Completed: false,
		}

		updatedTodo := &model.Todo{
			ID:          "todo-1",
			UserID:      userID,
			TenantID:    "tenant-123",
			Title:       inp.Title,
			Description: inp.Description,
			Completed:   true,
			IsPublic:    true,
		}

		mockTodoRepo.EXPECT().
			FindByID(ctx, "todo-1").
			Return(existingTodo, nil)

		mockTodoRepo.EXPECT().
			Update(ctx, gomock.Any()).
			Return(updatedTodo, nil)

		result, err := interactor.Update(ctx, userID, inp)

		require.NoError(t, err)
		assert.Equal(t, inp.Title, result.Title)
		assert.True(t, result.Completed)
	})

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		inp := &input.UpdateTodoInput{
			ID:    "non-existent",
			Title: "Updated Todo",
		}

		mockTodoRepo.EXPECT().
			FindByID(ctx, "non-existent").
			Return(nil, nil)

		_, err := interactor.Update(ctx, userID, inp)

		assert.Equal(t, ErrTodoNotFound, err)
	})

	t.Run("not owner", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		inp := &input.UpdateTodoInput{
			ID:    "todo-1",
			Title: "Updated Todo",
		}

		existingTodo := &model.Todo{
			ID:       "todo-1",
			UserID:   "different-user",
			TenantID: "tenant-123",
			Title:    "Original Todo",
		}

		mockTodoRepo.EXPECT().
			FindByID(ctx, "todo-1").
			Return(existingTodo, nil)

		_, err := interactor.Update(ctx, userID, inp)

		assert.Equal(t, ErrNotTodoOwner, err)
	})
}

func TestTodoInteractor_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTodoRepo := mock.NewMockITodoRepository(ctrl)
	mockUUID := mocku.NewMockUUIDGenerator()

	interactor := NewTodoInteractor(mockTodoRepo, mockUUID)

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		todoID := "todo-1"

		existingTodo := &model.Todo{
			ID:       todoID,
			UserID:   userID,
			TenantID: "tenant-123",
			Title:    "Todo to delete",
		}

		mockTodoRepo.EXPECT().
			FindByID(ctx, todoID).
			Return(existingTodo, nil)

		mockTodoRepo.EXPECT().
			Delete(ctx, todoID).
			Return(nil)

		err := interactor.Delete(ctx, userID, todoID)

		require.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		todoID := "non-existent"

		mockTodoRepo.EXPECT().
			FindByID(ctx, todoID).
			Return(nil, nil)

		err := interactor.Delete(ctx, userID, todoID)

		assert.Equal(t, ErrTodoNotFound, err)
	})

	t.Run("not owner", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		todoID := "todo-1"

		existingTodo := &model.Todo{
			ID:       todoID,
			UserID:   "different-user",
			TenantID: "tenant-123",
			Title:    "Todo to delete",
		}

		mockTodoRepo.EXPECT().
			FindByID(ctx, todoID).
			Return(existingTodo, nil)

		err := interactor.Delete(ctx, userID, todoID)

		assert.Equal(t, ErrNotTodoOwner, err)
	})
}
