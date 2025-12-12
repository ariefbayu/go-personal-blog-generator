package generator

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ariefbayu/personal-blog-generator/internal/models"
	"github.com/ariefbayu/personal-blog-generator/internal/repository"
	"github.com/gomarkdown/markdown"
)

// Post represents a blog post for template rendering
type Post struct {
	Title              string
	Slug               string
	Content            template.HTML
	Tags               []string
	CreatedAt          time.Time
	CreatedAtFormatted string
}

// IndexData represents data for the index page template
type IndexData struct {
	Posts []IndexPost
}

// IndexPost represents a simplified post for the index page
type IndexPost struct {
	Title              string
	Slug               string
	CreatedAt          time.Time
	CreatedAtFormatted string
}

// PortfolioItem represents a portfolio item for template rendering
type PortfolioItem struct {
	Title            string
	ShortDescription template.HTML
	ProjectURL       string
	GithubURL        string
	ShowcaseImage    string
	SortOrder        int
}

// PortfolioData represents data for the portfolio page template
type PortfolioData struct {
	PortfolioItems []PortfolioItem
}

// Page represents a static page for template rendering
type Page struct {
	Title   string
	Slug    string
	Content template.HTML
}

// PageData represents data for the page template
type PageData struct {
	Title   string
	Slug    string
	Content template.HTML
}

// GenerateStaticSite generates static HTML files for all published posts, portfolio, and pages
func GenerateStaticSite(postRepo *repository.PostRepository, portfolioRepo *repository.PortfolioRepository, pageRepo *repository.PageRepository) error {
	// Query all published posts
	posts, err := postRepo.GetPublishedPosts()
	if err != nil {
		return fmt.Errorf("failed to query posts: %w", err)
	}

	// Parse the template
	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Ensure output directory exists
	outputDir := "html-outputs"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	postCount := 0
	for _, post := range posts {
		// Convert markdown to HTML
		contentHTML := mdToHTML(post.Content)

		// Parse tags
		var tags []string
		if post.Tags != "" {
			tags = strings.Split(post.Tags, ",")
			for i, tag := range tags {
				tags[i] = strings.TrimSpace(tag)
			}
		}

		// Create post struct for template
		templatePost := Post{
			Title:              post.Title,
			Slug:               post.Slug,
			Content:            contentHTML,
			Tags:               tags,
			CreatedAt:          post.CreatedAt,
			CreatedAtFormatted: post.CreatedAt.Format("January 2, 2006"),
		}

		// Create output file
		filename := filepath.Join(outputDir, post.Slug+".html")
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", filename, err)
		}

		// Execute template
		err = tmpl.Execute(file, templatePost)
		file.Close() // Close regardless of error
		if err != nil {
			return fmt.Errorf("failed to execute template for post %s: %w", post.Slug, err)
		}

		postCount++
	}

	// Generate index page
	err = generateIndexPage(posts, outputDir)
	if err != nil {
		return fmt.Errorf("failed to generate index page: %w", err)
	}

	// Generate portfolio page
	err = generatePortfolioPage(portfolioRepo, outputDir)
	if err != nil {
		return fmt.Errorf("failed to generate portfolio page: %w", err)
	}

	// Generate static pages
	err = generatePages(pageRepo, outputDir)
	if err != nil {
		return fmt.Errorf("failed to generate pages: %w", err)
	}

	return nil
}

// mdToHTML converts markdown content to HTML
func mdToHTML(content string) template.HTML {
	// Convert markdown to HTML
	htmlBytes := markdown.ToHTML([]byte(content), nil, nil)
	return template.HTML(htmlBytes)
}

// generateIndexPage creates the index.html file with recent posts
func generateIndexPage(posts []models.Post, outputDir string) error {
	// Parse the index template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		return fmt.Errorf("failed to parse index template: %w", err)
	}

	// Prepare index data (limit to 10 most recent posts)
	limit := 10
	if len(posts) < limit {
		limit = len(posts)
	}

	indexPosts := make([]IndexPost, limit)
	for i := 0; i < limit; i++ {
		post := posts[i]
		indexPosts[i] = IndexPost{
			Title:              post.Title,
			Slug:               post.Slug,
			CreatedAt:          post.CreatedAt,
			CreatedAtFormatted: post.CreatedAt.Format("January 2, 2006"),
		}
	}

	indexData := IndexData{
		Posts: indexPosts,
	}

	// Create index.html file
	indexPath := filepath.Join(outputDir, "index.html")
	file, err := os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("failed to create index.html: %w", err)
	}
	defer file.Close()

	// Execute template
	err = tmpl.Execute(file, indexData)
	if err != nil {
		return fmt.Errorf("failed to execute index template: %w", err)
	}

	return nil
}

// generatePortfolioPage creates the portfolio.html file with all portfolio items
func generatePortfolioPage(portfolioRepo *repository.PortfolioRepository, outputDir string) error {
	// Parse the portfolio template
	tmpl, err := template.ParseFiles("templates/portfolio.html")
	if err != nil {
		return fmt.Errorf("failed to parse portfolio template: %w", err)
	}

	// Query all portfolio items ordered by sort_order
	portfolioItems, err := portfolioRepo.GetAllPortfolioItems()
	if err != nil {
		return fmt.Errorf("failed to query portfolio items: %w", err)
	}

	// Prepare portfolio data
	templateItems := make([]PortfolioItem, len(portfolioItems))
	for i, item := range portfolioItems {
		// Convert markdown in short description to HTML if needed
		shortDescHTML := mdToHTML(item.ShortDescription)

		templateItems[i] = PortfolioItem{
			Title:            item.Title,
			ShortDescription: shortDescHTML,
			ProjectURL:       item.ProjectURL,
			GithubURL:        item.GithubURL,
			ShowcaseImage:    item.ShowcaseImage,
			SortOrder:        item.SortOrder,
		}
	}

	portfolioData := PortfolioData{
		PortfolioItems: templateItems,
	}

	// Create portfolio.html file
	portfolioPath := filepath.Join(outputDir, "portfolio.html")
	file, err := os.Create(portfolioPath)
	if err != nil {
		return fmt.Errorf("failed to create portfolio.html: %w", err)
	}
	defer file.Close()

	// Execute template
	err = tmpl.Execute(file, portfolioData)
	if err != nil {
		return fmt.Errorf("failed to execute portfolio template: %w", err)
	}

	return nil
}

// generatePages creates HTML files for all static pages
func generatePages(pageRepo *repository.PageRepository, outputDir string) error {
	// Parse the page template
	tmpl, err := template.ParseFiles("templates/page.html")
	if err != nil {
		return fmt.Errorf("failed to parse page template: %w", err)
	}

	// Query all pages
	pages, err := pageRepo.GetAllPages()
	if err != nil {
		return fmt.Errorf("failed to query pages: %w", err)
	}

	// Generate HTML for each page
	for _, page := range pages {
		// Convert markdown content to HTML
		contentHTML := mdToHTML(page.Content)

		// Create page data for template
		pageData := PageData{
			Title:   page.Title,
			Slug:    page.Slug,
			Content: contentHTML,
		}

		// Create output file
		filename := filepath.Join(outputDir, page.Slug+".html")
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", filename, err)
		}

		// Execute template
		err = tmpl.Execute(file, pageData)
		file.Close() // Close regardless of error
		if err != nil {
			return fmt.Errorf("failed to execute page template for %s: %w", page.Slug, err)
		}
	}

	return nil
}