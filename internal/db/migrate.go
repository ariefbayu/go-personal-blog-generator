package db

import (
"database/sql"
"fmt"
)

type Migration struct {
	ID  string
	SQL string
}

var migrations = []Migration{
	{
		ID: "001_initial_schema",
		SQL: `
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    content TEXT,
    tags TEXT,
    published BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS portfolio_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    short_description TEXT,
    project_url TEXT,
    github_url TEXT,
    showcase_image TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    content TEXT,
    show_in_nav BOOLEAN DEFAULT 0,
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
		`,
	},
}

func Migrate(db *sql.DB) error {
	// Create schema_migrations table
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
id TEXT PRIMARY KEY,
applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
)`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	for _, mig := range migrations {
		// Check if already applied
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE id = ?", mig.ID).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration %s: %w", mig.ID, err)
		}
		if count > 0 {
			continue
		}

		// Apply in transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %s: %w", mig.ID, err)
		}

		_, err = tx.Exec(mig.SQL)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", mig.ID, err)
		}

		_, err = tx.Exec("INSERT INTO schema_migrations (id) VALUES (?)", mig.ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", mig.ID, err)
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", mig.ID, err)
		}
	}

	return nil
}
