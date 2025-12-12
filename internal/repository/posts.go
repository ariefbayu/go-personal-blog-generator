package repository

import (
"database/sql"

"github.com/ariefbayu/personal-blog-generator/internal/models"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) GetAllPosts() ([]models.Post, error) {
	rows, err := r.db.Query("SELECT id, title, slug, published, created_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Slug, &post.Published, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepository) CreatePost(post *models.Post) error {
	err := r.db.QueryRow("INSERT INTO posts (title, slug, content, tags, published) VALUES (?, ?, ?, ?, ?) RETURNING id", post.Title, post.Slug, post.Content, post.Tags, post.Published).Scan(&post.ID)
	return err
}

func (r *PostRepository) GetPostByID(id int64) (*models.Post, error) {
	var post models.Post
	err := r.db.QueryRow("SELECT id, title, slug, content, tags, published, created_at FROM posts WHERE id = ?", id).Scan(&post.ID, &post.Title, &post.Slug, &post.Content, &post.Tags, &post.Published, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) UpdatePost(post *models.Post) error {
	_, err := r.db.Exec("UPDATE posts SET title = ?, slug = ?, content = ?, tags = ?, published = ? WHERE id = ?", post.Title, post.Slug, post.Content, post.Tags, post.Published, post.ID)
	return err
}

func (r *PostRepository) DeletePost(id int64) error {
	result, err := r.db.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PostRepository) GetPublishedPosts() ([]models.Post, error) {
	rows, err := r.db.Query("SELECT id, title, slug, content, tags, published, created_at FROM posts WHERE published = true ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Slug, &post.Content, &post.Tags, &post.Published, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, rows.Err()
}
