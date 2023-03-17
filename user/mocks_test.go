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

// PasswordManager mock
type mockPasswordManager struct {
	password string
	err      error
}

func (m *mockPasswordManager) GenerateHash(password string) (string, error) {
	return m.password, m.err
}

func (m *mockPasswordManager) CompareHashWith(hashedPassword, password string) error {
	return m.err
}

// Client mock form mailer client
type mockMailerClient struct {
	err error
}

func (m *mockMailerClient) SendActivationEmail(email, name, activationToken string) error {
	return m.err
}

var _ user.Validator = mockValidator{}

// Validator mock
type mockValidator struct {
	err error
}

func (m mockValidator) ValidateRegistration(u *user.UserRegistrationInfo) error {
	return m.err
}

func (m mockValidator) ValidateLogin(u *user.UserLoginInfo) error {
	return m.err
}
