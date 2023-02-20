package domain

import "time"

// Transaction business layer transaction definition
type Transaction struct {
	ID          int
	FromAccount int
	ToAccount   int
	Amount      float64
	Type        string
	Status      string
	DateCreated time.Time
	DateUpdated time.Time
}
