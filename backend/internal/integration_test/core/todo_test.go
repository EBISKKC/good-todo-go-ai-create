package core

import (
	"context"
	"testing"

	"good-todo-go/internal/domain/model"
	infrarepo "good-todo-go/internal/infrastructure/repository"
	"good-todo-go/internal/integration_test/common"
	"good-todo-go/internal/pkg"
	"good-todo-go/internal/usecase"
	"good-todo-go/internal/usecase/input"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTodoIntegration_CRUD(t *testing.T) {
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
		Role:         "member",
	})

	// Set tenant context
	err = db.SetTenantContext(ctx, tenant.ID)
	require.NoError(t, err)

	// Create repository and interactor
	todoRepo := infrarepo.NewTodoRepository(db.AppClient)
	uuidGen := pkg.NewUUIDGenerator()
	todoInteractor := usecase.NewTodoInteractor(todoRepo, uuidGen)

	var createdTodoID string

	t.Run("Create todo", func(t *testing.T) {
		inp := &input.CreateTodoInput{
			Title:       "Integration Test Todo",
			Description: "This is a test todo",
			IsPublic:    false,
		}

		result, err := todoInteractor.Create(ctx, user.ID, tenant.ID, inp)
		require.NoError(t, err)

		assert.NotEmpty(t, result.ID)
		assert.Equal(t, inp.Title, result.Title)
		assert.Equal(t, inp.Description, result.Description)
		assert.False(t, result.Completed)
		assert.False(t, result.IsPublic)

		createdTodoID = result.ID
	})

	t.Run("List todos", func(t *testing.T) {
		result, err := todoInteractor.List(ctx, user.ID)
		require.NoError(t, err)

		assert.Len(t, result, 1)
		assert.Equal(t, createdTodoID, result[0].ID)
	})

	t.Run("Update todo", func(t *testing.T) {
		inp := &input.UpdateTodoInput{
			ID:          createdTodoID,
			Title:       "Updated Todo",
			Description: "Updated description",
			Completed:   true,
			IsPublic:    true,
		}

		result, err := todoInteractor.Update(ctx, user.ID, inp)
		require.NoError(t, err)

		assert.Equal(t, inp.Title, result.Title)
		assert.Equal(t, inp.Description, result.Description)
		assert.True(t, result.Completed)
		assert.True(t, result.IsPublic)
		assert.NotNil(t, result.CompletedAt)
	})

	t.Run("List public todos", func(t *testing.T) {
		result, err := todoInteractor.ListPublic(ctx, tenant.ID)
		require.NoError(t, err)

		assert.Len(t, result, 1)
		assert.Equal(t, createdTodoID, result[0].ID)
		assert.True(t, result[0].IsPublic)
	})

	t.Run("Delete todo", func(t *testing.T) {
		err := todoInteractor.Delete(ctx, user.ID, createdTodoID)
		require.NoError(t, err)

		// Verify deletion
		result, err := todoInteractor.List(ctx, user.ID)
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	// Cleanup after test
	err = db.CleanupTables(ctx)
	require.NoError(t, err)
}

func TestTodoIntegration_Authorization(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := common.SetupTestDatabase(t)
	defer db.Close()

	ctx := context.Background()

	// Cleanup before test
	err := db.CleanupTables(ctx)
	require.NoError(t, err)

	// Create test tenant
	tenant := common.CreateTestTenant(t, db.AdminClient, &common.TestTenant{
		Name: "Test Tenant",
		Slug: "test-tenant",
	})

	// Create two users in the same tenant
	user1 := common.CreateTestUser(t, db.AdminClient, &common.TestUser{
		TenantID:     tenant.ID,
		Email:        "user1@test.com",
		PasswordHash: "hash1",
		Name:         "User 1",
		Role:         "member",
	})

	user2 := common.CreateTestUser(t, db.AdminClient, &common.TestUser{
		TenantID:     tenant.ID,
		Email:        "user2@test.com",
		PasswordHash: "hash2",
		Name:         "User 2",
		Role:         "member",
	})

	// Set tenant context
	err = db.SetTenantContext(ctx, tenant.ID)
	require.NoError(t, err)

	// Create todo for user1
	todo := common.CreateTestTodo(t, db.AdminClient, &common.TestTodo{
		TenantID: tenant.ID,
		UserID:   user1.ID,
		Title:    "User 1's Todo",
		IsPublic: false,
	})

	todoRepo := infrarepo.NewTodoRepository(db.AppClient)
	uuidGen := pkg.NewUUIDGenerator()
	todoInteractor := usecase.NewTodoInteractor(todoRepo, uuidGen)

	t.Run("User cannot update other user's todo", func(t *testing.T) {
		inp := &input.UpdateTodoInput{
			ID:        todo.ID,
			Title:     "Hacked!",
			Completed: false,
		}

		_, err := todoInteractor.Update(ctx, user2.ID, inp)
		assert.Equal(t, usecase.ErrNotTodoOwner, err)
	})

	t.Run("User cannot delete other user's todo", func(t *testing.T) {
		err := todoInteractor.Delete(ctx, user2.ID, todo.ID)
		assert.Equal(t, usecase.ErrNotTodoOwner, err)
	})

	t.Run("User can see their own todos but not others' private todos", func(t *testing.T) {
		// User1 should see their todo
		user1Todos, err := todoInteractor.List(ctx, user1.ID)
		require.NoError(t, err)
		assert.Len(t, user1Todos, 1)

		// User2 should not see user1's private todo
		user2Todos, err := todoInteractor.List(ctx, user2.ID)
		require.NoError(t, err)
		assert.Empty(t, user2Todos)
	})

	t.Run("User can see public todos from same tenant", func(t *testing.T) {
		// Make todo public
		updatedTodo := &model.Todo{
			ID:       todo.ID,
			TenantID: tenant.ID,
			UserID:   user1.ID,
			Title:    todo.Title,
			IsPublic: true,
		}
		todoRepo := infrarepo.NewTodoRepository(db.AdminClient)
		_, err := todoRepo.Update(ctx, updatedTodo)
		require.NoError(t, err)

		// Both users should see the public todo
		publicTodos, err := todoInteractor.ListPublic(ctx, tenant.ID)
		require.NoError(t, err)
		assert.Len(t, publicTodos, 1)
	})

	// Cleanup after test
	err = db.CleanupTables(ctx)
	require.NoError(t, err)
}
