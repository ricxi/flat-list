package main

import (
	"context"
	"database/sql"
)

type ActivationTokenInfo struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

// ensure that repository implements the Repository interface
var _ Repository = new(repository)

type Repository interface {
	InsertToken(ctx context.Context, info *ActivationTokenInfo) error
	GetTokens(ctx context.Context, userID string) ([]string, error)
}

type repository struct {
	db *sql.DB
}

// InsertToken inserts a new activation token for a given user based on their id
func (r *repository) InsertToken(ctx context.Context, info *ActivationTokenInfo) error {
	query := "INSERT INTO activation_tokens (token, user_id) VALUES ($1, $2)"

	_, err := r.db.ExecContext(ctx, query, info.Token, info.UserID)
	return err

}

// GetTokens returns a string of tokens that share teh same user id
func (r *repository) GetTokens(ctx context.Context, userID string) ([]string, error) {
	query := "SELECT token FROM activation_tokens WHERE activation_tokens.user_id = $1"

	rows, err := r.db.QueryContext(ctx, query, userID)
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
