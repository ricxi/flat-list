package main

import (
	"crypto/rand"
	"encoding/base32"
)

// generate an activation token that is used
// to validate a newly registered user's account
func generateActivationToken() (string, error) {
	tokenBytes := make([]byte, 16)

	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	token := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(tokenBytes)

	return token, nil
}
