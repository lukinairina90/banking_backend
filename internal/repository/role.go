package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lukinairina90/banking_backend/internal/domain"
	"github.com/lukinairina90/banking_backend/internal/repository/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Roles repository layer struct.
type Roles struct {
	db *sqlx.DB
}

// NewRoles constructor for Roles repository layer.
func NewRoles(db *sqlx.DB) *Roles {
	return &Roles{db: db}
}

// GetByName returns the role as provided role name.
func (r Roles) GetByName(ctx context.Context, name string) (domain.Role, error) {
	fields := logrus.Fields{
		"layer":      "repository",
		"repository": "Role",
		"method":     "GetByName",
		"name":       name,
	}

	var role models.Role

	query := "SELECT * FROM roles WHERE name = $1"

	if err := r.db.GetContext(ctx, &role, query, name); err != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("execution getting name from roles query error")

		return domain.Role{}, errors.Wrap(err, "execution getting name from query error")
	}

	return role.ToDomain(), nil
}

// GetByID returns the role as provided role id.
func (r Roles) GetByID(ctx context.Context, id int) (domain.Role, error) {
	fields := logrus.Fields{
		"layer":      "repository",
		"repository": "Role",
		"method":     "GetByID",
		"id":         id,
	}

	var role models.Role

	query := "SELECT * FROM roles WHERE id = $1"

	if err := r.db.GetContext(ctx, &role, query, id); err != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("execution getting role by ID from roles query error")

		return domain.Role{}, errors.Wrap(err, "execution getting role by ID from roles query error")
	}

	return role.ToDomain(), nil
}
