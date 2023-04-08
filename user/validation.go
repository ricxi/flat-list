package user

import (
	"fmt"
)

type Validator interface {
	Registration(u UserRegistrationInfo) error
	Login(u UserLoginInfo) error
	NonEmptyString(name, field string) error
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

// NonEmptyString checks if the given string field is not empty.
// If it is, then it returns an error message with the field's name.
func (v *validator) NonEmptyString(name, field string) error {
	if field == "" {
		return fmt.Errorf("%w: %s", ErrMissingField, name)
	}

	return nil
}
