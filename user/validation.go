package user

import (
	"fmt"
)

type Validator struct{}

func (v *Validator) ValidateRegistration(u *UserRegistrationInfo) error {

	if u.Email == "" {
		return fmt.Errorf("%w: email", ErrMissingField)
	}

	if u.Password == "" {
		return fmt.Errorf("%w: password", ErrMissingField)
	}

	return nil
}

func (v *Validator) ValidateLogin(u *UserLoginInfo) error {
	if u.Email == "" {
		return fmt.Errorf("%w: email", ErrMissingField)
	}

	if u.Password == "" {
		return fmt.Errorf("%w: password", ErrMissingField)
	}

	return nil
}
