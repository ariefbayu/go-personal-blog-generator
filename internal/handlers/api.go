package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ariefbayu/personal-blog-generator/internal/models"
	"github.com/ariefbayu/personal-blog-generator/internal/repository"
)

type APIHandlers struct {
	postRepo *repository.PostRepository
}

func NewAPIHandlers(postRepo *repository.PostRepository) *APIHandlers {
	return &APIHandlers{postRepo: postRepo}
}

func (h *APIHandlers) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postRepo.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (h *APIHandlers) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

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
