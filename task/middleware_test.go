package task

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetAuthToken(t *testing.T) {
	// The Authorization token should be set and a Bearer token should be found
	t.Run("SuccessSetAndBearerTokenFound", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		r := httptest.NewRequest("", "/", nil)
		// this is incorrect; i think it should be a jwt
		r.Header.Set("Authorization", "Bearer 507f191e810c19729de860ea")

		actual, err := getAuthToken(r)
		require.NoError(err)

		expected := "507f191e810c19729de860ea"
		assert.Equal(expected, actual)
		assert.True(primitive.IsValidObjectID(actual))
	})

	t.Run("NoAuthorizationHeader", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		expectedErrString := "auth header is empty or missing"

		r := httptest.NewRequest("", "/", nil)

		actual, err := getAuthToken(r)
		require.Error(err)
		require.Equal("", actual)

		assert.EqualError(err, expectedErrString)
	})

	t.Run("InvalidBearerTokenValueIsEmptyString", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		expectedErrString := "invalid Bearer token"

		r := httptest.NewRequest("", "/", nil)
		r.Header.Set("Authorization", "Bearer")

		actual, err := getAuthToken(r)
		require.Error(err)
		require.Equal("", actual)
		assert.EqualError(err, expectedErrString)
	})
}

func TestMiddlewareAuthenticate(t *testing.T) {
	mockAuthHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(UserIDCtxKey)
		if userID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"user id not found"}`))
		}

		w.Write([]byte(`{"success":true}`))
	})

	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"userId":"507f191e810c19729de860ea"}`))
		}))
		defer ts.Close()

		expected := `{"success":true}`

		rr := httptest.NewRecorder()

		m := &Middleware{
			AuthEndpoint: ts.URL,
		}

		auth := m.Authenticate(mockAuthHandler)

		r := httptest.NewRequest(http.MethodPost, "/", nil)
		r.Header.Set("Authorization", "Bearer jwttoken")

		auth.ServeHTTP(rr, r)

		require.Equal(http.StatusOK, rr.Code)
		require.NotEmpty(rr.Body)
		assert.JSONEq(expected, rr.Body.String())
	})

	// I didn't define a test server here because the
	// middleware should return before reaching that point (but I should look here if I ever get a nil pointer dereference error)
	t.Run("FailNoAuthHeader", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		expected := `{"error":"auth header is empty or missing"}`

		rr := httptest.NewRecorder()

		auth := (&Middleware{}).Authenticate(mockAuthHandler)

		r := httptest.NewRequest(http.MethodPost, "/", nil)

		auth.ServeHTTP(rr, r)

		require.Equal(http.StatusUnauthorized, rr.Code)
		require.NotEmpty(rr.Body)
		assert.JSONEq(expected, rr.Body.String())
	})

	t.Run("FailErrorFromCallToAPI", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"no token provided"}`))
		}))
		defer ts.Close()

		expected := `{"error":"unable to authorize user"}`

		rr := httptest.NewRecorder()

		m := &Middleware{
			AuthEndpoint: ts.URL,
		}

		auth := m.Authenticate(mockAuthHandler)

		r := httptest.NewRequest(http.MethodPost, "/", nil)
		r.Header.Set("Authorization", "Bearer jwttoken")

		auth.ServeHTTP(rr, r)

		require.Equal(http.StatusUnauthorized, rr.Code)
		require.NotEmpty(rr.Body)
		assert.JSONEq(expected, rr.Body.String())
	})
}
