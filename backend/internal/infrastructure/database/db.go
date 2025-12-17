package database

import (
	"context"
	"database/sql"
	"fmt"

	"good-todo-go/internal/ent/generated"
	"good-todo-go/internal/infrastructure/environment"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
)

type Database struct {
	Client *generated.Client
	DB     *sql.DB
	env    *environment.Environment
}

func NewDatabase(env *environment.Environment) (*Database, error) {
	db, err := sql.Open("postgres", env.GetAppDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	drv := entsql.OpenDB(dialect.Postgres, db)
	client := generated.NewClient(generated.Driver(drv))

	return &Database{
		Client: client,
		DB:     db,
		env:    env,
	}, nil
}

func (d *Database) Close() error {
	if err := d.Client.Close(); err != nil {
		return err
	}
	return d.DB.Close()
}

// WithTenantContext executes a function within a tenant context
func (d *Database) WithTenantContext(ctx context.Context, tenantID string, fn func(ctx context.Context, client *generated.Client) error) error {
	tx, err := d.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Set tenant context for RLS
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("SET app.current_tenant_id = '%s'", tenantID)); err != nil {
		return fmt.Errorf("failed to set tenant context: %w", err)
	}

	// Create a new client with the transaction
	drv := entsql.OpenDB(dialect.Postgres, d.DB)
	client := generated.NewClient(generated.Driver(drv))

	if err := fn(ctx, client); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
