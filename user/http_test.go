package user

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// newRequestWithHeaders creates an *http.Request from the httptest package, and adds some headers to it.
func newRequestWithHeaders(method, target string, body io.Reader, headers map[string]string) *http.Request {
	r := httptest.NewRequest(method, target, body)
	for key, value := range headers {
		r.Header.Set(key, value)
	}
	return r
}

// newRequestWithJSONHeader creates an *http.Request from the httptest package,
// and adds a 'Content-Type: application/json' header to it.
func newRequestWithJSONHeader(method, target string, body io.Reader) *http.Request {
	return newRequestWithHeaders(
		method,
		target,
		body,
		map[string]string{"Content-Type": "application/json"},
	)
}

// These are technically integration tests because I'm going through the chi router...
// I could write tests for cases where the validation fails, and look through the registerUser method in service for more error cases
func TestHandleRegister(t *testing.T) {
	type req struct {
		method  string
		target  string
		body    string
		headers map[string]string
	}
	type expected struct {
		statusCode int
		body       string
	}
	testCases := []struct {
		name     string
		service  Service
		r        req
		expected expected
	}{
		{
			name: "Success",
			service: mockService{
				userID: "507f191e810c19729de860ea",
				err:    nil,
			},
			r: req{
				method: http.MethodPost,
				target: "/v1/user/register",
				body: `
				{
					"firstName": "Michael",
					"lastName": "Scott",
					"email": "michaelscott@dundermifflin.com",
					"password": "1234"
				}`,
				headers: map[string]string{
					"Content-Type": "application/json",
				},
			},
			expected: expected{
				statusCode: 201,
				body:       `{"id":"507f191e810c19729de860ea","success":true}`,
			},
		},
		{
			name: "FailDuplicateUserError",
			service: mockService{
				userID: "",
				err:    ErrDuplicateUser,
			},
			r: req{
				method: http.MethodPost,
				target: "/v1/user/register",
				body: `
				{
					"firstName": "Michael",
					"lastName": "Scott",
					"email": "michaelscott@dundermifflin.com",
					"password": "1234"
				}`,
				headers: map[string]string{
					"Content-Type": "application/json",
				},
			},
			expected: expected{
				statusCode: 400,
				body:       `{"error":"user already exists", "success":false}`,
			},
		},
		{
			name: "WrongContentTypeNoHeaders",
			service: mockService{
				userID: "",
				err:    nil,
			},
			r: req{
				method: http.MethodPost,
				target: "/v1/user/register",
				body: `
				{
					"firstName": "Michael",
					"lastName": "Scott",
					"email": "michaelscott@dundermifflin.com",
					"password": "1234"
				}`,
				headers: map[string]string{}, // this is intentionally left blank
			},
			expected: expected{
				statusCode: 400,
				body:       `{"error":"Invalid 'Content-Type': must be 'application/json'", "success":false}`,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHTTPHandler(tt.service)
			rr := httptest.NewRecorder() // I don't have to close the body for this or worry about a nil pointer dereference
			r := newRequestWithHeaders(tt.r.method, tt.r.target, strings.NewReader(tt.r.body), tt.r.headers)

			h.ServeHTTP(rr, r)
			assert.Equal(t, tt.expected.statusCode, rr.Code)
			assert.JSONEq(t, tt.expected.body, rr.Body.String())
		})
	}
}

// I wrote this differently than the TestHandleRegister
// function just to see which one would be more 'clear'
func TestHandleLogin(t *testing.T) {
	type expected struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name     string
		service  Service
		request  *http.Request
		expected expected
	}{
		{
			name: "Success",
			service: mockService{
				userInfo: &UserInfo{
					ID:        "60af7bf76c21d03b7c174f88",
					FirstName: "Michael",
					LastName:  "Scott",
					Email:     "michaelscott@dundermifflin.com",
					Password:  "1234",
					Activated: true,
					Token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjBhZjdiZjc2YzIxZDAzYjdjMTc0Zjg4IiwiaWF0IjoxNjIwNzgyMTQ2LCJleHAiOjE2MjA4NjkzNDZ9.P8IFC7zT_75fT4l4-uHsm8nVyLsDnsYbZwH3Pqkn8xU",
				},
				err: nil,
			},
			request: newRequestWithJSONHeader(
				http.MethodPost,
				"/v1/user/login",
				strings.NewReader(`
				{
					"email": "michaelscott@dundermifflin.com",
					"password": "1234"
				}`),
			),
			expected: expected{
				statusCode: 200,
				body: `
				{
					"user": {
						"id": "60af7bf76c21d03b7c174f88",
						"firstName": "Michael",
						"lastName": "Scott",
						"email": "michaelscott@dundermifflin.com",
						"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjBhZjdiZjc2YzIxZDAzYjdjMTc0Zjg4IiwiaWF0IjoxNjIwNzgyMTQ2LCJleHAiOjE2MjA4NjkzNDZ9.P8IFC7zT_75fT4l4-uHsm8nVyLsDnsYbZwH3Pqkn8xU"
					},
					"success": true
				}`,
			},
		},
		{
			// This returns ErrInvalidEmail, but why did I do this??
			name: "FailUserNotFound",
			service: mockService{
				userInfo: nil,
				err:      ErrInvalidEmail,
			},
			request: newRequestWithJSONHeader(
				http.MethodPost,
				"/v1/user/login",
				strings.NewReader(`
				{
					"email": "invalidemail@dundermifflin.com",
					"password": "1234"
				}
				`),
			),
			expected: expected{
				statusCode: 400,
				body: `
				{
					"error": "user with this email was not found",
					"success": false
				}`,
			},
		},
		{
			name: "UserNotActivated",
			service: mockService{
				userInfo: nil,
				err:      ErrUserNotActivated,
			},
			request: newRequestWithJSONHeader(
				http.MethodPost,
				"/v1/user/login",
				strings.NewReader(`
				{
					"email": "michaelscott@dundermifflin.com",
					"password": "1234"
				}
				`),
			),
			expected: expected{
				statusCode: 400,
				body: `
				{
					"error": "user has not activated their account",
					"success": false
				}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHTTPHandler(tt.service)
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, tt.request)

			assert.Equal(t, tt.expected.statusCode, rr.Code)
			assert.JSONEq(t, tt.expected.body, rr.Body.String())
		})
	}
}

func TestHandleActivate(t *testing.T) {
	type expected struct {
		statusCode int
		body       string
	}

	tests := []struct {
		name     string
		service  Service
		request  *http.Request
		expected expected
	}{
		{
			name: "Success",
			service: mockService{
				err: nil,
			},
			request: httptest.NewRequest(
				http.MethodPut,
				"/v1/user/activate/tokengoeshere",
				nil,
			),
			expected: expected{
				statusCode: 204,
				body:       "",
			},
		},
		{
			name: "UserNotFound",
			service: mockService{
				err: ErrUserNotFound,
			},
			request: httptest.NewRequest(
				http.MethodPut,
				"/v1/user/activate/tokengoeshere",
				nil,
			),
			expected: expected{
				statusCode: 400,
				body:       `{"error":"user not found", "success":false}`,
			},
		},
		{
			name: "NoTokenProvidedAsParameter",
			service: mockService{
				err: ErrUserNotFound,
			},
			request: httptest.NewRequest(
				http.MethodPut,
				"/v1/user/activate/", // trailing slash doesn't matter for chi
				nil,
			),
			expected: expected{
				statusCode: 404,
				body:       `{"error":"resource not found", "success":false}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHTTPHandler(tt.service)
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, tt.request)

			assert.Equal(t, tt.expected.statusCode, rr.Code)
			// This is kind of hacky, but I'll change it later
			if rr.Code >= 400 {
				assert.JSONEq(t, tt.expected.body, rr.Body.String())
			}
		})
	}
}

func TestHandleRestartActivation(t *testing.T) {
	type expected struct {
		statusCode int
		hasBody    bool
		body       string
	}

	testCases := []struct {
		name     string
		service  Service
		request  *http.Request
		expected expected
	}{
		{
			name: "Success",
			service: mockService{
				err: nil,
			},
			request: newRequestWithJSONHeader(
				http.MethodPost,
				"/v1/user/restart/activation",
				strings.NewReader(`
				{
					"email": "michaelscott@dundermifflin.com",
					"password": "1234"
				}
				`),
			),
			expected: expected{
				statusCode: 204,
				hasBody:    false,
				body:       "",
			},
		},
		{
			name: "InvalidPassword",
			service: mockService{
				err: ErrInvalidPassword,
			},
			request: newRequestWithJSONHeader(
				http.MethodPost,
				"/v1/user/restart/activation",
				strings.NewReader(`
				{
					"email": "michaelscott@dundermifflin.com",
					"password": ""
				}
				`),
			),
			expected: expected{
				statusCode: 400,
				hasBody:    true,
				body:       `{"error":"invalid password provided", "success":false}`,
			},
		},
	}

	for _, tt := range testCases {
		h := NewHTTPHandler(tt.service)
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, tt.request)

		assert.Equal(t, tt.expected.statusCode, rr.Code)

		if tt.expected.hasBody {
			assert.JSONEq(t, tt.expected.body, rr.Body.String())
		}
	}
}
