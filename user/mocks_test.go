package user

import (
	"context"
)

var _ Repository = &mockRepository{}

// Repository mock
type mockRepository struct {
	userID string
	user   *UserInfo
	err    error
}

func (m *mockRepository) CreateUser(ctx context.Context, u UserRegistrationInfo) (string, error) {
	return m.userID, m.err
}

func (m *mockRepository) GetUserByEmail(ctx context.Context, email string) (*UserInfo, error) {
	return m.user, m.err
}

func (m *mockRepository) GetUserByID(ctx context.Context, id string) (*UserInfo, error) {
	return m.user, m.err
}

func (m *mockRepository) UpdateUserByID(ctx context.Context, u UserInfo) error {
	return m.err
}

// Service mock
type MockService struct {
	userID   string
	userInfo *UserInfo
	err      error
}

func (m *MockService) RegisterUser(ctx context.Context, user UserRegistrationInfo) (string, error) {
	return m.userID, m.err
}

func (m *MockService) LoginUser(ctx context.Context, user *UserLoginInfo) (*UserInfo, error) {
	return m.userInfo, m.err
}

func (m *MockService) ActivateUser(ctx context.Context, activationToken string) error {
	return m.err
}

func (m *MockService) RestartActivation(ctx context.Context, u *UserLoginInfo) error {
	return m.err
}

// PasswordManager mock
type mockPasswordManager struct {
	hashedPassword string
	err            error
}

func (m *mockPasswordManager) GenerateHash(password string) (string, error) {
	return m.hashedPassword, m.err
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

var _ Validator = &mockValidator{}

// Validator mock
type mockValidator struct {
	err error
}

func (m *mockValidator) Registration(u UserRegistrationInfo) error {
	return m.err
}

func (m *mockValidator) Login(u UserLoginInfo) error {
	return m.err
}

type mockTokenClient struct {
	mockActivationToken string
	mockUserID          string
	err                 error
}

func (m *mockTokenClient) CreateActivationToken(ctx context.Context, userID string) (string, error) {
	return m.mockActivationToken, m.err
}

func (m *mockTokenClient) ValidateActivationToken(ctx context.Context, activationToken string) (string, error) {
	return m.mockUserID, m.err
}
