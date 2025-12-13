package db

import (
	"database/sql"
	"fmt"
	"sort"
)

func Migrate(db *sql.DB) error {
	// Create schema_migrations table
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		id TEXT PRIMARY KEY,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// Get migration IDs from embedded migrations
	var migrationIDs []string
	for id := range Migrations {
		migrationIDs = append(migrationIDs, id)
	}
	sort.Strings(migrationIDs)

	for _, id := range migrationIDs {
		// Check if already applied
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE id = ?", id).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration %s: %w", id, err)
		}
		if count > 0 {
			continue
		}

		// Get SQL from embedded
		sql, exists := Migrations[id]
		if !exists {
			return fmt.Errorf("migration %s not found in embedded migrations", id)
		}

		// Apply in transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %s: %w", id, err)
		}

		_, err = tx.Exec(sql)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", id, err)
		}

		_, err = tx.Exec("INSERT INTO schema_migrations (id) VALUES (?)", id)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", id, err)
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", id, err)
		}
	}

	return nil
}