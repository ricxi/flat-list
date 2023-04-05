// ! I kept a commented version of my old tests
// ! as a reference to see if I get better at testing
package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// I tried a out a code generator to make these
// tests this time
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
		name    string
		fields  fields
		args    args
		wantID  string
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				repository: &mockRepository{},
				mailer:     &mockMailerClient{},
				password:   &mockPasswordManager{},
				validate:   &mockValidator{},
				token:      &mockTokenClient{},
			},
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			// require := require.New(t)

			s := &service{
				repository: tt.fields.repository,
				mailer:     tt.fields.mailer,
				password:   tt.fields.password,
				validate:   tt.fields.validate,
				token:      tt.fields.token,
			}

			actualID, err := s.RegisterUser(tt.args.ctx, tt.args.u)
			if assert.NoError(err) {
				assert.True(primitive.IsValidObjectID(actualID))
			} else {
				fmt.Println("no error cases yet")
			}
			// if (err != nil) != tt.wantErr {
			// 	t.Errorf("service.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
			// 	return
			// }
			// if got != tt.want {
			// 	t.Errorf("service.RegisterUser() = %v, want %v", got, tt.want)
			// }
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
