package database

import (
	"context"
	"database/sql"
	"fmt"

	"good-todo-go/internal/ent/generated"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
)

// TenantContext holds a database connection with tenant context set
type TenantContext struct {
	Client   *generated.Client
	tx       *sql.Tx
	tenantID string
}

// NewTenantContext creates a new tenant context with RLS set
func NewTenantContext(ctx context.Context, db *sql.DB, tenantID string) (*TenantContext, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Set tenant context for RLS
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("SET LOCAL app.current_tenant_id = '%s'", tenantID)); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to set tenant context: %w", err)
	}

	// Create Ent client with transaction
	drv := entsql.OpenDB(dialect.Postgres, &txDB{tx: tx})
	client := generated.NewClient(generated.Driver(drv))

	return &TenantContext{
		Client:   client,
		tx:       tx,
		tenantID: tenantID,
	}, nil
}

// Commit commits the transaction
func (tc *TenantContext) Commit() error {
	return tc.tx.Commit()
}

// Rollback rolls back the transaction
func (tc *TenantContext) Rollback() error {
	return tc.tx.Rollback()
}

// Close closes the tenant context (rolls back if not committed)
func (tc *TenantContext) Close() error {
	return tc.tx.Rollback()
}

// txDB wraps a transaction to implement the sql.DB interface needed by ent
type txDB struct {
	tx *sql.Tx
}

func (t *txDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *txDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

func (t *txDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

func (t *txDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	// Return the existing transaction since we're already in one
	return t.tx, nil
}

func (t *txDB) Close() error {
	return nil
}
