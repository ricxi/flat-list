package user

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	registerPayload := `{"firstName": "Michael", "lastName": "Scott", "email": "michaelscott@dundermifflin.com", "password": "1234"}`

	req := httptest.NewRequest(http.MethodPost, "/v1/user/register", strings.NewReader(registerPayload))
	w := httptest.NewRecorder()

	s := MockService{userID: "abcdefg", err: nil}
	handler := httpHandler{service: &s}
	handler.handleLogin(w, req)

}
