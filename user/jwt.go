package user

import (
	"os"

	"github.com/golang-jwt/jwt/v4"
)

func generateJWT(claims jwt.MapClaims) (string, error) {
	secretKey, found := os.LookupEnv("JWT_SECRET_KEY")
	if !found {
		return "", ErrMissingEnvs
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
