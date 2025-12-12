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

	_, err = db.Exec(`
		CREATE TABLE portfolio_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			short_description TEXT NOT NULL,
			project_url TEXT,
			github_url TEXT,
			showcase_image TEXT,
			sort_order INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	// Create repositories
	postRepo := repository.NewPostRepository(db)
	portfolioRepo := repository.NewPortfolioRepository(db)

	// Insert test data
	testPost := &models.Post{
		Title:    "Test Post",
		Slug:     "test-post",
		Content:  "# Hello World\n\nThis is a **test** post.",
		Tags:     "test,blog",
		Published: true,
		CreatedAt: time.Now(),
	}
	err = postRepo.CreatePost(testPost)
	if err != nil {
		t.Fatal(err)
	}

	// Insert test portfolio data
	testPortfolioItem := &models.PortfolioItem{
		Title:            "Test Portfolio Item",
		ShortDescription: "A **test** portfolio item with markdown.",
		ProjectURL:       "https://example.com/project",
		GithubURL:        "https://github.com/test/repo",
		ShowcaseImage:    "/images/test.jpg",
		SortOrder:        1,
	}
	err = portfolioRepo.CreatePortfolioItem(testPortfolioItem)
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

	portfolioTemplateContent := `<!DOCTYPE html>
<html>
<head><title>My Portfolio</title></head>
<body>
<h1>My Portfolio</h1>
{{if .PortfolioItems}}
<div class="portfolio-grid">
{{range .PortfolioItems}}
<div class="portfolio-item">
{{if .ShowcaseImage}}<img src="{{.ShowcaseImage}}" alt="{{.Title}}" />{{end}}
<h3>{{.Title}}</h3>
<div>{{.ShortDescription}}</div>
{{if .ProjectURL}}<a href="{{.ProjectURL}}">View Project</a>{{end}}
{{if .GithubURL}}<a href="{{.GithubURL}}">View on GitHub</a>{{end}}
</div>
{{end}}
</div>
{{else}}
<p>No portfolio items</p>
{{end}}
</body>
</html>`
	err = os.WriteFile(filepath.Join(templateDir, "portfolio.html"), []byte(portfolioTemplateContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Generate static site
	err = GenerateStaticSite(postRepo, portfolioRepo)
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

	// Check if portfolio.html exists
	portfolioFile := "html-outputs/portfolio.html"
	if _, err := os.Stat(portfolioFile); os.IsNotExist(err) {
		t.Errorf("Portfolio file %s was not created", portfolioFile)
	}

	// Check portfolio.html content
	portfolioContent, err := os.ReadFile(portfolioFile)
	if err != nil {
		t.Fatal(err)
	}

	portfolioContentStr := string(portfolioContent)
	if !contains(portfolioContentStr, "Test Portfolio Item") {
		t.Error("Portfolio HTML does not contain portfolio item title")
	}
	if !contains(portfolioContentStr, "https://example.com/project") {
		t.Error("Portfolio HTML does not contain project URL")
	}
	if !contains(portfolioContentStr, "/images/test.jpg") {
		t.Error("Portfolio HTML does not contain showcase image")
	}
	if !contains(portfolioContentStr, "<strong>test</strong>") {
		t.Error("Portfolio HTML does not contain converted markdown in description")
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