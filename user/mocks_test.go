package user_test

import (
	"context"

	"github.com/ricxi/flat-list/user"
)

var _ user.Repository = &mockRepository{}

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

type mockPasswordManager struct {
	password string
	err      error
}

// GenerateHash creates a hash for a given password
func (m *mockPasswordManager) GenerateHash(password string) (string, error) {
	return m.password, m.err
}

func (m *mockPasswordManager) CompareHashWith(hashedPassword, password string) error {
	return m.err
}

type mockMailerClient struct {
	err error
}

func (m *mockMailerClient) SendActivationEmail(email, name, activationToken string) error {
	return m.err
}
