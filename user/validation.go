package user

import (
	"context"
	"fmt"
)

type ValidationService struct {
	Service Service
}

func (vs *ValidationService) RegisterUser(ctx context.Context, u *UserRegistrationInfo) (string, error) {

	if u.Email == "" {
		return "", fmt.Errorf("%w: email", ErrMissingField)
	}

	if u.Password == "" {
		return "", fmt.Errorf("%w: password", ErrMissingField)
	}

	return vs.Service.RegisterUser(ctx, u)
}

func (vs *ValidationService) LoginUser(ctx context.Context, u *UserLoginInfo) (*UserInfo, error) {
	if u.Email == "" {
		return nil, fmt.Errorf("%w: email", ErrMissingField)
	}

	if u.Password == "" {
		return nil, fmt.Errorf("%w: password", ErrMissingField)
	}

	return vs.Service.LoginUser(ctx, u)
}
