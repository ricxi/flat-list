package user

import (
	"errors"
)

var ErrMissingField = errors.New("missing field is required")
var ErrUserNotFound = errors.New("user not found")
var ErrDuplicateUser = errors.New("user already exists")
var ErrUserNotActivated = errors.New("user has not activated their account")
var ErrInvalidEmail = errors.New("user with this email was not found")
var ErrInvalidPassword = errors.New("invalid password provided")

// used by helper functions in service
var ErrMissingEnvs = errors.New("service: missing environment variables")
