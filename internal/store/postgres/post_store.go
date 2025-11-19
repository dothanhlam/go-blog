package postgres

import (
	"database/sql"
	"go-blog/internal/model"

	"github.com/lib/pq"
)

type PostStore struct {
	db *sql.DB
}

func NewPostStore(db *sql.DB) *PostStore {
	return &PostStore{db: db}
}

func (s *PostStore) Create(post *model.Post) (*model.Post, error) {
	query := `INSERT INTO posts (user_id, title, sub_title, image, tags, version) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`
	err := s.db.QueryRow(query, post.UserID, post.Title, post.SubTitle, post.Image, pq.Array(post.Tags), post.Version).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostStore) Update(post *model.Post) (*model.Post, error) {
	query := `UPDATE posts SET title = $1, sub_title = $2, image = $3, tags = $4, content_path = $5, version = $6, updated_at = NOW() WHERE id = $7 RETURNING updated_at`
	err := s.db.QueryRow(query, post.Title, post.SubTitle, post.Image, pq.Array(post.Tags), post.ContentPath, post.Version, post.ID).Scan(&post.UpdatedAt)
	return post, err
}

func (s *PostStore) GetByID(id int) (*model.Post, error) {
	post := &model.Post{}
	query := `SELECT id, user_id, title, sub_title, image, tags, content_path, version, created_at, updated_at FROM posts WHERE id = $1`
	err := s.db.QueryRow(query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.SubTitle,
		&post.Image,
		pq.Array(&post.Tags),
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
		SELECT id, user_id, title, sub_title, image, tags, content_path, version, created_at, updated_at 
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
			&post.ID, &post.UserID, &post.Title, &post.SubTitle, &post.Image, pq.Array(&post.Tags), &post.ContentPath,
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

func (s *PostStore) Search(query string, limit, offset int) ([]*model.Post, error) {
	// plainto_tsquery is used for user-provided search terms.
	// It's safer and handles multiple words well.
	sqlQuery := `
		SELECT id, user_id, title, sub_title, image, tags, content_path, version, created_at, updated_at
		FROM posts
		WHERE title_tsv @@ plainto_tsquery('english', $1)
		ORDER BY ts_rank(title_tsv, plainto_tsquery('english', $1)) DESC
		LIMIT $2 OFFSET $3`

	rows, err := s.db.Query(sqlQuery, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		post := &model.Post{}
		if err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.SubTitle, &post.Image, pq.Array(&post.Tags), &post.ContentPath,
			&post.Version, &post.CreatedAt, &post.UpdatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, rows.Err()
}