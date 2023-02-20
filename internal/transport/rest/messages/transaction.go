package messages

import "time"

// Transaction object representation of response connection with transactions functionality.
type Transaction struct {
	ID          int       `json:"id"`
	FromAccount int       `json:"from_account"`
	ToAccount   int       `json:"to_account"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}
