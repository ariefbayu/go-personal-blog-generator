package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/ariefbayu/personal-blog-generator/internal/models"
	"github.com/ariefbayu/personal-blog-generator/internal/repository"
)

type PageHandlers struct {
	pageRepo *repository.PageRepository
}

func NewPageHandlers(pageRepo *repository.PageRepository) *PageHandlers {
	return &PageHandlers{pageRepo: pageRepo}
}

func (h *PageHandlers) GetPagesHandler(w http.ResponseWriter, r *http.Request) {
	pages, err := h.pageRepo.GetAllPages()
	if err != nil {
		http.Error(w, "Failed to fetch pages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pages)
}

func (h *PageHandlers) GetPageHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/pages/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid page ID", http.StatusBadRequest)
		return
	}

	page, err := h.pageRepo.GetPageByID(id)
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(page)
}

func (h *PageHandlers) CreatePageHandler(w http.ResponseWriter, r *http.Request) {
	var page models.Page
	if err := json.NewDecoder(r.Body).Decode(&page); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Server-side validation
	if strings.TrimSpace(page.Title) == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(page.Slug) == "" {
		http.Error(w, "Slug is required", http.StatusBadRequest)
		return
	}

	// Validate slug format (lowercase, alphanumeric, hyphens only)
	slugRegex := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !slugRegex.MatchString(page.Slug) {
		http.Error(w, "Slug must contain only lowercase letters, numbers, and hyphens", http.StatusBadRequest)
		return
	}

	// Check if slug is unique
	existingPage, err := h.pageRepo.GetPageBySlug(page.Slug)
	if err == nil && existingPage != nil {
		http.Error(w, "Slug already exists", http.StatusConflict)
		return
	}

	err = h.pageRepo.CreatePage(&page)
	if err != nil {
		http.Error(w, "Failed to create page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": page.ID})
}

func (h *PageHandlers) UpdatePageHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/pages/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid page ID", http.StatusBadRequest)
		return
	}

	var page models.Page
	if err := json.NewDecoder(r.Body).Decode(&page); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	page.ID = id

	// Server-side validation
	if strings.TrimSpace(page.Title) == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(page.Slug) == "" {
		http.Error(w, "Slug is required", http.StatusBadRequest)
		return
	}

	// Validate slug format (lowercase, alphanumeric, hyphens only)
	slugRegex := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !slugRegex.MatchString(page.Slug) {
		http.Error(w, "Slug must contain only lowercase letters, numbers, and hyphens", http.StatusBadRequest)
		return
	}

	// Check if slug is unique (excluding current page)
	existingPage, err := h.pageRepo.GetPageBySlug(page.Slug)
	if err == nil && existingPage != nil && existingPage.ID != id {
		http.Error(w, "Slug already exists", http.StatusConflict)
		return
	}

	err = h.pageRepo.UpdatePage(&page)
	if err != nil {
		http.Error(w, "Failed to update page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Page updated successfully"})
}

func (h *PageHandlers) DeletePageHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/pages/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid page ID", http.StatusBadRequest)
		return
	}

	err = h.pageRepo.DeletePage(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			http.Error(w, "Page not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete page", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}