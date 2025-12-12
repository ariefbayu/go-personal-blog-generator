package models

import "time"

type PortfolioItem struct {
    ID               int64     `db:"id" json:"id"`
    Title            string    `db:"title" json:"title"`
    ShortDescription string    `db:"short_description" json:"short_description"`
    ProjectURL       string    `db:"project_url" json:"project_url"`
    GithubURL        string    `db:"github_url" json:"github_url"`
    ShowcaseImage    string    `db:"showcase_image" json:"showcase_image"`
    SortOrder        int       `db:"sort_order" json:"sort_order"`
    CreatedAt        time.Time `db:"created_at" json:"created_at"`
    UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}