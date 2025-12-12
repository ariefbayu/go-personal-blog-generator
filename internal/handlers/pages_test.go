package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ariefbayu/personal-blog-generator/internal/models"
	"github.com/ariefbayu/personal-blog-generator/internal/repository"
	"github.com/go-chi/chi/v5"
)

func TestPageHandlers(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer db.Close()

	// Create repository and handlers
	pageRepo := repository.NewPageRepository(db)
	handlers := NewPageHandlers(pageRepo)

	// Test data
	testPage := models.Page{
		Title:     "Test Page",
		Slug:      "test-page",
		Content:   "# Test Content\n\nThis is a test page.",
		ShowInNav: true,
		SortOrder: 1,
	}

	t.Run("CreatePage", func(t *testing.T) {
		data, _ := json.Marshal(testPage)
		req := httptest.NewRequest("POST", "/api/pages", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handlers.CreatePageHandler(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", w.Code)
		}

		var response map[string]int64
		json.Unmarshal(w.Body.Bytes(), &response)
		if response["id"] == 0 {
			t.Error("Expected non-zero ID in response")
		}
		testPage.ID = response["id"]
	})

	t.Run("GetPages", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/pages", nil)
		w := httptest.NewRecorder()

		handlers.GetPagesHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to decode JSON: %v", err)
		}

		pagesInterface, ok := response["pages"]
		if !ok {
			t.Fatal("Response missing 'pages' field")
		}

		pagesData, err := json.Marshal(pagesInterface)
		if err != nil {
			t.Fatalf("Failed to marshal pages: %v", err)
		}

		var pages []models.Page
		err = json.Unmarshal(pagesData, &pages)
		if err != nil {
			t.Fatalf("Failed to unmarshal pages: %v", err)
		}

		if len(pages) != 1 {
			t.Errorf("Expected 1 page, got %d", len(pages))
		}
		if pages[0].Title != testPage.Title {
			t.Errorf("Expected title %s, got %s", testPage.Title, pages[0].Title)
		}

		// Check pagination metadata
		if response["total"].(float64) != 1 {
			t.Errorf("Expected total 1, got %v", response["total"])
		}

		if response["page"].(float64) != 1 {
			t.Errorf("Expected page 1, got %v", response["page"])
		}

		if response["limit"].(float64) != 10 {
			t.Errorf("Expected limit 10, got %v", response["limit"])
		}

		if response["total_pages"].(float64) != 1 {
			t.Errorf("Expected total_pages 1, got %v", response["total_pages"])
		}
	})

	t.Run("GetPage", func(t *testing.T) {
		r := chi.NewRouter()
		r.Get("/api/pages/{id}", handlers.GetPageHandler)

		req := httptest.NewRequest("GET", "/api/pages/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var page models.Page
		json.Unmarshal(w.Body.Bytes(), &page)
		if page.ID != 1 {
			t.Errorf("Expected ID 1, got %d", page.ID)
		}
	})

	t.Run("UpdatePage", func(t *testing.T) {
		r := chi.NewRouter()
		r.Put("/api/pages/{id}", handlers.UpdatePageHandler)

		updatedPage := testPage
		updatedPage.Title = "Updated Test Page"
		data, _ := json.Marshal(updatedPage)

		req := httptest.NewRequest("PUT", "/api/pages/1", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("DeletePage", func(t *testing.T) {
		r := chi.NewRouter()
		r.Delete("/api/pages/{id}", handlers.DeletePageHandler)

		req := httptest.NewRequest("DELETE", "/api/pages/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status 204, got %d", w.Code)
		}
	})

	t.Run("ValidationErrors", func(t *testing.T) {
		// Test missing title
		invalidPage := models.Page{
			Slug:    "test-slug",
			Content: "Test content",
		}
		data, _ := json.Marshal(invalidPage)
		req := httptest.NewRequest("POST", "/api/pages", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handlers.CreatePageHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for missing title, got %d", w.Code)
		}

		// Test missing slug
		invalidPage2 := models.Page{
			Title:   "Test Title",
			Content: "Test content",
		}
		data, _ = json.Marshal(invalidPage2)
		req = httptest.NewRequest("POST", "/api/pages", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		handlers.CreatePageHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for missing slug, got %d", w.Code)
		}

		// Test invalid slug format
		invalidPage3 := models.Page{
			Title:   "Test Title",
			Slug:    "Invalid Slug!",
			Content: "Test content",
		}
		data, _ = json.Marshal(invalidPage3)
		req = httptest.NewRequest("POST", "/api/pages", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		handlers.CreatePageHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for invalid slug format, got %d", w.Code)
		}
	})

	t.Run("DuplicateSlug", func(t *testing.T) {
		// Create first page
		page1 := models.Page{
			Title:   "First Page",
			Slug:    "unique-slug",
			Content: "First content",
		}
		data, _ := json.Marshal(page1)
		req := httptest.NewRequest("POST", "/api/pages", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handlers.CreatePageHandler(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201 for first page, got %d", w.Code)
		}

		// Try to create second page with same slug
		page2 := models.Page{
			Title:   "Second Page",
			Slug:    "unique-slug",
			Content: "Second content",
		}
		data, _ = json.Marshal(page2)
		req = httptest.NewRequest("POST", "/api/pages", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		handlers.CreatePageHandler(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("Expected status 409 for duplicate slug, got %d", w.Code)
		}
	})
}