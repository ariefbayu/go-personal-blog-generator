package generator

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

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

// GenerateStaticSite generates static HTML files for all published posts
func GenerateStaticSite(repo *repository.PostRepository) error {
	// Query all published posts
	posts, err := repo.GetPublishedPosts()
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

	return nil
}

// mdToHTML converts markdown content to HTML
func mdToHTML(content string) template.HTML {
	// Convert markdown to HTML
	htmlBytes := markdown.ToHTML([]byte(content), nil, nil)
	return template.HTML(htmlBytes)
}