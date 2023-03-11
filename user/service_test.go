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
		service := user.NewService(&tc.mockRepo)
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
