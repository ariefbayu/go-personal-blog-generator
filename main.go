package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ariefbayu/personal-blog-generator/internal/db"
	"github.com/ariefbayu/personal-blog-generator/internal/handlers"
	"github.com/ariefbayu/personal-blog-generator/internal/repository"
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

	postRepo := repository.NewPostRepository(database)
	apiHandlers := handlers.NewAPIHandlers(postRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api/posts", apiHandlers.GetPostsHandler)
	r.Get("/admin/dashboard", handlers.ServeDashboard)
	r.Get("/admin/posts", handlers.ServePostsPage)
	r.Handle("/admin/*", http.StripPrefix("/admin/", http.FileServer(http.Dir("admin-files/"))))

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}