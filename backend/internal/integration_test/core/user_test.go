package core

import (
	"context"
	"testing"

	infrarepo "good-todo-go/internal/infrastructure/repository"
	"good-todo-go/internal/integration_test/common"
	"good-todo-go/internal/usecase"
	"good-todo-go/internal/usecase/input"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserIntegration_GetMe(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := common.SetupTestDatabase(t)
	defer db.Close()

	ctx := context.Background()

	// Cleanup before test
	err := db.CleanupTables(ctx)
	require.NoError(t, err)

	// Create test tenant and user
	tenant := common.CreateTestTenant(t, db.AdminClient, &common.TestTenant{
		Name: "Test Tenant",
		Slug: "test-tenant",
	})

	user := common.CreateTestUser(t, db.AdminClient, &common.TestUser{
		TenantID:     tenant.ID,
		Email:        "user@test.com",
		PasswordHash: "hash",
		Name:         "Test User",
		Role:         "admin",
	})

	// Set tenant context
	err = db.SetTenantContext(ctx, tenant.ID)
	require.NoError(t, err)

	// Create repository and interactor
	userRepo := infrarepo.NewUserRepository(db.AppClient)
	userInteractor := usecase.NewUserInteractor(userRepo)

	t.Run("Get existing user", func(t *testing.T) {
		result, err := userInteractor.GetMe(ctx, user.ID)
		require.NoError(t, err)

		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, user.Email, result.Email)
		assert.Equal(t, user.Name, result.Name)
		assert.Equal(t, "admin", result.Role)
		assert.True(t, result.EmailVerified)
	})

	t.Run("Get non-existent user", func(t *testing.T) {
		_, err := userInteractor.GetMe(ctx, "non-existent-id")
		assert.Equal(t, usecase.ErrUserNotFound, err)
	})

	// Cleanup after test
	err = db.CleanupTables(ctx)
	require.NoError(t, err)
}

func TestUserIntegration_UpdateMe(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := common.SetupTestDatabase(t)
	defer db.Close()

	ctx := context.Background()

	// Cleanup before test
	err := db.CleanupTables(ctx)
	require.NoError(t, err)

	// Create test tenant and user
	tenant := common.CreateTestTenant(t, db.AdminClient, &common.TestTenant{
		Name: "Test Tenant",
		Slug: "test-tenant",
	})

	user := common.CreateTestUser(t, db.AdminClient, &common.TestUser{
		TenantID:     tenant.ID,
		Email:        "user@test.com",
		PasswordHash: "hash",
		Name:         "Original Name",
		Role:         "member",
	})

	// Set tenant context
	err = db.SetTenantContext(ctx, tenant.ID)
	require.NoError(t, err)

	// Create repository and interactor
	userRepo := infrarepo.NewUserRepository(db.AppClient)
	userInteractor := usecase.NewUserInteractor(userRepo)

	t.Run("Update user name", func(t *testing.T) {
		inp := &input.UpdateUserInput{
			Name: "Updated Name",
		}

		result, err := userInteractor.UpdateMe(ctx, user.ID, inp)
		require.NoError(t, err)

		assert.Equal(t, "Updated Name", result.Name)
		assert.Equal(t, user.Email, result.Email) // Email should not change
	})

	t.Run("Update non-existent user", func(t *testing.T) {
		inp := &input.UpdateUserInput{
			Name: "New Name",
		}

		_, err := userInteractor.UpdateMe(ctx, "non-existent-id", inp)
		assert.Equal(t, usecase.ErrUserNotFound, err)
	})

	// Cleanup after test
	err = db.CleanupTables(ctx)
	require.NoError(t, err)
}
