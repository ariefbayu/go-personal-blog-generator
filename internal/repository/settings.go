package repository

import (
	"database/sql"
	"time"
)

type Settings struct {
	ID                int       `json:"id"`
	SiteName          string    `json:"site_name"`
	ShowPortfolioMenu bool      `json:"show_portfolio_menu"`
	ShowPostsMenu     bool      `json:"show_posts_menu"`
	MenuOrder         string    `json:"menu_order"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type SettingsRepository struct {
	db *sql.DB
}

func NewSettingsRepository(db *sql.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) GetSettings() (*Settings, error) {
	settings := &Settings{}
	err := r.db.QueryRow(`
		SELECT id, site_name, show_portfolio_menu, show_posts_menu, menu_order, created_at, updated_at
		FROM settings WHERE id = 1
	`).Scan(
		&settings.ID,
		&settings.SiteName,
		&settings.ShowPortfolioMenu,
		&settings.ShowPostsMenu,
		&settings.MenuOrder,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return settings, nil
}

func (r *SettingsRepository) UpdateSettings(settings *Settings) error {
	settings.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE settings SET
			site_name = ?,
			show_portfolio_menu = ?,
			show_posts_menu = ?,
			menu_order = ?,
			updated_at = ?
		WHERE id = 1
	`,
		settings.SiteName,
		settings.ShowPortfolioMenu,
		settings.ShowPostsMenu,
		settings.MenuOrder,
		settings.UpdatedAt,
	)
	return err
}
