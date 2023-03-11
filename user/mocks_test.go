package user_test

import (
	"context"

	"github.com/ricxi/flat-list/user"
)

// Repository mock
type mockRepository struct {
	userID string
	user   *user.UserInfo
	err    error
}

func (m *mockRepository) CreateUser(ctx context.Context, u *user.UserRegistrationInfo) (string, error) {
	return m.userID, m.err
}

func (m *mockRepository) GetUserByEmail(ctx context.Context, email string) (*user.UserInfo, error) {
	return m.user, m.err
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
