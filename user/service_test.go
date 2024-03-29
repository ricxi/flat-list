package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// I used a code generator to create the
// boilerplate for the table-driven tests
// this time, but I wrote the tests myself.
// Note, the logger will log to the terminal
// when some of these tests are run for error cases.
// ! No tests are written for create activation token
func Test_Service_RegisterUser(t *testing.T) {
	type fields struct {
		repository Repository
		mailer     MailerClient
		password   PasswordManager
		validate   Validator
		token      TokenClient
	}
	type args struct {
		ctx context.Context
		u   UserRegistrationInfo
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		expUserID string
		expErrStr string
	}{
		{
			name: "SuccessPassValidation",
			fields: fields{
				repository: &mockRepository{
					userID: "5ef7fdd91c19e3222b41b839",
					err:    nil,
				},
				mailer: &mockMailerClient{
					// The mailer's methods are called concurrently, so returning an error will not affect the service at all
					err: nil,
				},
				password: &mockPasswordManager{hashedPassword: "", err: nil},
				validate: &validator{}, // This is not the mock
				token:    &mockTokenClient{mockActivationToken: "", mockUserID: "", err: nil},
			},
			args: args{
				ctx: context.Background(),
				u: UserRegistrationInfo{
					Email:    "michaelscott@dundermifflin.com",
					Password: "1234",
				},
			},
		},
		{
			name: "SuccessWithAllFields",
			fields: fields{
				repository: &mockRepository{
					userID: "5ef7fdd91c19e3222b41b839",
					err:    nil,
				},
				mailer: &mockMailerClient{
					// The mailer's methods are called concurrently, so returning an error will not affect the service at all
					err: nil,
				},
				password: &mockPasswordManager{
					hashedPassword: "",
					err:            nil,
				},
				// This is not the mock
				validate: &validator{},
				token:    &mockTokenClient{mockActivationToken: "", mockUserID: "", err: nil},
			},
			args: args{
				ctx: context.Background(),
				u: UserRegistrationInfo{
					FirstName: "Michael",
					LastName:  "Scott",
					Email:     "michaelscott@dundermifflin.com",
					Password:  "1234",
				},
			},
		},
		{
			name: "FailErrMissingFieldEmail",
			fields: fields{
				repository: &mockRepository{
					userID: "5ef7fdd91c19e3222b41b839",
					err:    nil,
				},
				mailer: &mockMailerClient{
					// The mailer's methods are called concurrently, so returning an error will not affect the service at all
					err: nil,
				},
				password: &mockPasswordManager{hashedPassword: "", err: nil},
				// This is not the mock
				validate: &validator{},
				token:    &mockTokenClient{mockActivationToken: "", mockUserID: "", err: nil},
			},
			args: args{
				ctx: context.Background(),
				u: UserRegistrationInfo{
					FirstName: "Michael",
					LastName:  "Scott",
					Password:  "1234",
				},
			},
			expErrStr: "missing field is required: email",
		},
		{
			name: "FailErrMissingFieldPassword",
			fields: fields{
				repository: &mockRepository{
					userID: "5ef7fdd91c19e3222b41b839",
					err:    nil,
				},
				mailer: &mockMailerClient{
					// The mailer's methods are called concurrently, so returning an error will not affect the service at all
					err: errors.New("dummy error"),
				},
				password: &mockPasswordManager{hashedPassword: "", err: nil},
				// This is not the mock
				validate: &validator{},
				token:    &mockTokenClient{mockActivationToken: "", mockUserID: "", err: nil},
			},
			args: args{
				ctx: context.Background(),
				u: UserRegistrationInfo{
					FirstName: "Michael",
					LastName:  "Scott",
					Email:     "michaelscott@dundermifflin.com",
				},
			},
			expErrStr: "missing field is required: password",
		},
		{
			name: "FailPasswordGenerationError",
			fields: fields{
				repository: &mockRepository{
					// code should not reach here
					userID: "",
				},
				mailer: &mockMailerClient{
					// The mailer's methods are called concurrently, so returning an error will not affect the service at all
					err: errors.New("dummy error"),
				},
				password: &mockPasswordManager{
					hashedPassword: "does not matter because it will never be exposed outside of service",
					err:            bcrypt.ErrPasswordTooLong,
				},
				// This is not the mock
				validate: &validator{},
				token:    &mockTokenClient{mockActivationToken: "", mockUserID: "", err: nil},
			},
			args: args{
				ctx: context.Background(),
				u: UserRegistrationInfo{
					FirstName: "Michael",
					LastName:  "Scott",
					Email:     "michaelscott@dundermifflin.com",
					Password:  "1234",
				},
			},
			expErrStr: "bcrypt: password length exceeds 72 bytes",
		},
		{
			name: "FailDuplicateUser",
			fields: fields{
				repository: &mockRepository{
					userID: "",
					err:    ErrDuplicateUser,
				},
				mailer: &mockMailerClient{
					// The mailer's methods are called concurrently, so returning an error will not affect the service at all
					err: errors.New("dummy error"),
				},
				password: &mockPasswordManager{
					hashedPassword: "does not matter because it will never be exposed outside of service",
					err:            nil,
				},
				// This is not the mock
				validate: &validator{},
				token:    &mockTokenClient{mockActivationToken: "", mockUserID: "", err: nil},
			},
			args: args{
				ctx: context.Background(),
				u: UserRegistrationInfo{
					FirstName: "Michael",
					LastName:  "Scott",
					Email:     "michaelscott@dundermifflin.com",
					Password:  "1234",
				},
			},
			expErrStr: "user already exists",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			s := &service{
				repository: tt.fields.repository,
				mailer:     tt.fields.mailer,
				password:   tt.fields.password,
				validate:   tt.fields.validate,
				token:      tt.fields.token,
			}

			actualID, err := s.registerUser(tt.args.ctx, tt.args.u)
			if err != nil {
				assert.Error(err, "expected an error")
				assert.Empty(actualID, "got a ID but did not expect one")
				assert.EqualError(err, tt.expErrStr)
			} else {
				assert.NoError(err)
				assert.True(primitive.IsValidObjectID(actualID), "user id returned is not a valid mongo id")
			}
		})
	}
}

func Test_Service_LoginUser(t *testing.T) {
	type fields struct {
		repository Repository
		mailer     MailerClient
		password   PasswordManager
		validate   Validator
		token      TokenClient
	}
	type args struct {
		ctx context.Context
		u   UserLoginInfo
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		expUser   UserInfo
		expErrStr string
	}{
		{
			name: "Success",
			fields: fields{
				repository: &mockRepository{
					user: &UserInfo{
						ID:        "5ef7fdd91c19e3222b41b839",
						FirstName: "Michael",
						LastName:  "Scott",
						Email:     "michaelscott@dundermifflin.com",
						// CreatedAt: ,
						// UpdatedAt: ,
						Activated: true,
					},
				},
				mailer: &mockMailerClient{}, // LoginUser does not use any methods defined on this type
				password: &mockPasswordManager{
					hashedPassword: "",
					err:            nil,
				},
				validate: &validator{},       // This is not the mock
				token:    &mockTokenClient{}, // LoginUser method does not use any methods from the token client
			},
			args: args{
				ctx: context.Background(),
				u: UserLoginInfo{
					Email:    "michaelscott@dundermifflin.com",
					Password: "1234",
				},
			},
			expUser: UserInfo{
				ID:        "5ef7fdd91c19e3222b41b839",
				FirstName: "Michael",
				LastName:  "Scott",
				Email:     "michaelscott@dundermifflin.com",
			},
			expErrStr: "",
		},
		{
			// This test will pass regardless of the email used
			// because the mockRepository simulates an ErrUserNotFound error,
			// which the service layer will catch and return an ErrInvalidEmail error.
			// Perhaps I should rethink and refactor, or write an integration test
			name: "FailUserNotFoundForEmail",
			fields: fields{
				repository: &mockRepository{
					user: nil,
					err:  ErrUserNotFound,
				},
				mailer: &mockMailerClient{}, // LoginUser does not use any methods defined on this type
				password: &mockPasswordManager{
					hashedPassword: "",
					err:            nil,
				},
				validate: &validator{},       // This is not the mock
				token:    &mockTokenClient{}, // LoginUser method does not use any methods from the token client
			},
			args: args{
				ctx: context.Background(),
				u: UserLoginInfo{
					Email:    "michaelscott@dundermifflin.com",
					Password: "1234",
				},
			},
			expUser: UserInfo{
				ID:        "5ef7fdd91c19e3222b41b839",
				FirstName: "Michael",
				LastName:  "Scott",
				Email:     "michaelscott@dundermifflin.com",
			},
			expErrStr: "user with this email was not found",
		},
		{
			name: "FailUserNotActivated",
			fields: fields{
				repository: &mockRepository{
					user: &UserInfo{
						ID:        "5ef7fdd91c19e3222b41b839",
						FirstName: "Michael",
						LastName:  "Scott",
						Email:     "michaelscott@dundermifflin.com",
						// CreatedAt: ,
						// UpdatedAt: ,
						Activated: false,
					},
				},
				mailer: &mockMailerClient{}, // LoginUser does not use any methods defined on this type
				password: &mockPasswordManager{
					hashedPassword: "",
					err:            nil,
				},
				validate: &validator{},       // This is not the mock
				token:    &mockTokenClient{}, // LoginUser method does not use any methods from the token client
			},
			args: args{
				ctx: context.Background(),
				u: UserLoginInfo{
					Email:    "michaelscott@dundermifflin.com",
					Password: "1234",
				},
			},
			expUser: UserInfo{
				ID:        "5ef7fdd91c19e3222b41b839",
				FirstName: "Michael",
				LastName:  "Scott",
				Email:     "michaelscott@dundermifflin.com",
			},
			expErrStr: "user has not activated their account",
		},
		{
			name: "FailWrongPassword",
			fields: fields{
				repository: &mockRepository{
					user: &UserInfo{
						ID:        "5ef7fdd91c19e3222b41b839",
						FirstName: "Michael",
						LastName:  "Scott",
						Email:     "michaelscott@dundermifflin.com",
						// CreatedAt: ,
						// UpdatedAt: ,
						Activated: true,
					},
				},
				mailer: &mockMailerClient{}, // LoginUser does not use any methods defined on this type
				password: &mockPasswordManager{
					hashedPassword: "",
					err:            ErrInvalidPassword,
				},
				validate: &validator{},       // This is not the mock
				token:    &mockTokenClient{}, // LoginUser method does not use any methods from the token client
			},
			args: args{
				ctx: context.Background(),
				u: UserLoginInfo{
					Email:    "michaelscott@dundermifflin.com",
					Password: "1234",
				},
			},
			expUser: UserInfo{
				ID:        "5ef7fdd91c19e3222b41b839",
				FirstName: "Michael",
				LastName:  "Scott",
				Email:     "michaelscott@dundermifflin.com",
			},
			expErrStr: "invalid password provided",
		},
		{
			name: "FailMissingEmailField",
			fields: fields{
				repository: &mockRepository{userID: "", user: nil, err: nil},
				mailer:     &mockMailerClient{}, // LoginUser does not use any methods defined on this type
				password:   &mockPasswordManager{hashedPassword: "", err: nil},
				validate:   &validator{},       // This is not the mock
				token:      &mockTokenClient{}, // LoginUser method does not use any methods from the token client
			},
			args: args{
				ctx: context.Background(),
				u: UserLoginInfo{
					Password: "1234",
				},
			},
			expUser:   UserInfo{},
			expErrStr: "missing field is required: email",
		},
		{
			name: "FailMissingPasswordField",
			fields: fields{
				repository: &mockRepository{userID: "", user: nil, err: nil},
				mailer:     &mockMailerClient{}, // LoginUser does not use any methods defined on this type
				password:   &mockPasswordManager{hashedPassword: "", err: nil},
				validate:   &validator{},       // This is not the mock
				token:      &mockTokenClient{}, // LoginUser method does not use any methods from the token client
			},
			args: args{
				ctx: context.Background(),
				u: UserLoginInfo{
					Email: "michaelscott@dundermifflin.com",
				},
			},
			expUser:   UserInfo{},
			expErrStr: "missing field is required: password",
		},
	}

	t.Setenv("JWT_SECRET_KEY", "secrets")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			s := &service{
				repository: tt.fields.repository,
				mailer:     tt.fields.mailer,
				password:   tt.fields.password,
				validate:   tt.fields.validate,
				// token:      tt.fields.token,
			}

			actualUser, err := s.loginUser(tt.args.ctx, tt.args.u)
			if err != nil {
				assert.Error(err)
				assert.Nil(actualUser)
				assert.EqualError(err, tt.expErrStr)
			} else {
				assert.NoError(err)
				assert.Equal(tt.expUser.ID, actualUser.ID)
				assert.Equal(tt.expUser.FirstName, actualUser.FirstName)
				assert.Equal(tt.expUser.LastName, actualUser.LastName)
				assert.Equal(tt.expUser.Email, actualUser.Email)
				// assert.WithinDuration(*tt.wantUser.CreatedAt, *actualUser.CreatedAt, time.Second)
				// assert.WithinDuration(*tt.wantUser.UpdatedAt, *actualUser.UpdatedAt, time.Second)

				// Only checks that the jwt token field is not empty
				assert.NotEmpty(actualUser.Token)
			}
		})
	}
}

// The following depedencies are not used by this method:
// mailer     MailerClient
// password   PasswordManager
// validate   Validator
// token      TokenClient
func Test_Service_Authenticate(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name            string
		repository      Repository
		validate        Validator
		args            args
		generateUserJWT func(userID string) (string, error)
		expUserID       string
		expErr          string
	}{
		{
			name: "Success",
			repository: &mockRepository{
				user: &UserInfo{
					ID:        "5ef7fdd91c19e3222b41b839",
					Activated: true,
				},
			},
			generateUserJWT: func(userID string) (string, error) {
				return generateUserJWT(userID)
			},
			args: args{
				userID: "5ef7fdd91c19e3222b41b839",
				ctx:    context.Background(),
			},
			expUserID: "5ef7fdd91c19e3222b41b839",
		},
		{
			name: "FailUserNotActivated",
			repository: &mockRepository{
				user: &UserInfo{
					ID:        "5ef7fdd91c19e3222b41b839",
					Activated: false,
				},
			},
			generateUserJWT: func(userID string) (string, error) {
				return generateUserJWT(userID)
			},
			args: args{
				userID: "5ef7fdd91c19e3222b41b839",
				ctx:    context.Background(),
			},
			expUserID: "5ef7fdd91c19e3222b41b839",
			expErr:    "user has not activated their account",
		},
		{
			name: "FailInvalidJWT",
			repository: &mockRepository{
				user: &UserInfo{
					ID:        "5ef7fdd91c19e3222b41b839",
					Activated: false,
				},
			},
			generateUserJWT: func(userID string) (string, error) {
				signedJWT, err := generateUserJWT(userID)
				return signedJWT + "tamperedWith", err
			},
			args: args{
				userID: "5ef7fdd91c19e3222b41b839",
				ctx:    context.Background(),
			},
			expUserID: "5ef7fdd91c19e3222b41b839",
			expErr:    "invalid jwt",
		},
		{
			name: "FailUserNotFound",
			repository: &mockRepository{
				user: nil,
				err:  ErrUserNotFound,
			},
			generateUserJWT: func(userID string) (string, error) {
				return generateUserJWT(userID)
			},
			args: args{
				userID: "5ef7fdd91c19e3222b41b839",
				ctx:    context.Background(),
			},
			expUserID: "",
			expErr:    "user not found",
		},
		{
			name: "FailNoJWTEmptyString",
			repository: &mockRepository{
				user: nil,
				err:  ErrUserNotFound,
			},
			generateUserJWT: func(userID string) (string, error) {
				return "", nil
			},
			args: args{
				userID: "5ef7fdd91c19e3222b41b839",
				ctx:    context.Background(),
			},
			expUserID: "",
			expErr:    "missing field is required: jwt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			s := &service{
				validate:   &validator{},
				repository: tt.repository,
			}

			signedJWT, err := tt.generateUserJWT(tt.args.userID)
			require.NoError(err)

			actualUserID, err := s.authenticate(tt.args.ctx, signedJWT)
			if err != nil {
				require.Error(err)
				assert.EqualError(err, tt.expErr)
			} else {
				require.NoError(err)
				assert.Equal(tt.expUserID, actualUserID)
			}
		})
	}
}
