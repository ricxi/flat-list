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

type UserClaims struct {
	jwt.MapClaims
	UserID string
}

// generateUserJWT creates a signed jwt with the user's ID.
func generateUserJWT(userID string) (string, error) {
	userClaims := UserClaims{
		UserID: userID,
		MapClaims: jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	return generateJWT(userClaims)
}

// generateJWT creates a signed jwt and receives any type
// that embeds a struct that implements the jwt.Claims interface
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

func verifyUserJWT(signedJWT string, userClaims *UserClaims) error {
	secretKey, found := os.LookupEnv("JWT_SECRET_KEY")
	if !found {
		return ErrMissingEnvs
	}

	token, err := jwt.ParseWithClaims(
		signedJWT,
		userClaims,
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
