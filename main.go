package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ariefbayu/personal-blog-generator/internal/db"
	"github.com/ariefbayu/personal-blog-generator/internal/handlers"
	"github.com/ariefbayu/personal-blog-generator/internal/repository"
	"github.com/ariefbayu/personal-blog-generator/internal/utils"
)

//go:embed internal/handlers/assets/admin/**
var adminFS embed.FS

func main() {
	utils.LoadEnv()

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			dbPath = "./blog.db" // fallback
		} else {
			dbPath = filepath.Join(homeDir, ".personal-blog-generator", "blog.db")
		}
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
	portfolioRepo := repository.NewPortfolioRepository(database)
	pageRepo := repository.NewPageRepository(database)
	settingsRepo := repository.NewSettingsRepository(database)
	apiHandlers := handlers.NewAPIHandlers(postRepo, portfolioRepo, pageRepo, settingsRepo)
	portfolioHandlers := handlers.NewPortfolioHandlers(portfolioRepo)
	pageHandlers := handlers.NewPageHandlers(pageRepo)

	// Create sub-filesystem to strip internal/handlers/assets/admin prefix
	adminSubFS, err := fs.Sub(adminFS, "internal/handlers/assets/admin")
	if err != nil {
		log.Fatal(err)
	}
	handlers.AdminFS = adminSubFS
	handlers.DBPath = dbPath

	// Set TemplatePath and OutputPath for handlers
	templatePath := os.Getenv("TEMPLATE_PATH")
	if templatePath == "" {
		homeDir, _ := os.UserHomeDir()
		templatePath = filepath.Join(homeDir, ".personal-blog-generator", "templates")
	}
	handlers.TemplatePath = templatePath

	outputPath := os.Getenv("OUTPUT_PATH")
	if outputPath == "" {
		homeDir, _ := os.UserHomeDir()
		outputPath = filepath.Join(homeDir, "html-outputs")
	}
	handlers.OutputPath = outputPath

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// ROOT_PREFIX logic
	rootPrefix := os.Getenv("ROOT_PREFIX")
	if rootPrefix == "" {
		rootPrefix = "/admin"
	}
	if rootPrefix == "/" {
		rootPrefix = ""
	} else {
		if !strings.HasPrefix(rootPrefix, "/") {
			rootPrefix = "/" + rootPrefix
		}
		rootPrefix = strings.TrimSuffix(rootPrefix, "/")
	}
	handlers.RootPrefix = rootPrefix

	// Public routes
	r.Handle("/images/*", http.StripPrefix("/images/", http.FileServer(http.Dir(filepath.Join(outputPath, "images")))))

	// Protected routes (Admin & API) under ROOT_PREFIX
	r.Route(rootPrefix, func(r chi.Router) {
	        adminUser := os.Getenv("ADMIN_USERNAME")
	        adminPass := os.Getenv("ADMIN_PASSWORD")
	        if adminUser == "" {
	                adminUser = "admin"
	        }
	        if adminPass == "" {
	                adminPass = "admin"
	        }

	        r.Use(middleware.BasicAuth("Blog Admin Area", map[string]string{
	                adminUser: adminPass,
	        }))

	        // API routes
	        r.Get("/api/posts", apiHandlers.GetPostsHandler)
	        r.Post("/api/posts", apiHandlers.CreatePostHandler)
	        r.Get("/api/posts/{id}", apiHandlers.GetPostHandler)
	        r.Put("/api/posts/{id}", apiHandlers.UpdatePostHandler)
	        r.Delete("/api/posts/{id}", apiHandlers.DeletePostHandler)
	        r.Get("/api/portfolio", portfolioHandlers.GetPortfolioItemsHandler)
	        r.Post("/api/portfolio", portfolioHandlers.CreatePortfolioItemHandler)
	        r.Get("/api/portfolio/{id}", portfolioHandlers.GetPortfolioItemHandler)
	        r.Put("/api/portfolio/{id}", portfolioHandlers.UpdatePortfolioItemHandler)
	        r.Delete("/api/portfolio/{id}", portfolioHandlers.DeletePortfolioItemHandler)
	        r.Get("/api/pages", pageHandlers.GetPagesHandler)
	        r.Post("/api/pages", pageHandlers.CreatePageHandler)
	        r.Get("/api/pages/{id}", pageHandlers.GetPageHandler)
	        r.Put("/api/pages/{id}", pageHandlers.UpdatePageHandler)
	        r.Delete("/api/pages/{id}", pageHandlers.DeletePageHandler)
	        r.Get("/api/settings", apiHandlers.GetSettingsHandler)
	        r.Post("/api/settings", apiHandlers.UpdateSettingsHandler)
	        r.Get("/api/settings/templates", apiHandlers.GetTemplatesHandler)
	        r.Get("/api/settings/templates/content", apiHandlers.GetTemplateContentHandler)
	        r.Post("/api/settings/templates/save", apiHandlers.SaveTemplateHandler)
	        r.Post("/api/upload/image", handlers.UploadImageHandler)
	        r.Post("/api/publish", apiHandlers.PublishSiteHandler)

	        // Admin root redirects (must come before static assets)
	        r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	                http.Redirect(w, r, rootPrefix+"/dashboard", http.StatusFound)
	        })

	        // Admin page routes
	        r.Get("/dashboard", handlers.ServeDashboard)
	        r.Get("/posts", handlers.ServePostsPage)
	        r.Get("/posts/new", handlers.ServeNewPostPage)
	        r.Get("/posts/{id}/edit", handlers.ServeEditPostPage)
	        r.Get("/portfolio", handlers.ServePortfolioPage)
	        r.Get("/portfolio/new", handlers.ServeNewPortfolioPage)
	        r.Get("/portfolio/{id}/edit", handlers.ServeEditPortfolioPage)
	        r.Get("/pages", handlers.ServePagesPage)
	        r.Get("/pages/new", handlers.ServeNewPagePage)
	        r.Get("/pages/{id}/edit", handlers.ServeEditPagePage)
	        r.Get("/settings", handlers.ServeSettingsPage)
	        r.Get("/templates", handlers.ServeTemplatesPage)

	        // Admin static assets (must come after specific routes to avoid catching them)
	        r.Handle("/*", http.StripPrefix(rootPrefix+"/", http.FileServer(http.FS(handlers.AdminFS))))
	})
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
