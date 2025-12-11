package handlers

import (
"encoding/json"
"net/http"

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
