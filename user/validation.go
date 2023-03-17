package user

import (
	"fmt"
)

type Validator interface {
	ValidateRegistration(u *UserRegistrationInfo) error
	ValidateLogin(u *UserLoginInfo) error
}

func NewValidator() Validator {
	return validator{}
}

type validator struct{}

func (v validator) ValidateRegistration(u *UserRegistrationInfo) error {

	if u.Email == "" {
		return fmt.Errorf("%w: email", ErrMissingField)
	}

	if u.Password == "" {
		return fmt.Errorf("%w: password", ErrMissingField)
	}

	return nil
}

func (v validator) ValidateLogin(u *UserLoginInfo) error {
	if u.Email == "" {
		return fmt.Errorf("%w: email", ErrMissingField)
	}

	if u.Password == "" {
		return fmt.Errorf("%w: password", ErrMissingField)
	}

	return nil
}
