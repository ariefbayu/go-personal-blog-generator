package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

	// Get migration files
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}
	migrationDir := filepath.Join(wd, "internal", "db", "migrations")
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	for _, fileName := range migrationFiles {
		id := strings.TrimSuffix(fileName, ".sql")

		// Check if already applied
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE id = ?", id).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration %s: %w", id, err)
		}
		if count > 0 {
			continue
		}

		// Read SQL
		sqlPath := filepath.Join(migrationDir, fileName)
		sqlBytes, err := os.ReadFile(sqlPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", fileName, err)
		}
		sql := string(sqlBytes)

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