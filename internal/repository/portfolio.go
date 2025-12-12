package repository

import (
	"database/sql"

	"github.com/ariefbayu/personal-blog-generator/internal/models"
)

type PortfolioRepository struct {
	db *sql.DB
}

func NewPortfolioRepository(db *sql.DB) *PortfolioRepository {
	return &PortfolioRepository{db: db}
}

func (r *PortfolioRepository) GetAllPortfolioItems() ([]models.PortfolioItem, error) {
	rows, err := r.db.Query("SELECT id, title, short_description, project_url, github_url, showcase_image, sort_order, created_at, updated_at FROM portfolio_items ORDER BY sort_order ASC, created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.PortfolioItem
	for rows.Next() {
		var item models.PortfolioItem
		err := rows.Scan(&item.ID, &item.Title, &item.ShortDescription, &item.ProjectURL, &item.GithubURL, &item.ShowcaseImage, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *PortfolioRepository) GetPortfolioItemByID(id int64) (*models.PortfolioItem, error) {
	var item models.PortfolioItem
	err := r.db.QueryRow("SELECT id, title, short_description, project_url, github_url, showcase_image, sort_order, created_at, updated_at FROM portfolio_items WHERE id = ?", id).Scan(&item.ID, &item.Title, &item.ShortDescription, &item.ProjectURL, &item.GithubURL, &item.ShowcaseImage, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *PortfolioRepository) CreatePortfolioItem(item *models.PortfolioItem) error {
	err := r.db.QueryRow("INSERT INTO portfolio_items (title, short_description, project_url, github_url, showcase_image, sort_order) VALUES (?, ?, ?, ?, ?, ?) RETURNING id", item.Title, item.ShortDescription, item.ProjectURL, item.GithubURL, item.ShowcaseImage, item.SortOrder).Scan(&item.ID)
	return err
}

func (r *PortfolioRepository) UpdatePortfolioItem(item *models.PortfolioItem) error {
	_, err := r.db.Exec("UPDATE portfolio_items SET title = ?, short_description = ?, project_url = ?, github_url = ?, showcase_image = ?, sort_order = ? WHERE id = ?", item.Title, item.ShortDescription, item.ProjectURL, item.GithubURL, item.ShowcaseImage, item.SortOrder, item.ID)
	return err
}

func (r *PortfolioRepository) DeletePortfolioItem(id int64) error {
	result, err := r.db.Exec("DELETE FROM portfolio_items WHERE id = ?", id)
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