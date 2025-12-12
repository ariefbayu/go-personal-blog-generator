package repository

import (
	"database/sql"

	"github.com/ariefbayu/personal-blog-generator/internal/models"
)

type PageRepository struct {
	db *sql.DB
}

func NewPageRepository(db *sql.DB) *PageRepository {
	return &PageRepository{db: db}
}

func (r *PageRepository) GetAllPages() ([]models.Page, error) {
	rows, err := r.db.Query("SELECT id, title, slug, content, show_in_nav, sort_order, created_at, updated_at FROM pages ORDER BY sort_order ASC, created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []models.Page
	for rows.Next() {
		var page models.Page
		err := rows.Scan(&page.ID, &page.Title, &page.Slug, &page.Content, &page.ShowInNav, &page.SortOrder, &page.CreatedAt, &page.UpdatedAt)
		if err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}
	return pages, nil
}

func (r *PageRepository) GetPageByID(id int64) (*models.Page, error) {
	var page models.Page
	err := r.db.QueryRow("SELECT id, title, slug, content, show_in_nav, sort_order, created_at, updated_at FROM pages WHERE id = ?", id).Scan(&page.ID, &page.Title, &page.Slug, &page.Content, &page.ShowInNav, &page.SortOrder, &page.CreatedAt, &page.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *PageRepository) GetPageBySlug(slug string) (*models.Page, error) {
	var page models.Page
	err := r.db.QueryRow("SELECT id, title, slug, content, show_in_nav, sort_order, created_at, updated_at FROM pages WHERE slug = ?", slug).Scan(&page.ID, &page.Title, &page.Slug, &page.Content, &page.ShowInNav, &page.SortOrder, &page.CreatedAt, &page.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *PageRepository) GetPagesForNavigation() ([]models.Page, error) {
	rows, err := r.db.Query("SELECT id, title, slug, content, show_in_nav, sort_order, created_at, updated_at FROM pages WHERE show_in_nav = true ORDER BY sort_order ASC, created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []models.Page
	for rows.Next() {
		var page models.Page
		err := rows.Scan(&page.ID, &page.Title, &page.Slug, &page.Content, &page.ShowInNav, &page.SortOrder, &page.CreatedAt, &page.UpdatedAt)
		if err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}
	return pages, nil
}

func (r *PageRepository) CreatePage(page *models.Page) error {
	err := r.db.QueryRow("INSERT INTO pages (title, slug, content, show_in_nav, sort_order) VALUES (?, ?, ?, ?, ?) RETURNING id", page.Title, page.Slug, page.Content, page.ShowInNav, page.SortOrder).Scan(&page.ID)
	return err
}

func (r *PageRepository) UpdatePage(page *models.Page) error {
	_, err := r.db.Exec("UPDATE pages SET title = ?, slug = ?, content = ?, show_in_nav = ?, sort_order = ? WHERE id = ?", page.Title, page.Slug, page.Content, page.ShowInNav, page.SortOrder, page.ID)
	return err
}

func (r *PageRepository) DeletePage(id int64) error {
	result, err := r.db.Exec("DELETE FROM pages WHERE id = ?", id)
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