package messages

import "time"

// Event object representation of response connection with events functionality.
type Event struct {
	ID       int            `json:"id"`
	UserID   int            `json:"user_id"`
	Type     string         `json:"type"`
	Message  string         `json:"message"`
	Metadata map[string]any `json:"metadata"`
	DateTime time.Time      `json:"date_time"`
}
