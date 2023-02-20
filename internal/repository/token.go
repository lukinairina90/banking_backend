package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lukinairina90/banking_backend/internal/domain"
	"github.com/lukinairina90/banking_backend/internal/repository/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Tokens repository layer struct.
type Tokens struct {
	db *sqlx.DB
}

// NewTokens constructor for Tokens repository layer.
func NewTokens(db *sqlx.DB) *Tokens {
	return &Tokens{db: db}
}

// Create creates token in the database by provides token struct.
func (r Tokens) Create(ctx context.Context, token domain.RefreshSession) error {
	fields := logrus.Fields{
		"layer":      "repository",
		"repository": "Token",
		"method":     "Create",
		"token":      token,
	}

	_, err := r.db.ExecContext(ctx, "INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)", token.UserID, token.Token, token.ExpiresAt)
	if err != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("execution inserting into refresh_tokens query error")

		return errors.Wrap(err, "execution inserting into refresh_tokens query error")
	}
	return nil
}

// Get returns the token struct as provided token string.
func (r Tokens) Get(ctx context.Context, token string) (domain.RefreshSession, error) {
	fields := logrus.Fields{
		"layer":      "repository",
		"repository": "Token",
		"method":     "Get",
		"token":      token,
	}

	var refreshSession models.RefreshSession

	err := r.db.QueryRowContext(ctx, "SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token=$1", token).Scan(&refreshSession.ID, &refreshSession.UserID, &refreshSession.Token, &refreshSession.ExpiresAt)
	if err != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("scanning result into struct error")

		return domain.RefreshSession{}, errors.Wrap(err, "scanning result into struct error")
	}

	_, err = r.db.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE user_id=$1", refreshSession.UserID)
	if err != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("execution deleting from refresh_tokens query error")

		return domain.RefreshSession{}, errors.Wrap(err, "execution deleting from refresh_tokens query error")
	}

	domainRefreshSession := refreshSession.ToDomain()

	return domainRefreshSession, err
}
