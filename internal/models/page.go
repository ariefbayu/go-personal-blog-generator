package models

import "time"

type Page struct {
    ID         int64     `db:"id" json:"id"`
    Title      string    `db:"title" json:"title"`
    Slug       string    `db:"slug" json:"slug"`
    Content    string    `db:"content" json:"content"`
    ShowInNav  bool      `db:"show_in_nav" json:"show_in_nav"`
    SortOrder  int       `db:"sort_order" json:"sort_order"`
    CreatedAt  time.Time `db:"created_at" json:"created_at"`
    UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}