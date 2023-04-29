package user

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// These are technically integration tests because I'm going through the chi router...
// I could write tests for cases where the validation fails, and look through the registerUser method in service for more error cases
func TestRegisterUser(t *testing.T) {
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
			r := httptest.NewRequest(tt.r.method, tt.r.target, strings.NewReader(tt.r.body))
			for key, value := range tt.r.headers {
				r.Header.Set(key, value)
			}

			h.ServeHTTP(rr, r)
			assert.Equal(t, tt.expected.statusCode, rr.Code)
			assert.JSONEq(t, tt.expected.body, rr.Body.String())
		})
	}
}
