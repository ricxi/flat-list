package user

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"log"
	"os"

	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(ctx context.Context, user *UserRegistrationInfo) (string, error)
	LoginUser(ctx context.Context, user *UserLoginInfo) (*UserInfo, error)
}

type service struct {
	repository      Repository
	passwordService PasswordService
	client          Client
}

func NewService(repository Repository, passwordService PasswordService, client Client) Service {
	return &service{
		repository:      repository,
		passwordService: passwordService,
		client:          client,
	}
}

func (s *service) RegisterUser(ctx context.Context, u *UserRegistrationInfo) (string, error) {
	hashedPassword, err := s.passwordService.GenerateHash(u.Password)
	if err != nil {
		log.Println(err)
		return "", err
	}

	createdAt := time.Now().In(time.UTC)

	u.Password = ""
	u.HashedPassword = string(hashedPassword)
	u.CreatedAt = &createdAt
	u.UpdatedAt = &createdAt
	u.Activated = false

	userID, err := s.repository.CreateUser(ctx, u)
	if err != nil {
		log.Println(err)
		return "", err
	}

	activationToken, err := generateActivationToken()
	if err != nil {
		return "", err
	}

	go func() {
		if err := s.client.SendActivationEmail(u.Email, u.FirstName, activationToken); err != nil {
			log.Println(err)
		}
	}()

	return userID, nil
}

func (s *service) LoginUser(ctx context.Context, u *UserLoginInfo) (*UserInfo, error) {
	uInfo, err := s.repository.GetUserByEmail(ctx, u.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidEmail
		}
		log.Println(err)
		return nil, err
	}

	if err := s.passwordService.CompareHashWith(uInfo.HashedPassword, u.Password); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidPassword
		}

		log.Println(err)
		return nil, err
	}

	uInfo.Password = ""
	uInfo.HashedPassword = ""

	token, err := generateJWT(jwt.MapClaims{
		"user_id": uInfo.ID,
		"email":   uInfo.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	uInfo.Token = token

	return uInfo, nil
}

func generateJWT(claims jwt.MapClaims) (string, error) {
	secretKey, found := os.LookupEnv("JWT_SECRET_KEY")
	if !found {
		return "", ErrMissingEnvs
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

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
