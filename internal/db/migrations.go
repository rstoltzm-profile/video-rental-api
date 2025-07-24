package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Helper to apply a migration if not already applied
func applyMigration(pool *pgxpool.Pool, version string, sql string) error {
	var alreadyApplied bool
	err := pool.QueryRow(context.Background(),
		`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)`, version).Scan(&alreadyApplied)
	if err != nil {
		return fmt.Errorf("failed to check migration version %s: %w", version, err)
	}
	if alreadyApplied {
		log.Printf("Migration %s already applied.", version)
		return nil
	}
	_, err = pool.Exec(context.Background(), sql)
	if err != nil {
		return fmt.Errorf("failed to apply migration %s: %w", version, err)
	}
	_, err = pool.Exec(context.Background(),
		`INSERT INTO schema_migrations (version) VALUES ($1)`, version)
	if err != nil {
		return fmt.Errorf("failed to record migration version %s: %w", version, err)
	}
	log.Printf("Migration %s applied successfully.", version)
	return nil
}

func RunMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()

	// 1. Ensure migrations table exists
	_, err := pool.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version VARCHAR(255) PRIMARY KEY,
            applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// 2. Ensure Pagila schema is present
	var exists bool
	err = pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name = 'customer')`).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check schema: %w", err)
	}
	if !exists {
		return fmt.Errorf("Pagila schema not found - please run docker-compose up in pagila directory first")
	}

	// 3. Mark initial schema as migrated
	if err := applyMigration(pool, "pagila-initial", `SELECT 1`); err != nil {
		return err
	}

	// 4. Example: Add a new migration for a placeholder table
	if err := applyMigration(pool, "2025-07-23-placeholder-table", `
        CREATE TABLE IF NOT EXISTS placeholder (
            id SERIAL PRIMARY KEY,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `); err != nil {
		return err
	}

	// 5. Delete Table
	if err := applyMigration(pool, "2025-07-23-drop-placeholder-table", `
		DROP TABLE IF EXISTS placeholder
	`); err != nil {
		return err
	}

	return nil
}

func GetCurrentMigration(pool *pgxpool.Pool) (string, string, error) {
	var version string
	var appliedAt time.Time
	err := pool.QueryRow(context.Background(),
		`SELECT version, applied_at FROM schema_migrations ORDER BY applied_at DESC LIMIT 1`,
	).Scan(&version, &appliedAt)
	if err != nil {
		return "", "", err
	}
	return version, appliedAt.Format(time.RFC3339), nil
}
