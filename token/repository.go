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
	DB *sql.DB
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

// InsertToken inserts a new activation token for a given user based on their id
func (r *Repository) InsertToken(ctx context.Context, info *ActivationTokenInfo) error {
	query := "INSERT INTO activation_tokens (token, user_id) VALUES ($1, $2)"

	_, err := r.DB.ExecContext(ctx, query, info.Token, info.UserID)
	return err

}

// GetTokens returns a string of tokens that share the same user id
func (r *Repository) GetTokens(ctx context.Context, userID string) ([]string, error) {
	query := "SELECT token FROM activation_tokens WHERE activation_tokens.user_id = $1"

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return tokens, nil
}
