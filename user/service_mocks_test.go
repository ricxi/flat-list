package user_test

import (
	"context"

	"github.com/ricxi/flat-list/user"
)

type Repository struct {
	UserID string
	User   *user.UserInfo
	Err    error
}

func (m *Repository) CreateUser(ctx context.Context, u *user.UserRegistrationInfo) (string, error) {
	return m.UserID, m.Err
}

func (m *Repository) GetUserByEmail(ctx context.Context, email string) (*user.UserInfo, error) {
	return m.User, m.Err
}
