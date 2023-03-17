package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockService struct {
	userID   string
	userInfo *UserInfo
	err      error
}

func (m *mockService) RegisterUser(ctx context.Context, user *UserRegistrationInfo) (string, error) {
	return m.userID, m.err
}

func (m *mockService) LoginUser(ctx context.Context, user *UserLoginInfo) (*UserInfo, error) {
	return m.userInfo, m.err
}

func TestRegisterUser(t *testing.T) {
	var body bytes.Buffer
	jsonStr := `{"firstName": "Michael", "lastName": "Scott", "email": "michaelscott@dundermifflin.com", "password": "1234"}`

	if err := json.NewEncoder(&body).Encode(jsonStr); err != nil {
		t.Fatal(err)

	}

	req := httptest.NewRequest(http.MethodPost, "/v1/user/register", &body)
	w := httptest.NewRecorder()

	s := mockService{userID: "abcdefg", err: nil}
	handler := httpHandler{service: &s}
	handler.handleLogin(w, req)

}
