package postgres

import (
	"database/sql"
	"go-blog/internal/model"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(user *model.User) (*model.User, error) {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := s.db.QueryRow(query, user.Username, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserStore) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	// We must select the password_hash to compare it later.
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE email = $1`
	err := s.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password, // The password_hash from DB is scanned into the Password field.
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err // This will correctly return sql.ErrNoRows if user not found.
	}
	return user, nil
}

func (s *UserStore) GetByID(id int) (*model.User, error) {
	// Implementation for getting a user by ID
	// This is a placeholder; actual implementation would query the database.
	return nil, nil
}