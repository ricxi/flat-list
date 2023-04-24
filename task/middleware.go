package task

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	res "github.com/ricxi/flat-list/shared/response"
)

type Middleware struct {
	AuthEndpoint string
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getAuthToken(r)
		if err != nil {
			res.SendErrorJSON(w, err.Error(), http.StatusUnauthorized)
			return
		}

		reqBody := new(bytes.Buffer)
		if err := json.NewEncoder(reqBody).Encode(map[string]string{"token": token}); err != nil {
			res.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequest(http.MethodPost, m.AuthEndpoint, reqBody)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			res.SendErrorJSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			res.SendErrorJSON(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if resp.StatusCode != http.StatusOK {
			res.SendErrorJSON(w, "unable to authorize user", http.StatusUnauthorized)
			return
		}

		userID, err := getUserID(resp)
		if err != nil {
			res.SendErrorJSON(w, "unable to authorize user", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDCtxKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getAuthToken obtains the 'Authorization' header from
// the request, which should be a 'Bearer' token whose value
// is a valid mongo object id.
func getAuthToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("auth header is empty or missing")
	}

	slicedAuthHeader := strings.Split(authHeader, " ")
	if len(slicedAuthHeader) != 2 || slicedAuthHeader[0] != "Bearer" {
		return "", errors.New("invalid Bearer token")
	}

	return slicedAuthHeader[1], nil
}

// getUserID obtains the user id in the response body
func getUserID(resp *http.Response) (string, error) {
	authData := make(map[string]string)
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&authData); err != nil {
		return "", err
	}

	userID, ok := authData["userId"]
	if !ok || userID == "" {
		return "", errors.New("invalid or missing user id")
	}

	return userID, nil
}
