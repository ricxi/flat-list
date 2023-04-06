package user

import (
	"fmt"
)

type Validator interface {
	Registration(u UserRegistrationInfo) error
	Login(u UserLoginInfo) error
}

func NewValidator() *validator {
	return &validator{}
}

type validator struct{}

func (v *validator) Registration(u UserRegistrationInfo) error {
	if u.Email == "" {
		return fmt.Errorf("%w: email", ErrMissingField)
	}

	if u.Password == "" {
		return fmt.Errorf("%w: password", ErrMissingField)
	}

	return nil
}

func (v *validator) Login(u UserLoginInfo) error {
	if u.Email == "" {
		return fmt.Errorf("%w: email", ErrMissingField)
	}

	if u.Password == "" {
		return fmt.Errorf("%w: password", ErrMissingField)
	}

	return nil
}
