package model

import "time"

// Post represents the metadata for a blog post stored in the database.
type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	ContentPath string    `json:"-"` // Path to the markdown file in storage (local or S3)
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PostHistory tracks changes to a post.
type PostHistory struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	Version   int       `json:"version"`
	ContentPath string    `json:"content_path"` // Path to the historical version of the markdown file
	CreatedAt time.Time `json:"created_at"`
}