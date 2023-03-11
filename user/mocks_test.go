package user_test

import (
	"context"

	"github.com/ricxi/flat-list/user"
)

type Repository struct {
	UserID string
	user   *user.UserInfo
	Err    error
}

func (m *Repository) CreateUser(ctx context.Context, u *user.UserRegistrationInfo) (string, error) {
	return m.UserID, m.Err
}

func (m *Repository) GetUserByEmail(ctx context.Context, email string) (*user.UserInfo, error) {
	return m.user, m.Err
}

type mockPasswordService struct {
	password string
	err      error
}

// GenerateHash creates a hash for a given password
func (m *mockPasswordService) GenerateHash(password string) (string, error) {
	return m.password, m.err
}

func (m *mockPasswordService) CompareHashWith(hashedPassword, password string) error {
	return m.err
}
