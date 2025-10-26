package postgres

import (
	"database/sql"
	"go-blog/internal/model"
)

type PostStore struct {
	db *sql.DB
}

func NewPostStore(db *sql.DB) *PostStore {
	return &PostStore{db: db}
}

func (s *PostStore) Create(post *model.Post) (*model.Post, error) {
	query := `INSERT INTO posts (user_id, title, version) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := s.db.QueryRow(query, post.UserID, post.Title, post.Version).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostStore) Update(post *model.Post) (*model.Post, error) {
	query := `UPDATE posts SET title = $1, content_path = $2, version = $3, updated_at = NOW() WHERE id = $4 RETURNING updated_at`
	err := s.db.QueryRow(query, post.Title, post.ContentPath, post.Version, post.ID).Scan(&post.UpdatedAt)
	return post, err
}

func (s *PostStore) GetByID(id int) (*model.Post, error) {
	post := &model.Post{}
	query := `SELECT id, user_id, title, content_path, version, created_at, updated_at FROM posts WHERE id = $1`
	err := s.db.QueryRow(query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.ContentPath,
		&post.Version,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostStore) List(limit, offset int) ([]*model.Post, error) {
	query := `
		SELECT id, user_id, title, content_path, version, created_at, updated_at 
		FROM posts 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2`

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		post := &model.Post{}
		if err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.ContentPath,
			&post.Version, &post.CreatedAt, &post.UpdatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostStore) CreateHistory(history *model.PostHistory) error {
	query := `INSERT INTO post_history (post_id, version, content_path) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(query, history.PostID, history.Version, history.ContentPath)
	return err
}