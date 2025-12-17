package common

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"good-todo-go/internal/ent/generated"
	"good-todo-go/internal/infrastructure/environment"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
)

// TestDatabase holds database connections for testing
type TestDatabase struct {
	AdminClient *generated.Client
	AdminDB     *sql.DB
	AppClient   *generated.Client
	AppDB       *sql.DB
	env         *environment.Environment
}

// SetupTestDatabase creates database connections for testing
// Requires PostgreSQL to be running (via docker compose)
func SetupTestDatabase(t *testing.T) *TestDatabase {
	t.Helper()

	env := environment.NewEnvironment()

	// Check if test should be skipped
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration test")
	}

	// Admin connection (without RLS)
	adminDB, err := sql.Open("postgres", env.GetAdminDSN())
	if err != nil {
		t.Skipf("Failed to connect to admin database (is PostgreSQL running?): %v", err)
	}

	if err := adminDB.Ping(); err != nil {
		t.Skipf("Failed to ping admin database: %v", err)
	}

	adminDrv := entsql.OpenDB(dialect.Postgres, adminDB)
	adminClient := generated.NewClient(generated.Driver(adminDrv))

	// App connection (with RLS)
	appDB, err := sql.Open("postgres", env.GetAppDSN())
	if err != nil {
		t.Skipf("Failed to connect to app database: %v", err)
	}

	if err := appDB.Ping(); err != nil {
		t.Skipf("Failed to ping app database: %v", err)
	}

	appDrv := entsql.OpenDB(dialect.Postgres, appDB)
	appClient := generated.NewClient(generated.Driver(appDrv))

	return &TestDatabase{
		AdminClient: adminClient,
		AdminDB:     adminDB,
		AppClient:   appClient,
		AppDB:       appDB,
		env:         env,
	}
}

// Close closes all database connections
func (td *TestDatabase) Close() {
	td.AdminClient.Close()
	td.AdminDB.Close()
	td.AppClient.Close()
	td.AppDB.Close()
}

// CleanupTables removes all data from test tables
func (td *TestDatabase) CleanupTables(ctx context.Context) error {
	// Delete in correct order due to foreign keys
	if _, err := td.AdminDB.ExecContext(ctx, "DELETE FROM todos"); err != nil {
		return fmt.Errorf("failed to clean todos: %w", err)
	}
	if _, err := td.AdminDB.ExecContext(ctx, "DELETE FROM users"); err != nil {
		return fmt.Errorf("failed to clean users: %w", err)
	}
	if _, err := td.AdminDB.ExecContext(ctx, "DELETE FROM tenants"); err != nil {
		return fmt.Errorf("failed to clean tenants: %w", err)
	}
	return nil
}

// SetTenantContext sets the tenant context for RLS
func (td *TestDatabase) SetTenantContext(ctx context.Context, tenantID string) error {
	_, err := td.AppDB.ExecContext(ctx, fmt.Sprintf("SET app.current_tenant_id = '%s'", tenantID))
	return err
}

// ClearTenantContext clears the tenant context
func (td *TestDatabase) ClearTenantContext(ctx context.Context) error {
	_, err := td.AppDB.ExecContext(ctx, "RESET app.current_tenant_id")
	return err
}
