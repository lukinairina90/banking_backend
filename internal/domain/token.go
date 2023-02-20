package domain

import "time"

// RefreshSession business layer refreshSession definition
type RefreshSession struct {
	ID        int
	UserID    int
	Token     string
	ExpiresAt time.Time
}
