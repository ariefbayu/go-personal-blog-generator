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

	_, err = db.Exec(`
		CREATE TABLE pages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			slug TEXT UNIQUE NOT NULL,
			content TEXT NOT NULL,
			show_in_nav BOOLEAN DEFAULT TRUE,
			sort_order INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE settings (
			id INTEGER PRIMARY KEY DEFAULT 1,
			site_name TEXT DEFAULT 'My Personal Blog',
			show_portfolio_menu BOOLEAN DEFAULT TRUE,
			show_posts_menu BOOLEAN DEFAULT TRUE,
			menu_order TEXT DEFAULT '["posts", "portfolio", "pages"]',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	// Insert default settings
	_, err = db.Exec(`INSERT INTO settings (id, site_name, show_portfolio_menu, show_posts_menu, menu_order) VALUES (1, 'My Personal Blog', 1, 1, '["posts", "portfolio", "pages"]')`)
	if err != nil {
		t.Fatal(err)
	}

	// Create repositories
	postRepo := repository.NewPostRepository(db)
	portfolioRepo := repository.NewPortfolioRepository(db)
	pageRepo := repository.NewPageRepository(db)
	settingsRepo := repository.NewSettingsRepository(db)

	// Insert test data
	testPost1 := &models.Post{
		Title:     "Test Post 1",
		Slug:      "test-post-1",
		Content:   "This is a test post content",
		Tags:      "test,blog",
		Published: true,
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	err = postRepo.CreatePost(testPost1)
	if err != nil {
		t.Fatal(err)
	}

	testPost2 := &models.Post{
		Title:     "Test Post 2",
		Slug:      "test-post-2",
		Content:   "This is a test post content for the second post",
		Tags:      "test,blog",
		Published: true,
		CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
	}
	err = postRepo.CreatePost(testPost2)
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

	// Insert test page data
	testPage := &models.Page{
		Title:     "About Us",
		Slug:      "about",
		Content:   "# About Us\n\nThis is a **test** page with markdown content.",
		ShowInNav: true,
		SortOrder: 1,
	}
	err = pageRepo.CreatePage(testPage)
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

	templateContent := `<h1>{{.Title}}</h1>
<div>{{.Content}}</div>
<p>Published: {{.CreatedAtFormatted}}</p>
{{if .Tags}}
Tags: {{range .Tags}}<span>{{.}}</span> {{end}}
{{end}}`
	err = os.WriteFile(filepath.Join(templateDir, "post.html"), []byte(templateContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	indexTemplateContent := `<!DOCTYPE html>
<html>
<head><title>My Blog</title></head>
<body>
<nav><ul>{{range .NavLinks}}<li><a href="{{.URL}}">{{.Title}}</a></li>{{end}}</ul></nav>
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

	// Create header template
	headerTemplateContent := `<!DOCTYPE html>
<html>
<head><title>My Blog</title></head>
<body>
<nav><ul>{{range .NavLinks}}<li><a href="{{.URL}}">{{.Title}}</a></li>{{end}}</ul></nav>`
	err = os.WriteFile(filepath.Join(templateDir, "header.html"), []byte(headerTemplateContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create footer template
	footerTemplateContent := `
</body>
</html>`
	err = os.WriteFile(filepath.Join(templateDir, "footer.html"), []byte(footerTemplateContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	portfolioTemplateContent := `<!DOCTYPE html>
<html>
<head><title>My Portfolio</title></head>
<body>
<nav><ul>{{range .NavLinks}}<li><a href="{{.URL}}">{{.Title}}</a></li>{{end}}</ul></nav>
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

	pageTemplateContent := `<!DOCTYPE html>
<html>
<head><title>{{.Title}}</title></head>
<body>
<nav><ul>{{range .NavLinks}}<li><a href="{{.URL}}">{{.Title}}</a></li>{{end}}</ul></nav>
<h1>{{.Title}}</h1>
<div>{{.Content}}</div>
</body>
</html>`
	err = os.WriteFile(filepath.Join(templateDir, "page.html"), []byte(pageTemplateContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	postsTemplateContent := `<h1>All Posts</h1>
{{if .Posts}}
<div class="posts-list">
{{range .Posts}}
<article class="post-item">
<h2><a href="/posts/{{.Slug}}.html">{{.Title}}</a></h2>
<time>{{.CreatedAtFormatted}}</time>
{{if .Tags}}{{range .Tags}}<span class="tag">{{.}}</span>{{end}}{{end}}
<p>{{.Excerpt}}</p>
</article>
{{end}}
</div>
{{else}}
<p>No posts</p>
{{end}}`
	err = os.WriteFile(filepath.Join(templateDir, "posts.html"), []byte(postsTemplateContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Generate static site
	err = GenerateStaticSite(postRepo, portfolioRepo, pageRepo, settingsRepo, "./templates", "./html-outputs")
	if err != nil {
		t.Fatalf("GenerateStaticSite failed: %v", err)
	}

	// Check if output file exists
	outputFile := "html-outputs/test-post-1.html"
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file %s was not created", outputFile)
	}

	// Check file content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)
	if !contains(contentStr, "Test Post 1") {
		t.Error("Generated HTML does not contain post title")
	}
	if !contains(contentStr, "This is a test post content") {
		t.Error("Generated HTML does not contain post content")
	}
	// if !contains(contentStr, "Home") {
	// 	t.Error("Generated HTML does not contain navigation Home link")
	// }
	// if !contains(contentStr, "Portfolio") {
	// 	t.Error("Generated HTML does not contain navigation Portfolio link")
	// }
	// if !contains(contentStr, "About Us") {
	// 	t.Error("Generated HTML does not contain navigation About Us link")
	// }

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
	if !contains(indexContentStr, "Test Post 2") {
		t.Error("Index HTML does not contain post title")
	}
	if !contains(indexContentStr, "/test-post-2.html") {
		t.Error("Index HTML does not contain link to post")
	}
	if !contains(indexContentStr, "Home") {
		t.Error("Index HTML does not contain navigation Home link")
	}
	// if !contains(indexContentStr, "Portfolio") {
	// 	t.Error("Index HTML does not contain navigation Portfolio link")
	// }
	// if !contains(indexContentStr, "About Us") {
	// 	t.Error("Index HTML does not contain navigation About Us link")
	// }

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

	// Check if about.html exists
	pageFile := "html-outputs/about.html"
	if _, err := os.Stat(pageFile); os.IsNotExist(err) {
		t.Errorf("Page file %s was not created", pageFile)
	}

	// Check about.html content
	pageContent, err := os.ReadFile(pageFile)
	if err != nil {
		t.Fatal(err)
	}

	pageContentStr := string(pageContent)
	if !contains(pageContentStr, "About Us") {
		t.Error("Page HTML does not contain page title")
	}
	if !contains(pageContentStr, "<h1>About Us</h1>") {
		t.Error("Page HTML does not contain converted markdown heading")
	}
	if !contains(pageContentStr, "<strong>test</strong>") {
		t.Error("Page HTML does not contain converted markdown in content")
	}

	// Check if posts.html exists
	postsFile := "html-outputs/posts.html"
	if _, err := os.Stat(postsFile); os.IsNotExist(err) {
		t.Errorf("Posts file %s was not created", postsFile)
	}

	// Check posts.html content
	postsContent, err := os.ReadFile(postsFile)
	if err != nil {
		t.Fatal(err)
	}

	postsContentStr := string(postsContent)
	if !contains(postsContentStr, "All Posts") {
		t.Error("Posts HTML does not contain page title")
	}
	if !contains(postsContentStr, "Test Post 1") {
		t.Error("Posts HTML does not contain first post title")
	}
	if !contains(postsContentStr, "Test Post 2") {
		t.Error("Posts HTML does not contain second post title")
	}
	if !contains(postsContentStr, "/posts/test-post-1.html") {
		t.Error("Posts HTML does not contain first post link")
	}
	if !contains(postsContentStr, "/posts/test-post-2.html") {
		t.Error("Posts HTML does not contain second post link")
	}
	if !contains(postsContentStr, "This is a test post content") {
		t.Error("Posts HTML does not contain post excerpt")
	}
	if !contains(postsContentStr, "2023-01-01") {
		t.Error("Posts HTML does not contain post date")
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

func TestCopyStaticAssets(t *testing.T) {
	// Create temporary directories for testing
	tempDir := t.TempDir()
	templatePath := filepath.Join(tempDir, "templates")
	staticPath := filepath.Join(templatePath, "static")
	cssPath := filepath.Join(staticPath, "css")
	outputPath := filepath.Join(tempDir, "output")

	// Create directories
	if err := os.MkdirAll(templatePath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(cssPath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a test CSS file
	testCSSContent := `
.heading-1 {
    font-size: 2rem;
    font-weight: bold;
}
.text-body {
    font-size: 1rem;
}
`
	cssFile := filepath.Join(cssPath, "styles.css")
	if err := os.WriteFile(cssFile, []byte(testCSSContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Run copyStaticAssets
	err := copyStaticAssets(templatePath, outputPath)
	if err != nil {
		t.Fatalf("copyStaticAssets failed: %v", err)
	}

	// Verify CSS file was copied to output directory
	outputCSSFile := filepath.Join(outputPath, "css", "styles.css")
	copiedContent, err := os.ReadFile(outputCSSFile)
	if err != nil {
		t.Fatalf("Failed to read copied CSS file: %v", err)
	}

	if string(copiedContent) != testCSSContent {
		t.Errorf("Copied CSS content does not match original\nExpected: %s\nGot: %s", testCSSContent, string(copiedContent))
	}
}

func TestCopyStaticAssetsNoCSS(t *testing.T) {
	// Create temporary directories without CSS
	tempDir := t.TempDir()
	templatePath := filepath.Join(tempDir, "templates")
	outputPath := filepath.Join(tempDir, "output")

	// Create directories
	if err := os.MkdirAll(templatePath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		t.Fatal(err)
	}

	// Run copyStaticAssets - should succeed even if no CSS directory exists
	err := copyStaticAssets(templatePath, outputPath)
	if err != nil {
		t.Fatalf("copyStaticAssets should succeed when no CSS directory exists: %v", err)
	}

	// Verify no CSS directory was created in output
	outputCSSDir := filepath.Join(outputPath, "css")
	if _, err := os.Stat(outputCSSDir); !os.IsNotExist(err) {
		t.Error("CSS directory should not be created when source doesn't exist")
	}
}
