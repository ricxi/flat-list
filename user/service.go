package user

import (
	"context"
	"errors"
	"log"

	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(ctx context.Context, user UserRegistrationInfo) (string, error)
	LoginUser(ctx context.Context, user UserLoginInfo) (*UserInfo, error)
	ActivateUser(ctx context.Context, activationToken string) error
	RestartActivation(ctx context.Context, u UserLoginInfo) error
	Authenticate(ctx context.Context, signedJWT string) (string, error)
}

// service is instantiated using a builder (see builder.go file)
type service struct {
	repository Repository
	mailer     MailerClient
	password   PasswordManager
	validate   Validator
	token      TokenClient
}

func (s *service) RegisterUser(ctx context.Context, u UserRegistrationInfo) (string, error) {
	if err := s.validate.Registration(u); err != nil {
		return "", err
	}

	hashedPassword, err := s.password.GenerateHash(u.Password)
	if err != nil {
		log.Println(err)
		return "", err
	}
	u.Password = ""
	u.HashedPassword = string(hashedPassword)

	createdAt := time.Now().In(time.UTC)
	u.CreatedAt = &createdAt
	u.UpdatedAt = &createdAt

	u.Activated = false

	userID, err := s.repository.CreateUser(ctx, u)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// This line makes a grpc call to an external api
	activationToken, err := s.token.CreateActivationToken(context.Background(), userID)
	if err != nil {
		log.Println(err)
		return "", err
	}

	go func() {
		// send an activation email if a token is successfully generated
		if err := s.mailer.SendActivationEmail(u.Email, u.FirstName, activationToken); err != nil {
			log.Println(err)
		}
	}()

	return userID, nil
}

func (s *service) LoginUser(ctx context.Context, u UserLoginInfo) (*UserInfo, error) {
	if err := s.validate.Login(u); err != nil {
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

	// Should I compare the password before checking if the user has activated their account?
	if err := s.password.CompareHashWith(uInfo.HashedPassword, u.Password); err != nil {
		return nil, err
	}

	uInfo.Password = ""
	uInfo.HashedPassword = ""

	token, err := generateUserJWT(uInfo.ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	uInfo.Token = token

	return uInfo, nil
}

func (s *service) ActivateUser(ctx context.Context, activationToken string) error {
	userID, err := s.token.ValidateActivationToken(context.Background(), activationToken)
	if err != nil {
		log.Println(err)
		return err
	}

	var userUpdate UserInfo

	updateTime := time.Now().In(time.UTC)
	userUpdate.ID = userID
	userUpdate.Activated = true
	userUpdate.UpdatedAt = &updateTime

	if err := s.repository.UpdateUserByID(ctx, userUpdate); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// RestartActivation generates a new activation token and sends a new activation email to a user
// so long as they provide their email and a valid password (basically their login info).
// It is a route that is accessed by users who did receive a valid activation token or email due
// to unforseen or other cirumstances.
func (s *service) RestartActivation(ctx context.Context, u UserLoginInfo) error {
	if err := s.validate.Login(u); err != nil {
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

	if err := s.password.CompareHashWith(uInfo.HashedPassword, u.Password); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidPassword
		}

		log.Println(err)
		return err
	}

	activationToken, err := s.token.CreateActivationToken(context.Background(), uInfo.ID)
	if err != nil {
		log.Println(err)
		return err
	}

	go func() {
		if err := s.mailer.SendActivationEmail(uInfo.Email, uInfo.FirstName, activationToken); err != nil {
			log.Println(err)
		}
	}()

	return nil
}

// Authenticate receives a signed jwt, extracts user data from it, verifies the
// jwt, then checks that the user exists in the database and if their account
// has been activated. It returns the user's ID if everything is successful.
func (s *service) Authenticate(ctx context.Context, signedJWT string) (string, error) {
	if err := s.validate.NonEmptyString("jwt", signedJWT); err != nil {
		return "", err
	}

	// Should I add validation for UserClaims?
	var userClaims UserClaims
	if err := verifyUserJWT(signedJWT, &userClaims); err != nil {
		return "", err
	}

	uInfo, err := s.repository.GetUserByID(ctx, userClaims.UserID)
	if err != nil {
		return "", err
	}

	if !uInfo.Activated {
		return "", ErrUserNotActivated
	}

	return uInfo.ID, nil
}
