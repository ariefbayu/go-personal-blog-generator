package main

import (
	"log"
	"os"

	"github.com/ariefbayu/personal-blog-generator/internal/db"
	"github.com/ariefbayu/personal-blog-generator/internal/utils"
)

func main() {
	utils.LoadEnv()

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./blog.db" // default
	}

	database, err := db.Connect(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	err = db.Migrate(database)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connected and migrated successfully")
}