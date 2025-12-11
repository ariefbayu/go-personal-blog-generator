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
	apiHandlers := NewAPIHandlers(postRepo)

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

func TestGetPostHandler(t *testing.T) {
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
	_, err = testDB.Exec("INSERT INTO posts (title, slug, content, tags, published) VALUES (?, ?, ?, ?, ?)", "Test Post", "test-post", "Test content", "test", true)
	if err != nil {
		t.Fatalf("Failed to insert test post: %v", err)
	}

	// Setup handler
	postRepo := repository.NewPostRepository(testDB)
	apiHandlers := NewAPIHandlers(postRepo)

	// Test request
	req := httptest.NewRequest("GET", "/api/posts/1", nil)
	w := httptest.NewRecorder()

	apiHandlers.GetPostHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var post models.Post
	err = json.NewDecoder(w.Body).Decode(&post)
	if err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	if post.Title != "Test Post" {
		t.Errorf("Expected title 'Test Post', got '%s'", post.Title)
	}
}

func TestUpdatePostHandler(t *testing.T) {
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
	_, err = testDB.Exec("INSERT INTO posts (title, slug, content, tags, published) VALUES (?, ?, ?, ?, ?)", "Test Post", "test-post", "Test content", "test", true)
	if err != nil {
		t.Fatalf("Failed to insert test post: %v", err)
	}

	// Setup handler
	postRepo := repository.NewPostRepository(testDB)
	apiHandlers := NewAPIHandlers(postRepo)

	// Test data
	updateData := models.Post{
		Title:     "Updated Post",
		Slug:      "updated-post",
		Content:   "Updated content",
		Tags:      "updated",
		Published: false,
	}

	// Test request
	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", "/api/posts/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	apiHandlers.UpdatePostHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify post was updated
	post, err := postRepo.GetPostByID(1)
	if err != nil {
		t.Fatalf("Failed to get post: %v", err)
	}

	if post.Title != "Updated Post" {
		t.Errorf("Expected title 'Updated Post', got '%s'", post.Title)
	}
}

func TestDeletePostHandler(t *testing.T) {
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
	_, err = testDB.Exec("INSERT INTO posts (title, slug, content, tags, published) VALUES (?, ?, ?, ?, ?)", "Test Post", "test-post", "Test content", "test", true)
	if err != nil {
		t.Fatalf("Failed to insert test post: %v", err)
	}

	// Setup handler
	postRepo := repository.NewPostRepository(testDB)
	apiHandlers := NewAPIHandlers(postRepo)

	// Test request
	req := httptest.NewRequest("DELETE", "/api/posts/1", nil)
	w := httptest.NewRecorder()

	apiHandlers.DeletePostHandler(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	// Verify post was deleted
	_, err = postRepo.GetPostByID(1)
	if err == nil {
		t.Errorf("Expected post to be deleted, but it still exists")
	}
}
