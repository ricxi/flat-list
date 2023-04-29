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

func (m *mockRepository) createUser(ctx context.Context, u UserRegistrationInfo) (string, error) {
	return m.userID, m.err
}

func (m *mockRepository) getUserByEmail(ctx context.Context, email string) (*UserInfo, error) {
	return m.user, m.err
}

func (m *mockRepository) getUserByID(ctx context.Context, id string) (*UserInfo, error) {
	return m.user, m.err
}

func (m *mockRepository) updateUserByID(ctx context.Context, u UserInfo) error {
	return m.err
}

// Service mock
type mockService struct {
	userID   string
	userInfo *UserInfo
	err      error
}

func (m mockService) registerUser(ctx context.Context, user UserRegistrationInfo) (string, error) {
	return m.userID, m.err
}

func (m mockService) loginUser(ctx context.Context, user UserLoginInfo) (*UserInfo, error) {
	return m.userInfo, m.err
}

func (m mockService) activateUser(ctx context.Context, activationToken string) error {
	return m.err
}

func (m mockService) restartActivation(ctx context.Context, u UserLoginInfo) error {
	return m.err
}

func (m mockService) authenticate(ctx context.Context, signedJWT string) (string, error) {
	return m.userID, m.err
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

func (m *mockValidator) NonEmptyString(name, field string) error {
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
