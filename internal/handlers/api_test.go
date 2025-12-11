package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ariefbayu/personal-blog-generator/internal/db"
	"github.com/ariefbayu/personal-blog-generator/internal/models"
	"github.com/ariefbayu/personal-blog-generator/internal/repository"
)

func TestGetPostsHandler(t *testing.T) {
	// Change to project root for migration paths
	oldWd, _ := os.Getwd()
	os.Chdir("../../../personal-blog-generator")
	defer os.Chdir(oldWd)

	// Setup test database
	testDB, err := db.Connect(":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}
	defer testDB.Close()

	err = db.Migrate(testDB)
	if err != nil {
		t.Fatalf("Failed to migrate test DB: %v", err)
	}

	// Insert test data
	_, err = testDB.Exec("INSERT INTO posts (title, slug, published) VALUES (?, ?, ?)", "Test Post", "test-post", true)
	if err != nil {
		t.Fatalf("Failed to insert test post: %v", err)
	}

	// Setup handler
	postRepo := repository.NewPostRepository(testDB)
	apiHandlers := NewAPIHandlers(postRepo)

	// Test request
	req := httptest.NewRequest("GET", "/api/posts", nil)
	w := httptest.NewRecorder()

	apiHandlers.GetPostsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var posts []models.Post
	err = json.NewDecoder(w.Body).Decode(&posts)
	if err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	if len(posts) != 1 {
		t.Errorf("Expected 1 post, got %d", len(posts))
	}

	if posts[0].Title != "Test Post" {
		t.Errorf("Expected title 'Test Post', got '%s'", posts[0].Title)
	}
}
