package utils

import (
"os"
"testing"
)

func TestLoadEnv(t *testing.T) {
	// Create a temporary .env file
	tempEnv := "TEST_VAR=test_value\n"
	err := os.WriteFile(".env", []byte(tempEnv), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp .env: %v", err)
	}
	defer os.Remove(".env") // clean up

	LoadEnv()

	if os.Getenv("TEST_VAR") != "test_value" {
		t.Errorf("Expected TEST_VAR to be 'test_value', got '%s'", os.Getenv("TEST_VAR"))
	}
}
