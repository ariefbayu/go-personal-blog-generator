package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ariefbayu/personal-blog-generator/internal/models"
	"github.com/ariefbayu/personal-blog-generator/internal/repository"
)

type PortfolioHandlers struct {
	portfolioRepo *repository.PortfolioRepository
}

func NewPortfolioHandlers(portfolioRepo *repository.PortfolioRepository) *PortfolioHandlers {
	return &PortfolioHandlers{portfolioRepo: portfolioRepo}
}

func (h *PortfolioHandlers) GetPortfolioItemsHandler(w http.ResponseWriter, r *http.Request) {
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

	items, total, err := h.portfolioRepo.GetPortfolioItemsPaginated(limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch portfolio items", http.StatusInternalServerError)
		return
	}

	totalPages := (total + limit - 1) / limit

	response := map[string]interface{}{
		"items":       items,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *PortfolioHandlers) GetPortfolioItemHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/portfolio/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid portfolio item ID", http.StatusBadRequest)
		return
	}

	item, err := h.portfolioRepo.GetPortfolioItemByID(id)
	if err != nil {
		http.Error(w, "Portfolio item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *PortfolioHandlers) CreatePortfolioItemHandler(w http.ResponseWriter, r *http.Request) {
	var item models.PortfolioItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Server-side validation
	if strings.TrimSpace(item.Title) == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(item.ShortDescription) == "" {
		http.Error(w, "Short description is required", http.StatusBadRequest)
		return
	}

	// Validate URLs if provided
	if item.ProjectURL != "" {
		if _, err := url.ParseRequestURI(item.ProjectURL); err != nil {
			http.Error(w, "Invalid project URL format", http.StatusBadRequest)
			return
		}
	}
	if item.GithubURL != "" {
		if _, err := url.ParseRequestURI(item.GithubURL); err != nil {
			http.Error(w, "Invalid GitHub URL format", http.StatusBadRequest)
			return
		}
	}

	err := h.portfolioRepo.CreatePortfolioItem(&item)
	if err != nil {
		http.Error(w, "Failed to create portfolio item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": item.ID})
}

func (h *PortfolioHandlers) UpdatePortfolioItemHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/portfolio/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid portfolio item ID", http.StatusBadRequest)
		return
	}

	var item models.PortfolioItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	item.ID = id

	// Server-side validation
	if strings.TrimSpace(item.Title) == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(item.ShortDescription) == "" {
		http.Error(w, "Short description is required", http.StatusBadRequest)
		return
	}

	// Validate URLs if provided
	if item.ProjectURL != "" {
		if _, err := url.ParseRequestURI(item.ProjectURL); err != nil {
			http.Error(w, "Invalid project URL format", http.StatusBadRequest)
			return
		}
	}
	if item.GithubURL != "" {
		if _, err := url.ParseRequestURI(item.GithubURL); err != nil {
			http.Error(w, "Invalid GitHub URL format", http.StatusBadRequest)
			return
		}
	}

	err = h.portfolioRepo.UpdatePortfolioItem(&item)
	if err != nil {
		http.Error(w, "Failed to update portfolio item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Portfolio item updated successfully"})
}

func (h *PortfolioHandlers) DeletePortfolioItemHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/portfolio/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid portfolio item ID", http.StatusBadRequest)
		return
	}

	err = h.portfolioRepo.DeletePortfolioItem(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			http.Error(w, "Portfolio item not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete portfolio item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
