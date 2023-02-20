package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// metadata type map[string]any
type metadata map[string]any

// Value implements the driver.Valuer interface, automatically called on a selection from the database.
func (c metadata) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan implements sql.Scanner interface, automatically called on saving to the database.
func (c metadata) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &c)
}

// Event object representation of the database table events
type Event struct {
	ID       int       `db:"id"`
	UserID   int       `db:"user_id"`
	Type     string    `db:"type"`
	Message  string    `db:"message"`
	Metadata metadata  `db:"metadata"`
	DateTime time.Time `db:"time"`
}
