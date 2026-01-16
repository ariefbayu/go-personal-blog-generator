package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ariefbayu/personal-blog-generator/internal/generator"
	"github.com/ariefbayu/personal-blog-generator/internal/models"
	"github.com/ariefbayu/personal-blog-generator/internal/repository"
)

type FileNode struct {
	Name     string     `json:"name"`
	Type     string     `json:"type"` // "file" or "dir"
	Path     string     `json:"path"`
	Editable bool       `json:"editable,omitempty"`
	Children []FileNode `json:"children,omitempty"`
}

func buildFileTree(dir string) ([]FileNode, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var nodes []FileNode
	for _, entry := range entries {
		relPath, _ := filepath.Rel(dir, filepath.Join(dir, entry.Name()))
		node := FileNode{
			Name: entry.Name(),
			Path: relPath,
		}
		if entry.IsDir() {
			node.Type = "dir"
			children, err := buildFileTree(filepath.Join(dir, entry.Name()))
			if err != nil {
				return nil, err
			}
			node.Children = children
		} else {
			node.Type = "file"
			ext := filepath.Ext(entry.Name())
			if ext == ".html" || ext == ".css" || ext == ".js" || ext == ".txt" || ext == ".md" {
				node.Editable = true
			}
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

type APIHandlers struct {
	postRepo      *repository.PostRepository
	portfolioRepo *repository.PortfolioRepository
	pageRepo      *repository.PageRepository
	settingsRepo  *repository.SettingsRepository
}

func NewAPIHandlers(postRepo *repository.PostRepository, portfolioRepo *repository.PortfolioRepository, pageRepo *repository.PageRepository, settingsRepo *repository.SettingsRepository) *APIHandlers {
	return &APIHandlers{
		postRepo:      postRepo,
		portfolioRepo: portfolioRepo,
		pageRepo:      pageRepo,
		settingsRepo:  settingsRepo,
	}
}

func (h *APIHandlers) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	posts, total, err := h.postRepo.GetPostsPaginated(limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	totalPages := (total + limit - 1) / limit

	response := map[string]interface{}{
		"posts":       posts,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *APIHandlers) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Set creation time
	post.CreatedAt = time.Now()

	err := h.postRepo.CreatePost(&post)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			http.Error(w, "Slug already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": post.ID})
}

func (h *APIHandlers) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.postRepo.GetPostByID(id)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (h *APIHandlers) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	post.ID = id

	// Set update time
	post.UpdatedAt = time.Now()

	err = h.postRepo.UpdatePost(&post)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			http.Error(w, "Slug already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post updated successfully"})
}

func (h *APIHandlers) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	err = h.postRepo.DeletePost(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *APIHandlers) PublishSiteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get paths from environment variables
	templatePath := os.Getenv("TEMPLATE_PATH")
	if templatePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			templatePath = "./templates" // fallback
		} else {
			templatePath = filepath.Join(homeDir, ".personal-blog-generator", "templates")
		}
	}
	log.Printf("DEBUG: templatePath = %s", templatePath)
	outputPath := os.Getenv("OUTPUT_PATH")
	if outputPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			outputPath = "./html-outputs" // fallback
		} else {
			outputPath = filepath.Join(homeDir, "html-outputs")
		}
	}

	// Generate the static site
	err := generator.GenerateStaticSite(h.postRepo, h.portfolioRepo, h.pageRepo, h.settingsRepo, templatePath, outputPath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Generation failed: %s", err.Error())})
		return
	}

	// Count published posts for the response
	posts, err := h.postRepo.GetAllPosts()
	if err != nil {
		// If counting fails, just return success without count
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Site generated successfully"})
		return
	}

	publishedCount := 0
	for _, post := range posts {
		if post.Published {
			publishedCount++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Site generated successfully",
		"count":   publishedCount,
	})
}

func (h *APIHandlers) GetSettingsHandler(w http.ResponseWriter, r *http.Request) {
	settings, err := h.settingsRepo.GetSettings()
	if err != nil {
		http.Error(w, "Failed to get settings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func (h *APIHandlers) UpdateSettingsHandler(w http.ResponseWriter, r *http.Request) {
	var settings repository.Settings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.settingsRepo.UpdateSettings(&settings)
	if err != nil {
		http.Error(w, "Failed to update settings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Settings updated successfully"})
}

func (h *APIHandlers) GetTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	templatePath := os.Getenv("TEMPLATE_PATH")
	if templatePath == "" {
		templatePath = "./templates"
	}
	tree, err := buildFileTree(templatePath)
	if err != nil {
		http.Error(w, "Failed to list templates", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tree)
}

func (h *APIHandlers) GetTemplateContentHandler(w http.ResponseWriter, r *http.Request) {
	pathParam := r.URL.Query().Get("path")
	if pathParam == "" {
		http.Error(w, "Missing path parameter", http.StatusBadRequest)
		return
	}
	templatePath := os.Getenv("TEMPLATE_PATH")
	if templatePath == "" {
		templatePath = "./templates"
	}
	fullPath := filepath.Join(templatePath, pathParam)
	if !strings.HasPrefix(fullPath, templatePath+string(filepath.Separator)) && fullPath != templatePath {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	content, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(content)
}

func (h *APIHandlers) SaveTemplateHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	templatePath := os.Getenv("TEMPLATE_PATH")
	if templatePath == "" {
		templatePath = "./templates"
	}
	fullPath := filepath.Join(templatePath, req.Path)
	if !strings.HasPrefix(fullPath, templatePath+string(filepath.Separator)) && fullPath != templatePath {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	// Create backup
	bakPath := fullPath + ".bak"
	if _, err := os.Stat(fullPath); err == nil {
		err = os.Rename(fullPath, bakPath)
		if err != nil {
			http.Error(w, "Failed to create backup", http.StatusInternalServerError)
			return
		}
	}
	// Write new content
	err := os.WriteFile(fullPath, []byte(req.Content), 0644)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "File saved successfully"})
}
