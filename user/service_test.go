package user_test

import (
	"context"
	"testing"

	"github.com/ricxi/flat-list/user"
)

func TestRegisterUser(t *testing.T) {
	testCases := []struct {
		name              string
		mockRepo          Repository
		uRegistrationInfo user.UserRegistrationInfo
		expectedUserID    string
		expectedError     error
	}{
		{
			// hard code the id for determinism?
			name: "RegisterSuccess",
			mockRepo: Repository{
				UserID: "5ef7fdd91c19e3222b41b839",
				Err:    nil,
			},
			uRegistrationInfo: user.UserRegistrationInfo{
				FirstName: "Michael",
				LastName:  "Scott",
				Email:     "michaelscott@dundermifflin.com",
				Password:  "1234",
			},
			expectedUserID: "5ef7fdd91c19e3222b41b839",
			expectedError:  nil,
		},
		{
			// unit tests might not be the best for this?
			name: "RegisterFailDuplicateUser",
			mockRepo: Repository{
				UserID: "5ef7fdd91c19e3222b41b839",
				Err:    user.ErrDuplicateUser,
			},
			uRegistrationInfo: user.UserRegistrationInfo{
				FirstName: "Michael",
				LastName:  "Scott",
				Email:     "michaelscott@dundermifflin.com",
				Password:  "1234",
			},
			expectedUserID: "",
			expectedError:  user.ErrDuplicateUser,
		},
	}

	for _, tc := range testCases {
		service := user.NewService(&tc.mockRepo, &mockPasswordService{})
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

func TestLoginUser(t *testing.T) {
	testCases := []struct {
		name             string
		expectedUserInfo user.UserInfo
		uLoginInfo       user.UserLoginInfo
		expectedErr      error
	}{
		{
			name: "LoginUserSuccess",
			uLoginInfo: user.UserLoginInfo{
				Email:    "michaelscott@dundermifflin.com",
				Password: "1234",
			},
			expectedUserInfo: user.UserInfo{
				ID:        "5ef7fdd91c19e3222b41b839",
				FirstName: "Michael",
				LastName:  "Scott",
				Email:     "michaelscott@dundermifflin.com",
				Token:     "",
			},
			expectedErr: nil,
		},
	}

	// setup environment variables
	t.Setenv("JWT_SECRET_KEY", "testsecrets")

	for _, tc := range testCases {
		mockRepo := Repository{
			user: &tc.expectedUserInfo,
			Err:  nil,
		}
		service := user.NewService(&mockRepo, &mockPasswordService{err: nil})

		uInfo, err := service.LoginUser(context.Background(), &tc.uLoginInfo)
		if err != nil {
			t.Errorf("did not expect an error, but got one: %v", err)
		}

		if uInfo.ID != tc.expectedUserInfo.ID {
			t.Errorf("expected %s, but got %s", uInfo.ID, tc.expectedUserInfo.ID)
		}
	}
}
