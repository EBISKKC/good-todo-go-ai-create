package common

import (
	"context"
	"testing"

	"good-todo-go/internal/domain/model"
	"good-todo-go/internal/ent/generated"
	"good-todo-go/internal/ent/generated/user"

	"github.com/google/uuid"
)

// TestTenant represents test tenant data
type TestTenant struct {
	ID   string
	Name string
	Slug string
}

// TestUser represents test user data
type TestUser struct {
	ID           string
	TenantID     string
	Email        string
	PasswordHash string
	Name         string
	Role         string
}

// TestTodo represents test todo data
type TestTodo struct {
	ID       string
	TenantID string
	UserID   string
	Title    string
	IsPublic bool
}

// CreateTestTenant creates a test tenant using admin client
func CreateTestTenant(t *testing.T, client *generated.Client, tenant *TestTenant) *model.Tenant {
	t.Helper()

	if tenant.ID == "" {
		tenant.ID = uuid.New().String()
	}

	created, err := client.Tenant.Create().
		SetID(tenant.ID).
		SetName(tenant.Name).
		SetSlug(tenant.Slug).
		Save(context.Background())
	if err != nil {
		t.Fatalf("Failed to create test tenant: %v", err)
	}

	return &model.Tenant{
		ID:        created.ID,
		Name:      created.Name,
		Slug:      created.Slug,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
	}
}

// CreateTestUser creates a test user using admin client
func CreateTestUser(t *testing.T, client *generated.Client, u *TestUser) *model.User {
	t.Helper()

	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	role := user.RoleMember
	if u.Role == "admin" {
		role = user.RoleAdmin
	}

	created, err := client.User.Create().
		SetID(u.ID).
		SetTenantID(u.TenantID).
		SetEmail(u.Email).
		SetPasswordHash(u.PasswordHash).
		SetName(u.Name).
		SetRole(role).
		SetEmailVerified(true).
		Save(context.Background())
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return &model.User{
		ID:            created.ID,
		TenantID:      created.TenantID,
		Email:         created.Email,
		PasswordHash:  created.PasswordHash,
		Name:          created.Name,
		Role:          model.UserRole(created.Role),
		EmailVerified: created.EmailVerified,
		CreatedAt:     created.CreatedAt,
		UpdatedAt:     created.UpdatedAt,
	}
}

// CreateTestTodo creates a test todo using admin client
func CreateTestTodo(t *testing.T, client *generated.Client, todo *TestTodo) *model.Todo {
	t.Helper()

	if todo.ID == "" {
		todo.ID = uuid.New().String()
	}

	created, err := client.Todo.Create().
		SetID(todo.ID).
		SetTenantID(todo.TenantID).
		SetUserID(todo.UserID).
		SetTitle(todo.Title).
		SetIsPublic(todo.IsPublic).
		Save(context.Background())
	if err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	return &model.Todo{
		ID:        created.ID,
		TenantID:  created.TenantID,
		UserID:    created.UserID,
		Title:     created.Title,
		IsPublic:  created.IsPublic,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
	}
}
