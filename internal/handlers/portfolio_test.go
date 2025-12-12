package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ariefbayu/personal-blog-generator/internal/db"
	"github.com/ariefbayu/personal-blog-generator/internal/models"
	"github.com/ariefbayu/personal-blog-generator/internal/repository"
	"github.com/go-chi/chi/v5"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Change to project root for migration paths
	oldWd, _ := os.Getwd()
	os.Chdir("../../../personal-blog-generator")
	defer os.Chdir(oldWd)

	// Setup test database
	testDB, err := db.Connect(":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}

	err = db.Migrate(testDB)
	if err != nil {
		t.Fatalf("Failed to migrate test DB: %v", err)
	}

	return testDB
}

func TestPortfolioHandlers(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer db.Close()

	// Create repository and handlers
	portfolioRepo := repository.NewPortfolioRepository(db)
	handlers := NewPortfolioHandlers(portfolioRepo)

	// Test data
	testItem := models.PortfolioItem{
		Title:            "Test Portfolio Item",
		ShortDescription: "A test portfolio item",
		ProjectURL:       "https://example.com",
		GithubURL:        "https://github.com/test/repo",
		ShowcaseImage:    "/images/test.jpg",
		SortOrder:        1,
	}

	t.Run("CreatePortfolioItem", func(t *testing.T) {
		data, _ := json.Marshal(testItem)
		req := httptest.NewRequest("POST", "/api/portfolio", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handlers.CreatePortfolioItemHandler(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", w.Code)
		}

		var response map[string]int64
		json.Unmarshal(w.Body.Bytes(), &response)
		if response["id"] == 0 {
			t.Error("Expected non-zero ID in response")
		}
		testItem.ID = response["id"]
	})

	t.Run("GetPortfolioItems", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/portfolio", nil)
		w := httptest.NewRecorder()

		handlers.GetPortfolioItemsHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var items []models.PortfolioItem
		json.Unmarshal(w.Body.Bytes(), &items)
		if len(items) != 1 {
			t.Errorf("Expected 1 item, got %d", len(items))
		}
		if items[0].Title != testItem.Title {
			t.Errorf("Expected title %s, got %s", testItem.Title, items[0].Title)
		}
	})

	t.Run("GetPortfolioItem", func(t *testing.T) {
		r := chi.NewRouter()
		r.Get("/api/portfolio/{id}", handlers.GetPortfolioItemHandler)

		req := httptest.NewRequest("GET", "/api/portfolio/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var item models.PortfolioItem
		json.Unmarshal(w.Body.Bytes(), &item)
		if item.ID != 1 {
			t.Errorf("Expected ID 1, got %d", item.ID)
		}
	})

	t.Run("UpdatePortfolioItem", func(t *testing.T) {
		r := chi.NewRouter()
		r.Put("/api/portfolio/{id}", handlers.UpdatePortfolioItemHandler)

		updatedItem := testItem
		updatedItem.Title = "Updated Test Portfolio Item"
		data, _ := json.Marshal(updatedItem)

		req := httptest.NewRequest("PUT", "/api/portfolio/1", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("DeletePortfolioItem", func(t *testing.T) {
		r := chi.NewRouter()
		r.Delete("/api/portfolio/{id}", handlers.DeletePortfolioItemHandler)

		req := httptest.NewRequest("DELETE", "/api/portfolio/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status 204, got %d", w.Code)
		}
	})

	t.Run("ValidationErrors", func(t *testing.T) {
		// Test missing title
		invalidItem := models.PortfolioItem{
			ShortDescription: "Description without title",
		}
		data, _ := json.Marshal(invalidItem)
		req := httptest.NewRequest("POST", "/api/portfolio", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handlers.CreatePortfolioItemHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for missing title, got %d", w.Code)
		}

		// Test invalid URL
		invalidURLItem := models.PortfolioItem{
			Title:            "Test Item",
			ShortDescription: "Test description",
			ProjectURL:       "invalid-url",
		}
		data, _ = json.Marshal(invalidURLItem)
		req = httptest.NewRequest("POST", "/api/portfolio", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		handlers.CreatePortfolioItemHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for invalid URL, got %d", w.Code)
		}
	})
}