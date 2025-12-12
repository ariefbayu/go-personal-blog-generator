package generator

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ariefbayu/personal-blog-generator/internal/models"
	"github.com/ariefbayu/personal-blog-generator/internal/repository"
	_ "modernc.org/sqlite"
)

func TestGenerateStaticSite(t *testing.T) {
	// Create a temporary database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			slug TEXT UNIQUE NOT NULL,
			content TEXT NOT NULL,
			tags TEXT,
			published BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	// Create repository
	repo := repository.NewPostRepository(db)

	// Insert test data
	testPost := &models.Post{
		Title:    "Test Post",
		Slug:     "test-post",
		Content:  "# Hello World\n\nThis is a **test** post.",
		Tags:     "test,blog",
		Published: true,
		CreatedAt: time.Now(),
	}
	err = repo.CreatePost(testPost)
	if err != nil {
		t.Fatal(err)
	}

	// Create template directory and file
	templateDir := "templates"
	err = os.MkdirAll(templateDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(templateDir)

	templateContent := `<!DOCTYPE html>
<html>
<head><title>{{.Title}}</title></head>
<body>
<h1>{{.Title}}</h1>
<div>{{.Content}}</div>
<p>Published: {{.CreatedAtFormatted}}</p>
{{if .Tags}}
Tags: {{range .Tags}}<span>{{.}}</span> {{end}}
{{end}}
</body>
</html>`
	err = os.WriteFile(filepath.Join(templateDir, "post.html"), []byte(templateContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	indexTemplateContent := `<!DOCTYPE html>
<html>
<head><title>My Blog</title></head>
<body>
<h1>My Blog</h1>
{{if .Posts}}
<ul>
{{range .Posts}}
<li><a href="/{{.Slug}}.html">{{.Title}}</a> - {{.CreatedAtFormatted}}</li>
{{end}}
</ul>
{{else}}
<p>No posts</p>
{{end}}
</body>
</html>`
	err = os.WriteFile(filepath.Join(templateDir, "index.html"), []byte(indexTemplateContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Generate static site
	err = GenerateStaticSite(repo)
	if err != nil {
		t.Fatalf("GenerateStaticSite failed: %v", err)
	}

	// Check if output file exists
	outputFile := "html-outputs/test-post.html"
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file %s was not created", outputFile)
	}

	// Check file content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)
	if !contains(contentStr, "Test Post") {
		t.Error("Generated HTML does not contain post title")
	}
	if !contains(contentStr, "<h1>Hello World</h1>") {
		t.Error("Generated HTML does not contain converted markdown")
	}
	if !contains(contentStr, "<strong>test</strong>") {
		t.Error("Generated HTML does not contain bold formatting")
	}

	// Check if index.html exists
	indexFile := "html-outputs/index.html"
	if _, err := os.Stat(indexFile); os.IsNotExist(err) {
		t.Errorf("Index file %s was not created", indexFile)
	}

	// Check index.html content
	indexContent, err := os.ReadFile(indexFile)
	if err != nil {
		t.Fatal(err)
	}

	indexContentStr := string(indexContent)
	if !contains(indexContentStr, "Test Post") {
		t.Error("Index HTML does not contain post title")
	}
	if !contains(indexContentStr, "/test-post.html") {
		t.Error("Index HTML does not contain link to post")
	}

	// Clean up
	os.RemoveAll("html-outputs")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsAt(s, substr)))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}