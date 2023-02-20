package models

import (
	"github.com/lukinairina90/banking_backend/internal/domain"
)

// Role object representation of the database table roles
type Role struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

// ToDomain converts Role to domain.Role
func (r Role) ToDomain() domain.Role {
	return domain.Role{
		ID:   r.Id,
		Name: r.Name,
	}
}
