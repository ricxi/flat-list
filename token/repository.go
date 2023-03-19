package token

import (
	"context"
	"database/sql"
)

type ActivationTokenInfo struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return Repository{db: db}
}

func Connect(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// InsertActivationToken inserts a new activation token for a given user based on their id
func (r *Repository) InsertActivationToken(ctx context.Context, info *ActivationTokenInfo) error {
	query := "INSERT INTO activation_tokens (token, user_id) VALUES ($1, $2)"

	_, err := r.db.ExecContext(ctx, query, info.Token, info.UserID)

	return err
}

// GetUserID receives an activation token and searches for the user id associated with it
func (r *Repository) GetUserID(ctx context.Context, activationToken string) (string, error) {
	query := "SELECT user_id FROM activation_tokens WHERE activation_tokens.token = $1"

	row := r.db.QueryRowContext(ctx, query, activationToken)
	var userID string
	if err := row.Scan(&userID); err != nil {
		return "", err
	}

	return userID, nil
}
