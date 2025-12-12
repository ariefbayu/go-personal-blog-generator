package handlers

import (
	"bytes"
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
	portfolioRepo := repository.NewPortfolioRepository(testDB)
	pageRepo := repository.NewPageRepository(testDB)
	apiHandlers := NewAPIHandlers(postRepo, portfolioRepo, pageRepo)

	// Test request
	req := httptest.NewRequest("GET", "/api/posts", nil)
	w := httptest.NewRecorder()

	apiHandlers.GetPostsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	err = json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	postsInterface, ok := response["posts"]
	if !ok {
		t.Fatal("Response missing 'posts' field")
	}

	postsData, err := json.Marshal(postsInterface)
	if err != nil {
		t.Fatalf("Failed to marshal posts: %v", err)
	}

	var posts []models.Post
	err = json.Unmarshal(postsData, &posts)
	if err != nil {
		t.Fatalf("Failed to unmarshal posts: %v", err)
	}

	if len(posts) != 1 {
		t.Errorf("Expected 1 post, got %d", len(posts))
	}

	if posts[0].Title != "Test Post" {
		t.Errorf("Expected title 'Test Post', got '%s'", posts[0].Title)
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
}

func TestCreatePostHandler(t *testing.T) {
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

	// Setup handler
	postRepo := repository.NewPostRepository(testDB)
	portfolioRepo := repository.NewPortfolioRepository(testDB)
	pageRepo := repository.NewPageRepository(testDB)
	apiHandlers := NewAPIHandlers(postRepo, portfolioRepo, pageRepo)

	// Test data
	postData := models.Post{
		Title:     "Test Post",
		Slug:      "test-post",
		Content:   "Test content",
		Tags:      "test,example",
		Published: true,
	}

	// Test request
	body, _ := json.Marshal(postData)
	req := httptest.NewRequest("POST", "/api/posts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	apiHandlers.CreatePostHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response map[string]int64
	err = json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["id"] == 0 {
		t.Errorf("Expected non-zero id, got %d", response["id"])
	}

	// Verify post was created
	posts, err := postRepo.GetAllPosts()
	if err != nil {
		t.Fatalf("Failed to get posts: %v", err)
	}

	if len(posts) != 1 {
		t.Errorf("Expected 1 post, got %d", len(posts))
	}

	if posts[0].Title != "Test Post" {
		t.Errorf("Expected title 'Test Post', got '%s'", posts[0].Title)
	}
}
