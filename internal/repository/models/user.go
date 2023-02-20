package models

import (
	"time"

	"github.com/lukinairina90/banking_backend/internal/domain"
)

// User object representation of the database table users
type User struct {
	ID           int       `db:"id"`
	Name         string    `db:"name"`
	Surname      string    `db:"surname"`
	Email        string    `db:"email"`
	Password     string    `db:"password"`
	RoleId       int       `db:"role_id"`
	Blocked      bool      `db:"blocked"`
	RegisteredAt time.Time `db:"registered_at"`
}

// ToDomain converts User to domain.User
func (u User) ToDomain() domain.User {
	return domain.User{
		ID:           u.ID,
		Name:         u.Name,
		Surname:      u.Surname,
		Email:        u.Email,
		Password:     u.Password,
		RoleId:       u.RoleId,
		Blocked:      u.Blocked,
		RegisteredAt: u.RegisteredAt,
	}
}

type NameSurname struct {
	Name    string `db:"name"`
	Surname string `db:"surname"`
}
