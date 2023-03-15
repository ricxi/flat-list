package user

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ricxi/flat-list/mailer/activate"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Service interface {
	RegisterUser(ctx context.Context, user *UserRegistrationInfo) (string, error)
	LoginUser(ctx context.Context, user *UserLoginInfo) (*UserInfo, error)
}

type service struct {
	repository      Repository
	passwordService PasswordService
}

func NewService(repository Repository, passwordService PasswordService) Service {
	return &service{
		repository:      repository,
		passwordService: passwordService,
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
		err := sendActivationEmailGRPC(u.Email, u.FirstName, activationToken)
		if err != nil {
			log.Println(err)
		}
		// if err := sendActivationEmail(u.Email, u.FirstName, activationToken); err != nil {
		// 	log.Println(err)
		// }
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

// sendActivationEmail calls the mailer api to send an email
// to a newly registered user
func sendActivationEmail(email, firstName, activationToken string) error {
	activationInfo := struct {
		From            string `json:"from"`
		To              string `json:"to"`
		FirstName       string `json:"firstName"`
		ActivationToken string `json:"activationToken"`
	}{
		From:            "the.team@flatlist.com",
		To:              email,
		FirstName:       firstName,
		ActivationToken: activationToken,
	}

	reqBody := new(bytes.Buffer)
	if err := json.NewEncoder(reqBody).Encode(&activationInfo); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:5000/v1/mailer/activate", reqBody)
	if err != nil {
		return err
	}

	c := http.Client{}

	_, err = c.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func sendActivationEmailGRPC(email, name, activationToken string) error {
	cc, err := grpc.Dial(":5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	c := activate.NewMailerServiceClient(cc)
	if _, err := c.SendEmail(context.Background(), &activate.Request{
		From:            "the.team@flat-list.com",
		To:              email,
		FirstName:       name,
		ActivationToken: activationToken,
	}); err != nil {
		return err
	}

	return nil
}
