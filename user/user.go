package user

import (
	"time"
)

// UserInfo is sent out as a response
type UserInfo struct {
	ID             string     `json:"id"`
	FirstName      string     `json:"firstName"`
	LastName       string     `json:"lastName"`
	Email          string     `json:"email"`
	Password       string     `json:"-"`
	HashedPassword string     `json:"j"`
	Activated      bool       `json:"-"`
	CreatedAt      *time.Time `json:"-"`
	UpdatedAt      *time.Time `json:"-"`
	Token          string     `json:"token"`
}

// UserRegistrationInfo stores request
// data for registering a new user
type UserRegistrationInfo struct {
	FirstName      string     `json:"firstName"`
	LastName       string     `json:"lastName"`
	Email          string     `json:"email"`
	Password       string     `json:"password"`
	HashedPassword string     `json:"-"`
	Activated      bool       `json:"activated"`
	CreatedAt      *time.Time `json:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt"`
}

// UserLoginInfo stores request data
// for logging in a user
type UserLoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
