// ! I kept a commented version of my old tests
// ! as a reference to see if I get better at testing
package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// I used a code generator to create the
// boilerplate for the table-driven tests
// this time, but I wrote the tests myself
func Test_service_RegisterUser(t *testing.T) {
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
				},
				mailer: &mockMailerClient{
					// Returning an error does not affect the service at all
					err: errors.New("dummy error"),
				},
				password: &mockPasswordManager{
					hashedPassword: "does not matter",
					err:            nil,
				},
				// This is not the mock
				validate: &validator{},
				token: &mockTokenClient{
					mockActivationToken: "does not matter",
					mockUserID:          "does not matter",
					err:                 nil,
				},
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
				},
				mailer: &mockMailerClient{
					// Returning an error does not affect the service at all
					err: errors.New("dummy error"),
				},
				password: &mockPasswordManager{
					hashedPassword: "does not matter",
					err:            nil,
				},
				// This is not the mock
				validate: &validator{},
				token: &mockTokenClient{
					mockActivationToken: "does not matter",
					mockUserID:          "does not matter",
					err:                 nil,
				},
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
				},
				mailer: &mockMailerClient{
					// Returning an error does not affect the service at all
					err: errors.New("dummy error"),
				},
				password: &mockPasswordManager{
					hashedPassword: "does not matter",
					err:            nil,
				},
				// This is not the mock
				validate: &validator{},
				token: &mockTokenClient{
					mockActivationToken: "does not matter",
					mockUserID:          "does not matter",
					err:                 nil,
				},
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
				},
				mailer: &mockMailerClient{
					// Returning an error does not affect the service at all
					err: errors.New("dummy error"),
				},
				password: &mockPasswordManager{
					hashedPassword: "does not matter",
					err:            nil,
				},
				// This is not the mock
				validate: &validator{},
				token: &mockTokenClient{
					mockActivationToken: "does not matter",
					mockUserID:          "does not matter",
					err:                 nil,
				},
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
					// Returning an error here does not affect the service at all
					err: errors.New("dummy error"),
				},
				password: &mockPasswordManager{
					hashedPassword: "does not matter because it will never be exposed outside of service",
					err:            bcrypt.ErrPasswordTooLong,
				},
				// This is not the mock
				validate: &validator{},
				token: &mockTokenClient{
					mockActivationToken: "does not matter",
					mockUserID:          "does not matter",
					err:                 nil,
				},
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
					// Returning an error here does not affect the service at all
					err: errors.New("dummy error"),
				},
				password: &mockPasswordManager{
					hashedPassword: "does not matter because it will never be exposed outside of service",
					err:            nil,
				},
				// This is not the mock
				validate: &validator{},
				token: &mockTokenClient{
					mockActivationToken: "does not matter",
					mockUserID:          "does not matter",
					err:                 nil,
				},
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

			actualID, err := s.RegisterUser(tt.args.ctx, tt.args.u)
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

// func TestServiceRegisterUser(t *testing.T) {
// 	testCases := []struct {
// 		name              string
// 		uRegistrationInfo UserRegistrationInfo
// 		expectedUserID    string
// 		inputRepoUserID   string
// 		inputRepoErr      error
// 		expectedError     error
// 	}{
// 		{
// 			// hard code the id for determinism?
// 			name: "Success",
// 			uRegistrationInfo: UserRegistrationInfo{
// 				FirstName: "Michael",
// 				LastName:  "Scott",
// 				Email:     "michaelscott@dundermifflin.com",
// 				Password:  "1234",
// 			},
// 			inputRepoUserID: "5ef7fdd91c19e3222b41b839",
// 			inputRepoErr:    nil,
// 			expectedUserID:  "5ef7fdd91c19e3222b41b839",
// 			expectedError:   nil,
// 		},
// 		{
// 			// unit tests might not be the best for this?
// 			name: "FailDuplicateUser",
// 			uRegistrationInfo: UserRegistrationInfo{
// 				FirstName: "Michael",
// 				LastName:  "Scott",
// 				Email:     "michaelscott@dundermifflin.com",
// 				Password:  "1234",
// 			},
// 			inputRepoUserID: "5ef7fdd91c19e3222b41b839",
// 			inputRepoErr:    ErrDuplicateUser,
// 			expectedUserID:  "",
// 			expectedError:   ErrDuplicateUser,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		mockRepo := mockRepository{
// 			userID: tc.inputRepoUserID,
// 			err:    tc.inputRepoErr,
// 		}
// 		service := NewServiceBuilder().
// 			Repository(&mockRepo).
// 			MailerClient(&mockMailerClient{}).
// 			TokenClient(&mockTokenClient{}).
// 			PasswordManager(&mockPasswordManager{}).
// 			Validator(&mockValidator{}).
// 			Build()

// 		userID, err := service.RegisterUser(context.Background(), &tc.uRegistrationInfo)
// 		t.Run(tc.name, func(t *testing.T) {
// 			if err != tc.expectedError {
// 				t.Errorf("expected %v, but got %v", err, tc.expectedError)
// 			}

// 			if userID != tc.expectedUserID {
// 				t.Errorf("expected %s, but got %s", tc.expectedUserID, userID)
// 			}
// 		})
// 	}
// }

// // Does not test JWT features yet
// func TestServiceLoginUser(t *testing.T) {
// 	testCases := []struct {
// 		name              string
// 		uLoginInfo        *UserLoginInfo
// 		inputRepoUserInfo *UserInfo
// 		inputRepoErr      error
// 		inputPasswordErr  error
// 		expectedUserInfo  *UserInfo
// 		expectedErr       error
// 	}{
// 		{
// 			name: "Success",
// 			uLoginInfo: &UserLoginInfo{
// 				Email:    "michaelscott@dundermifflin.com",
// 				Password: "1234",
// 			},
// 			inputRepoUserInfo: &UserInfo{
// 				ID:        "5ef7fdd91c19e3222b41b839",
// 				FirstName: "Michael",
// 				LastName:  "Scott",
// 				Email:     "michaelscott@dundermifflin.com",
// 				Token:     "",
// 			},
// 			inputRepoErr:     nil,
// 			inputPasswordErr: nil,
// 			expectedUserInfo: &UserInfo{
// 				ID:        "5ef7fdd91c19e3222b41b839",
// 				FirstName: "Michael",
// 				LastName:  "Scott",
// 				Email:     "michaelscott@dundermifflin.com",
// 				Token:     "",
// 			},
// 			expectedErr: nil,
// 		},
// 		{
// 			name: "FailedInvalidEmail",
// 			uLoginInfo: &UserLoginInfo{
// 				Email:    "michaelscott@dundermifflin.com",
// 				Password: "1234",
// 			},
// 			inputRepoUserInfo: &UserInfo{},
// 			inputRepoErr:      ErrUserNotFound,
// 			inputPasswordErr:  nil,
// 			expectedUserInfo:  &UserInfo{},
// 			expectedErr:       ErrInvalidEmail,
// 		},
// 		{
// 			name: "FailedWrongPassword",
// 			uLoginInfo: &UserLoginInfo{
// 				Email:    "michaelscott@dundermifflin.com",
// 				Password: "1234",
// 			},
// 			inputRepoUserInfo: &UserInfo{
// 				ID:             "5ef7fdd91c19e3222b41b839",
// 				FirstName:      "Michael",
// 				LastName:       "Scott",
// 				Email:          "michaelscott@dundermifflin.com",
// 				HashedPassword: "doesntmatterwhatiputhere",
// 				Activated:      true,
// 			},
// 			inputRepoErr:     nil,
// 			inputPasswordErr: bcrypt.ErrMismatchedHashAndPassword,
// 			expectedUserInfo: &UserInfo{},
// 			expectedErr:      ErrInvalidPassword,
// 		},
// 	}

// 	// setup environment variables
// 	t.Setenv("JWT_SECRET_KEY", "testsecrets")

// 	for _, tc := range testCases {
// 		mockRepo := mockRepository{
// 			user: tc.inputRepoUserInfo,
// 			err:  tc.inputRepoErr,
// 		}
// 		mockPasswordManager := mockPasswordManager{
// 			err: tc.inputPasswordErr,
// 		}
// 		service := NewServiceBuilder().
// 			Repository(&mockRepo).
// 			MailerClient(&mockMailerClient{err: nil}).
// 			TokenClient(&mockTokenClient{}).
// 			PasswordManager(&mockPasswordManager).
// 			Validator(&mockValidator{}).
// 			Build()

// 		t.Run(tc.name, func(t *testing.T) {
// 			uInfo, err := service.LoginUser(context.Background(), tc.uLoginInfo)
// 			if err != nil {
// 				if tc.expectedErr != nil && tc.expectedErr != err {
// 					t.Errorf("expected error %q, but got error %q", tc.expectedErr.Error(), err.Error())
// 				}
// 				// if nil != uInfo {
// 				// 	t.Errorf("did not expect user info, but got %v", uInfo)
// 				// }
// 			} else {
// 				if uInfo.ID != tc.expectedUserInfo.ID {
// 					t.Errorf("expected %s, but got %s", uInfo.ID, tc.expectedUserInfo.ID)
// 				}
// 				// if uInfo.FirstName != tc.expectedUserInfo.FirstName {
// 				// 	t.Errorf("expected %s, but got %s", uInfo.FirstName, tc.expectedUserInfo.FirstName)
// 				// }
// 				// if uInfo.LastName != tc.expectedUserInfo.LastName {
// 				// 	t.Errorf("expected %s, but got %s", uInfo.LastName, tc.expectedUserInfo.LastName)
// 				// }

// 				// if uInfo.Email != tc.expectedUserInfo.Email {
// 				// 	t.Errorf("expected %s, but got %s", uInfo.Email, tc.expectedUserInfo.Email)
// 				// }
// 			}
// 		})
// 	}
// }
