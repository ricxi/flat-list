package main

import (
	"encoding/base32"
	"math/rand"
)

type ActivationToken string

func generateActivationToken() (ActivationToken, error) {
	tokenBytes := make([]byte, 16)

	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	token := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(tokenBytes)

	return ActivationToken(token), nil
}
