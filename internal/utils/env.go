package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to get user home directory:", err)
	}

	envPath := filepath.Join(homeDir, ".personal-blog-generator", ".env")

	// Check if the .env file exists
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		log.Fatal(fmt.Sprintf("Configuration file not found at %s. Please create the file with your configuration settings. Example:\n\n# Database configuration\nDB_PATH=./blog.db\n\n# Server configuration\nAPP_PORT=8080", envPath))
	}

	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file from", envPath, ":", err)
	}
}
