package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home dir: %v", err)
	}

	// Create config directory
	configDir := filepath.Join(homeDir, ".personal-blog-generator")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Create a temporary .env file in the config directory
	envPath := filepath.Join(configDir, ".env")
	tempEnv := "TEST_VAR=test_value\n"
	err = os.WriteFile(envPath, []byte(tempEnv), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp .env: %v", err)
	}
	defer os.Remove(envPath) // clean up

	LoadEnv()

	if os.Getenv("TEST_VAR") != "test_value" {
		t.Errorf("Expected TEST_VAR to be 'test_value', got '%s'", os.Getenv("TEST_VAR"))
	}
}
