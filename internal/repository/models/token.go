package models

import (
	"time"

	"github.com/lukinairina90/banking_backend/internal/domain"
)

// RefreshSession object representation of the database table refresh_tokens
type RefreshSession struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
}

// ToDomain converts RefreshSession to domain.RefreshSession
func (s RefreshSession) ToDomain() domain.RefreshSession {
	return domain.RefreshSession{
		ID:        s.ID,
		UserID:    s.UserID,
		Token:     s.Token,
		ExpiresAt: s.ExpiresAt,
	}
}
