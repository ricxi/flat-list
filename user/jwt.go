package user

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrInvalidJWT          = errors.New("invalid jwt")
	ErrInvalidJWTSignature = errors.New("invalid jwt signature")
)

// generateUserJWT creates a signed jwt with the user's ID.
func generateUserJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	return generateJWT(claims)
}

func generateJWT(claims jwt.Claims) (string, error) {
	secretKey, found := os.LookupEnv("JWT_SECRET_KEY")
	if !found {
		return "", ErrMissingEnvs
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedJWT, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("error signing token: %w ", err)
	}

	return signedJWT, nil
}

func verifyJWT(signedJWT string) error {
	secretKey, found := os.LookupEnv("JWT_SECRET_KEY")
	if !found {
		return ErrMissingEnvs
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(
		signedJWT,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return ErrInvalidJWTSignature
		}
	}

	if !token.Valid {
		return ErrInvalidJWT
	}

	return nil
}
