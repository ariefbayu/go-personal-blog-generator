package db

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConnect(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	db, err := Connect(dbPath)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer db.Close()

	// Run migrations to create tables
	err = Migrate(db)
	if err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	// Verify tables were created
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name IN ('posts', 'portfolio_items', 'pages')").Scan(&count)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if count != 3 {
		t.Errorf("Expected 3 tables, got %d", count)
	}

	// Check if file was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}
}