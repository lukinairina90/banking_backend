package domain

import (
	"time"
)

// User business layer user definition
type User struct {
	ID           int
	Name         string
	Surname      string
	Email        string
	Password     string
	RoleId       int
	Blocked      bool
	RegisteredAt time.Time
}

// SignUpInput business signUpInput user definition
type SignUpInput struct {
	Name     string
	Surname  string
	Email    string
	Password string
}

// SignInInput business signInInput user definition
type SignInInput struct {
	Email    string
	Password string
}
