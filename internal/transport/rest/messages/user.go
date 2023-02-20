package messages

import (
	"errors"
	"time"

	"github.com/lukinairina90/banking_backend/internal/domain"
)

var ErrUserNotFound = errors.New("user with such credentials not found")

// User object representation of response connection with users functionality.
type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RoleId       int       `json:"role_id"`
	Blocked      bool      `json:"blocked"`
	RegisteredAt time.Time `json:"registered_at"`
}

// SignUpInput object representation of response.
type SignUpInput struct {
	Name     string `json:"name" binding:"required,gte=2"`
	Surname  string `json:"surname" binding:"required,gte=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6"`
}

// ToDomain converts SignUpInput to domain.SignUpInput
func (i SignUpInput) ToDomain() domain.SignUpInput {
	return domain.SignUpInput{
		Name:     i.Name,
		Surname:  i.Surname,
		Email:    i.Email,
		Password: i.Password,
	}
}

// SignInInput object representation of response.
type SignInInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6"`
}

// ToDomain converts SignInInput to domain.SignInInput
func (i SignInInput) ToDomain() domain.SignInInput {
	return domain.SignInInput{
		Email:    i.Email,
		Password: i.Password,
	}
}
