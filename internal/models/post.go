package models

import "time"

type Post struct {
	ID            int64     `db:"id" json:"id"`
	Title         string    `db:"title" json:"title"`
	Slug          string    `db:"slug" json:"slug"`
	Content       string    `db:"content" json:"content"`
	Tags          string    `db:"tags" json:"tags"`
	FeaturedImage string    `db:"featured_image" json:"featuredImage"`
	Published     bool      `db:"published" json:"published"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
