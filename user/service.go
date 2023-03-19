package user

import (
	"context"
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
	ActivateUser(ctx context.Context, activationToken string) error
}

type service struct {
	repository      Repository
	client          Client
	passwordManager PasswordManager
	v               Validator
	tc              *tokenClient
}

func NewService(
	repository Repository,
	client Client,
	passwordManager PasswordManager,
	validator Validator,
) Service {
	s := service{
		repository:      repository,
		client:          client,
		passwordManager: passwordManager,
		v:               validator,
	}

	tc, err := NewTokenClient("5003")
	if err != nil {
		log.Fatalln(err)
	}
	s.tc = tc

	return &s
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

	activationToken, err := s.tc.CreateActivationToken(context.Background(), userID)
	if err != nil {
		log.Println(err)
		return "", err
	}

	go func() {
		// send an activation email if a token is successfully generated
		if err := s.client.SendActivationEmail(u.Email, u.FirstName, activationToken); err != nil {
			log.Println(err)
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
	userID, err := s.tc.ValidateActivationToken(context.Background(), activationToken)
	if err != nil {
		log.Println(err)
		return err
	}

	var userUpdate UserInfo

	updateTime := time.Now().In(time.UTC)
	userUpdate.ID = userID
	userUpdate.Activated = true
	userUpdate.UpdatedAt = &updateTime

	if err := s.repository.UpdateUserByID(ctx, &userUpdate); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// RestartActivation generates a new activation token and sends a new activation email to a user
// so long as they provide their email and a valid password (basically their login info).
// It is a route that is accessed by users who did receive a valid activation token or email due
// to unforseen or other cirumstances.
func (s *service) RestartActivation(ctx context.Context, u *UserLoginInfo) error {
	if err := s.v.ValidateLogin(u); err != nil {
		return err
	}

	uInfo, err := s.repository.GetUserByEmail(ctx, u.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return ErrInvalidEmail
		}
		log.Println(err)
		return err
	}

	if err := s.passwordManager.CompareHashWith(uInfo.HashedPassword, u.Password); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidPassword
		}

		log.Println(err)
		return err
	}

	activationToken, err := s.tc.CreateActivationToken(context.Background(), uInfo.ID)
	if err != nil {
		log.Println(err)
		return err
	}

	go func() {
		if err := s.client.SendActivationEmail(uInfo.Email, uInfo.FirstName, activationToken); err != nil {
			log.Println(err)
		}
	}()

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
