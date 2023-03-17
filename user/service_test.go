// ! nil pointers might occur because
// ! there are a few depedencies that have
// ! not been mocked yet
package user_test

import (
	"context"
	"testing"

	"github.com/ricxi/flat-list/user"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterUser(t *testing.T) {
	testCases := []struct {
		name              string
		uRegistrationInfo user.UserRegistrationInfo
		expectedUserID    string
		inputRepoUserID   string
		inputRepoErr      error
		expectedError     error
	}{
		{
			// hard code the id for determinism?
			name: "RegisterSuccess",
			uRegistrationInfo: user.UserRegistrationInfo{
				FirstName: "Michael",
				LastName:  "Scott",
				Email:     "michaelscott@dundermifflin.com",
				Password:  "1234",
			},
			inputRepoUserID: "5ef7fdd91c19e3222b41b839",
			inputRepoErr:    nil,
			expectedUserID:  "5ef7fdd91c19e3222b41b839",
			expectedError:   nil,
		},
		{
			// unit tests might not be the best for this?
			name: "RegisterFailDuplicateUser",
			uRegistrationInfo: user.UserRegistrationInfo{
				FirstName: "Michael",
				LastName:  "Scott",
				Email:     "michaelscott@dundermifflin.com",
				Password:  "1234",
			},
			inputRepoUserID: "5ef7fdd91c19e3222b41b839",
			inputRepoErr:    user.ErrDuplicateUser,
			expectedUserID:  "",
			expectedError:   user.ErrDuplicateUser,
		},
	}

	for _, tc := range testCases {
		mockRepo := mockRepository{
			userID: tc.inputRepoUserID,
			err:    tc.inputRepoErr,
		}
		service := user.NewService(&mockRepo, &mockPasswordManager{}, &mockMailerClient{})

		userID, err := service.RegisterUser(context.Background(), &tc.uRegistrationInfo)
		t.Run(tc.name, func(t *testing.T) {
			if err != tc.expectedError {
				t.Errorf("expected %v, but got %v", err, tc.expectedError)
			}

			if userID != tc.expectedUserID {
				t.Errorf("expected %s, but got %s", tc.expectedUserID, userID)
			}
		})
	}
}

// Does not test JWT features yet
func TestLoginUser(t *testing.T) {
	testCases := []struct {
		name              string
		uLoginInfo        *user.UserLoginInfo
		inputRepoUserInfo *user.UserInfo
		inputRepoErr      error
		inputPasswordErr  error
		expectedUserInfo  *user.UserInfo
		expectedErr       error
	}{
		{
			name: "LoginUserSuccess",
			uLoginInfo: &user.UserLoginInfo{
				Email:    "michaelscott@dundermifflin.com",
				Password: "1234",
			},
			inputRepoUserInfo: &user.UserInfo{
				ID:        "5ef7fdd91c19e3222b41b839",
				FirstName: "Michael",
				LastName:  "Scott",
				Email:     "michaelscott@dundermifflin.com",
				Token:     "",
			},
			inputRepoErr:     nil,
			inputPasswordErr: nil,
			expectedUserInfo: &user.UserInfo{
				ID:        "5ef7fdd91c19e3222b41b839",
				FirstName: "Michael",
				LastName:  "Scott",
				Email:     "michaelscott@dundermifflin.com",
				Token:     "",
			},
			expectedErr: nil,
		},
		{
			name: "LoginUserFailedInvalidEmail",
			uLoginInfo: &user.UserLoginInfo{
				Email:    "michaelscott@dundermifflin.com",
				Password: "1234",
			},
			inputRepoUserInfo: nil,
			inputRepoErr:      user.ErrUserNotFound,
			inputPasswordErr:  nil,
			expectedUserInfo:  nil,
			expectedErr:       user.ErrInvalidEmail,
		},
		{
			name: "LoginUserFailedWrongPassword",
			uLoginInfo: &user.UserLoginInfo{
				Email:    "michaelscott@dundermifflin.com",
				Password: "1234",
			},
			inputRepoUserInfo: &user.UserInfo{
				ID:             "5ef7fdd91c19e3222b41b839",
				FirstName:      "Michael",
				LastName:       "Scott",
				Email:          "michaelscott@dundermifflin.com",
				HashedPassword: "doesntmatterwhatiputhere",
			},
			inputRepoErr:     nil,
			inputPasswordErr: bcrypt.ErrMismatchedHashAndPassword,
			expectedUserInfo: nil,
			expectedErr:      user.ErrInvalidPassword,
		},
	}

	// setup environment variables
	t.Setenv("JWT_SECRET_KEY", "testsecrets")

	for _, tc := range testCases {
		mockRepo := mockRepository{
			user: tc.inputRepoUserInfo,
			err:  tc.inputRepoErr,
		}
		service := user.NewService(&mockRepo, &mockPasswordManager{err: tc.inputPasswordErr}, &mockMailerClient{})

		uInfo, err := service.LoginUser(context.Background(), tc.uLoginInfo)
		t.Run(tc.name, func(t *testing.T) {
			if nil == err {
				if uInfo.ID != tc.expectedUserInfo.ID {
					t.Errorf("expected %s, but got %s", uInfo.ID, tc.expectedUserInfo.ID)
				}
				if uInfo.FirstName != tc.expectedUserInfo.FirstName {
					t.Errorf("expected %s, but got %s", uInfo.FirstName, tc.expectedUserInfo.FirstName)
				}
				if uInfo.LastName != tc.expectedUserInfo.LastName {
					t.Errorf("expected %s, but got %s", uInfo.LastName, tc.expectedUserInfo.LastName)
				}

				if uInfo.Email != tc.expectedUserInfo.Email {
					t.Errorf("expected %s, but got %s", uInfo.Email, tc.expectedUserInfo.Email)
				}
			} else {
				if tc.expectedErr != err {
					t.Errorf("expected error %q, but got error %q", tc.expectedErr.Error(), err.Error())
				}
				if nil != uInfo {
					t.Errorf("did not expect user info, but got %v", uInfo)
				}
			}
		})
	}
}
