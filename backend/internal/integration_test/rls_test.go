package integration_test

import (
	"context"
	"testing"

	"good-todo-go/internal/ent/generated/todo"
	"good-todo-go/internal/ent/generated/user"
	"good-todo-go/internal/integration_test/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRLS_TenantIsolation tests that RLS properly isolates data between tenants
func TestRLS_TenantIsolation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := common.SetupTestDatabase(t)
	defer db.Close()

	ctx := context.Background()

	// Cleanup before test
	err := db.CleanupTables(ctx)
	require.NoError(t, err, "Failed to cleanup tables")

	// Create two tenants
	tenant1 := common.CreateTestTenant(t, db.AdminClient, &common.TestTenant{
		Name: "Tenant 1",
		Slug: "tenant-1",
	})

	tenant2 := common.CreateTestTenant(t, db.AdminClient, &common.TestTenant{
		Name: "Tenant 2",
		Slug: "tenant-2",
	})

	// Create users for each tenant
	user1 := common.CreateTestUser(t, db.AdminClient, &common.TestUser{
		TenantID:     tenant1.ID,
		Email:        "user1@tenant1.com",
		PasswordHash: "hash1",
		Name:         "User 1",
		Role:         "admin",
	})

	user2 := common.CreateTestUser(t, db.AdminClient, &common.TestUser{
		TenantID:     tenant2.ID,
		Email:        "user2@tenant2.com",
		PasswordHash: "hash2",
		Name:         "User 2",
		Role:         "admin",
	})

	// Create todos for each tenant
	todo1 := common.CreateTestTodo(t, db.AdminClient, &common.TestTodo{
		TenantID: tenant1.ID,
		UserID:   user1.ID,
		Title:    "Tenant 1 Todo",
		IsPublic: false,
	})

	todo2 := common.CreateTestTodo(t, db.AdminClient, &common.TestTodo{
		TenantID: tenant2.ID,
		UserID:   user2.ID,
		Title:    "Tenant 2 Todo",
		IsPublic: false,
	})

	t.Run("Users can only see their tenant's users", func(t *testing.T) {
		// Set tenant context to tenant 1
		err := db.SetTenantContext(ctx, tenant1.ID)
		require.NoError(t, err)

		// Query users - should only see tenant 1's users
		users, err := db.AppClient.User.Query().All(ctx)
		require.NoError(t, err)

		assert.Len(t, users, 1, "Should only see 1 user in tenant 1")
		assert.Equal(t, user1.Email, users[0].Email)

		// Set tenant context to tenant 2
		err = db.SetTenantContext(ctx, tenant2.ID)
		require.NoError(t, err)

		// Query users - should only see tenant 2's users
		users, err = db.AppClient.User.Query().All(ctx)
		require.NoError(t, err)

		assert.Len(t, users, 1, "Should only see 1 user in tenant 2")
		assert.Equal(t, user2.Email, users[0].Email)
	})

	t.Run("Users can only see their tenant's todos", func(t *testing.T) {
		// Set tenant context to tenant 1
		err := db.SetTenantContext(ctx, tenant1.ID)
		require.NoError(t, err)

		// Query todos - should only see tenant 1's todos
		todos, err := db.AppClient.Todo.Query().All(ctx)
		require.NoError(t, err)

		assert.Len(t, todos, 1, "Should only see 1 todo in tenant 1")
		assert.Equal(t, todo1.Title, todos[0].Title)

		// Set tenant context to tenant 2
		err = db.SetTenantContext(ctx, tenant2.ID)
		require.NoError(t, err)

		// Query todos - should only see tenant 2's todos
		todos, err = db.AppClient.Todo.Query().All(ctx)
		require.NoError(t, err)

		assert.Len(t, todos, 1, "Should only see 1 todo in tenant 2")
		assert.Equal(t, todo2.Title, todos[0].Title)
	})

	t.Run("Cannot access other tenant's user by ID", func(t *testing.T) {
		// Set tenant context to tenant 1
		err := db.SetTenantContext(ctx, tenant1.ID)
		require.NoError(t, err)

		// Try to get tenant 2's user - should not find it due to RLS
		u, err := db.AppClient.User.Query().
			Where(user.IDEQ(user2.ID)).
			Only(ctx)

		assert.Nil(t, u, "Should not be able to access other tenant's user")
		assert.Error(t, err, "Should return error when accessing other tenant's user")
	})

	t.Run("Cannot access other tenant's todo by ID", func(t *testing.T) {
		// Set tenant context to tenant 1
		err := db.SetTenantContext(ctx, tenant1.ID)
		require.NoError(t, err)

		// Try to get tenant 2's todo - should not find it due to RLS
		td, err := db.AppClient.Todo.Query().
			Where(todo.IDEQ(todo2.ID)).
			Only(ctx)

		assert.Nil(t, td, "Should not be able to access other tenant's todo")
		assert.Error(t, err, "Should return error when accessing other tenant's todo")
	})

	// Cleanup after test
	err = db.CleanupTables(ctx)
	require.NoError(t, err, "Failed to cleanup tables after test")
}

// TestRLS_EmptyTenantContext tests behavior when tenant context is not set
func TestRLS_EmptyTenantContext(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := common.SetupTestDatabase(t)
	defer db.Close()

	ctx := context.Background()

	// Cleanup before test
	err := db.CleanupTables(ctx)
	require.NoError(t, err, "Failed to cleanup tables")

	// Create a tenant and user
	tenant := common.CreateTestTenant(t, db.AdminClient, &common.TestTenant{
		Name: "Test Tenant",
		Slug: "test-tenant",
	})

	common.CreateTestUser(t, db.AdminClient, &common.TestUser{
		TenantID:     tenant.ID,
		Email:        "user@test.com",
		PasswordHash: "hash",
		Name:         "Test User",
		Role:         "admin",
	})

	common.CreateTestTodo(t, db.AdminClient, &common.TestTodo{
		TenantID: tenant.ID,
		UserID:   "user-id",
		Title:    "Test Todo",
		IsPublic: false,
	})

	t.Run("Users table allows empty tenant context for email verification", func(t *testing.T) {
		// Clear tenant context
		err := db.ClearTenantContext(ctx)
		require.NoError(t, err)

		// Should be able to query users with empty tenant context
		// (This is needed for email verification flow)
		users, err := db.AppClient.User.Query().All(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, users, "Should be able to see users with empty tenant context")
	})

	t.Run("Todos table requires tenant context", func(t *testing.T) {
		// Clear tenant context
		err := db.ClearTenantContext(ctx)
		require.NoError(t, err)

		// Should not be able to see todos without tenant context
		todos, err := db.AppClient.Todo.Query().All(ctx)
		require.NoError(t, err)
		assert.Empty(t, todos, "Should not see todos without tenant context")
	})

	// Cleanup after test
	err = db.CleanupTables(ctx)
	require.NoError(t, err, "Failed to cleanup tables after test")
}
