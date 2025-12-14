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
	NavigationData
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

// IndexData represents data for the index page template
type IndexData struct {
	Posts []IndexPost
	NavigationData
	PortfolioItems []PortfolioItem
}

// PostsData represents data for the posts listing page template
type PostsData struct {
	Posts []PostItem
	NavigationData
}

// PostItem represents a post item for the posts listing page
type PostItem struct {
	Title              string
	Slug               string
	CreatedAt          time.Time
	CreatedAtFormatted string
	Tags               []string
	Excerpt            string
}

// IndexPost represents a simplified post for the index page
type IndexPost struct {
	Title              string
	Slug               string
	CreatedAt          time.Time
	CreatedAtFormatted string
}

// PortfolioData represents data for the portfolio page template
type PortfolioData struct {
	PortfolioItems []PortfolioItem
	NavigationData
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
	NavigationData
}

// NavLink represents a navigation link
type NavLink struct {
	Title     string
	URL       string
	SortOrder int
}

// NavigationData represents navigation data for templates
type NavigationData struct {
	NavLinks []NavLink
}

// buildNavigationData builds navigation links from pages and standard links
func buildNavigationData(pageRepo *repository.PageRepository) ([]NavLink, error) {
	// Query pages that should show in navigation
	pages, err := pageRepo.GetPagesForNavigation()
	if err != nil {
		return nil, fmt.Errorf("failed to query navigation pages: %w", err)
	}

	var navLinks []NavLink

	// Add standard links
	navLinks = append(navLinks, NavLink{
		Title:     "Home",
		URL:       "/index.html",
		SortOrder: 0,
	})

	navLinks = append(navLinks, NavLink{
		Title:     "Blog",
		URL:       "/posts.html",
		SortOrder: 1,
	})

	// hide portfolio link for now
	// navLinks = append(navLinks, NavLink{
	// 	Title:     "Portfolio",
	// 	URL:       "/portfolio.html",
	// 	SortOrder: 2,
	// })

	// Add page links
	for _, page := range pages {
		navLinks = append(navLinks, NavLink{
			Title:     page.Title,
			URL:       "/" + page.Slug + ".html",
			SortOrder: page.SortOrder + 10, // Offset to put after standard links
		})
	}

	// Sort by sort order
	for i := 0; i < len(navLinks)-1; i++ {
		for j := i + 1; j < len(navLinks); j++ {
			if navLinks[i].SortOrder > navLinks[j].SortOrder {
				navLinks[i], navLinks[j] = navLinks[j], navLinks[i]
			}
		}
	}

	return navLinks, nil
}

// GenerateStaticSite generates static HTML files for all published posts, portfolio, and pages
func GenerateStaticSite(postRepo *repository.PostRepository, portfolioRepo *repository.PortfolioRepository, pageRepo *repository.PageRepository, templatePath, outputPath string) error {
	// Query all published posts
	posts, err := postRepo.GetPublishedPosts()
	if err != nil {
		return fmt.Errorf("failed to query posts: %w", err)
	}

	// Build navigation data
	navLinks, err := buildNavigationData(pageRepo)
	if err != nil {
		return fmt.Errorf("failed to build navigation data: %w", err)
	}
	navData := NavigationData{NavLinks: navLinks}

	// Parse the header, post content, and footer templates
	tmpl, err := template.ParseFiles(
		filepath.Join(templatePath, "header.html"),
		filepath.Join(templatePath, "post.html"),
		filepath.Join(templatePath, "footer.html"),
	)
	if err != nil {
		return fmt.Errorf("failed to parse post templates: %w", err)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputPath, 0755); err != nil {
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
			NavigationData:     navData,
		}

		// Create output file
		filename := filepath.Join(outputPath, post.Slug+".html")
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", filename, err)
		}

		// Execute templates in sequence: header, post content, footer
		if err := tmpl.ExecuteTemplate(file, "header.html", templatePost); err != nil {
			return fmt.Errorf("failed to execute header template for post %s: %w", post.Slug, err)
		}
		if err := tmpl.ExecuteTemplate(file, "post.html", templatePost); err != nil {
			return fmt.Errorf("failed to execute post template for post %s: %w", post.Slug, err)
		}
		if err := tmpl.ExecuteTemplate(file, "footer.html", templatePost); err != nil {
			return fmt.Errorf("failed to execute footer template for post %s: %w", post.Slug, err)
		}

		postCount++
	}

	// Generate index page
	err = generateIndexPage(posts, portfolioRepo, outputPath, templatePath, navData)
	if err != nil {
		return fmt.Errorf("failed to generate index page: %w", err)
	}

	// Generate posts listing page
	err = generatePostsPage(posts, outputPath, templatePath, navData)
	if err != nil {
		return fmt.Errorf("failed to generate posts page: %w", err)
	}

	// Generate portfolio page
	err = generatePortfolioPage(portfolioRepo, outputPath, templatePath, navData)
	if err != nil {
		return fmt.Errorf("failed to generate portfolio page: %w", err)
	}

	// Generate static pages
	err = generatePages(pageRepo, outputPath, templatePath, navData)
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
func generateIndexPage(posts []models.Post, portfolioRepo *repository.PortfolioRepository, outputPath, templatePath string, navData NavigationData) error {
	// Parse the header, index content, and footer templates
	tmpl, err := template.ParseFiles(
		filepath.Join(templatePath, "header.html"),
		filepath.Join(templatePath, "index.html"),
		filepath.Join(templatePath, "footer.html"),
	)
	if err != nil {
		return fmt.Errorf("failed to parse index templates: %w", err)
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

	indexData := IndexData{
		Posts:          indexPosts,
		NavigationData: navData,
		PortfolioItems: templateItems,
	}

	// Create index.html file
	indexPath := filepath.Join(outputPath, "index.html")
	file, err := os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("failed to create index.html: %w", err)
	}
	defer file.Close()

	// Execute templates in sequence: header, index content, footer
	if err := tmpl.ExecuteTemplate(file, "header.html", indexData); err != nil {
		return fmt.Errorf("failed to execute header template: %w", err)
	}
	if err := tmpl.ExecuteTemplate(file, "index.html", indexData); err != nil {
		return fmt.Errorf("failed to execute index template: %w", err)
	}
	if err := tmpl.ExecuteTemplate(file, "footer.html", indexData); err != nil {
		return fmt.Errorf("failed to execute footer template: %w", err)
	}

	return nil
}

// generatePostsPage creates the posts.html file with all posts
func generatePostsPage(posts []models.Post, outputPath, templatePath string, navData NavigationData) error {
	// Sort posts by created date descending (newest first)
	for i := 0; i < len(posts)-1; i++ {
		for j := i + 1; j < len(posts); j++ {
			if posts[i].CreatedAt.Before(posts[j].CreatedAt) {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}
	}

	// Parse the header, posts content, and footer templates
	tmpl, err := template.ParseFiles(
		filepath.Join(templatePath, "header.html"),
		filepath.Join(templatePath, "posts.html"),
		filepath.Join(templatePath, "footer.html"),
	)

	if err != nil {
		return fmt.Errorf("failed to parse posts templates: %w", err)
	}

	// Prepare posts data
	postItems := make([]PostItem, len(posts))
	for i, post := range posts {
		// Parse tags
		var tags []string
		if post.Tags != "" {
			tags = strings.Split(post.Tags, ",")
			for j, tag := range tags {
				tags[j] = strings.TrimSpace(tag)
			}
		}

		// Create excerpt from content (first 200 characters)
		content := strings.ReplaceAll(post.Content, "\n", " ")
		excerpt := content
		if len(content) > 200 {
			excerpt = content[:200] + "..."
		}

		postItems[i] = PostItem{
			Title:              post.Title,
			Slug:               post.Slug,
			CreatedAt:          post.CreatedAt,
			CreatedAtFormatted: post.CreatedAt.Format("2006-01-02"),
			Tags:               tags,
			Excerpt:            excerpt,
		}
	}

	postsData := PostsData{
		Posts:          postItems,
		NavigationData: navData,
	}

	// Create posts.html file
	postsPath := filepath.Join(outputPath, "posts.html")
	file, err := os.Create(postsPath)
	if err != nil {
		return fmt.Errorf("failed to create posts.html: %w", err)
	}
	defer file.Close()

	// Execute templates in sequence: header, posts content, footer
	if err := tmpl.ExecuteTemplate(file, "header.html", postsData); err != nil {
		return fmt.Errorf("failed to execute header template: %w", err)
	}
	if err := tmpl.ExecuteTemplate(file, "posts.html", postsData); err != nil {
		return fmt.Errorf("failed to execute posts template: %w", err)
	}
	if err := tmpl.ExecuteTemplate(file, "footer.html", postsData); err != nil {
		return fmt.Errorf("failed to execute footer template: %w", err)
	}

	return nil
}

// generatePortfolioPage creates the portfolio.html file with all portfolio items
func generatePortfolioPage(portfolioRepo *repository.PortfolioRepository, outputPath, templatePath string, navData NavigationData) error {
	// Parse the portfolio template
	tmpl, err := template.ParseFiles(filepath.Join(templatePath, "portfolio.html"))
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
		NavigationData: navData,
	}

	// Create portfolio.html file
	portfolioPath := filepath.Join(outputPath, "portfolio.html")
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
func generatePages(pageRepo *repository.PageRepository, outputPath, templatePath string, navData NavigationData) error {
	// Parse the page template
	tmpl, err := template.ParseFiles(filepath.Join(templatePath, "page.html"))
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
			Title:          page.Title,
			Slug:           page.Slug,
			Content:        contentHTML,
			NavigationData: navData,
		}

		// Create output file
		filename := filepath.Join(outputPath, page.Slug+".html")
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
