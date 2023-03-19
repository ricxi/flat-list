package user

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(ctx context.Context, user *UserRegistrationInfo) (string, error)
	LoginUser(ctx context.Context, user *UserLoginInfo) (*UserInfo, error)
	ActivateUser(ctx context.Context, activationToken string) error
}

type service struct {
	repository      Repository
	client          Client
	passwordManager PasswordManager
	v               Validator
}

func NewService(
	repository Repository,
	client Client,
	passwordManager PasswordManager,
	validator Validator,
) Service {
	return &service{
		repository:      repository,
		client:          client,
		passwordManager: passwordManager,
		v:               validator,
	}
}

func (s *service) RegisterUser(ctx context.Context, u *UserRegistrationInfo) (string, error) {
	if err := s.v.ValidateRegistration(u); err != nil {
		return "", err
	}

	hashedPassword, err := s.passwordManager.GenerateHash(u.Password)
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

	go func() {
		c := http.Client{Timeout: (3 * time.Second)}
		req, err := http.NewRequest(http.MethodPost, "http://localhost:5002/v1/token/activation/"+userID, nil)
		if err != nil {
			log.Println(err)
		}

		res, err := c.Do(req)
		if err != nil {
			log.Println(err)
		}

		defer res.Body.Close()
		activationTokenData := struct {
			Token string `json:"token"`
		}{}
		if err := json.NewDecoder(res.Body).Decode(&activationTokenData); err != nil {
			log.Println(err)
		}

		if activationTokenData.Token != "" {
			if err := s.client.SendActivationEmail(u.Email, u.FirstName, activationTokenData.Token); err != nil {
				log.Println(err)
			}
		}
	}()

	return userID, nil
}

func (s *service) LoginUser(ctx context.Context, u *UserLoginInfo) (*UserInfo, error) {
	if err := s.v.ValidateLogin(u); err != nil {
		return nil, err
	}

	uInfo, err := s.repository.GetUserByEmail(ctx, u.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidEmail
		}
		log.Println(err)
		return nil, err
	}

	if !uInfo.Activated {
		return nil, ErrUserNotActivated
	}

	if err := s.passwordManager.CompareHashWith(uInfo.HashedPassword, u.Password); err != nil {
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

func (s *service) ActivateUser(ctx context.Context, activationToken string) error {

	// search database for token
	// compare token
	// search database for user
	// set user's activate status to true

	return nil
}

func generateJWT(claims jwt.MapClaims) (string, error) {
	secretKey, found := os.LookupEnv("JWT_SECRET_KEY")
	if !found {
		return "", ErrMissingEnvs
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
