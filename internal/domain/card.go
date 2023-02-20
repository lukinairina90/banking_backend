package domain

import "time"

// Card business layer card definition
type Card struct {
	Id             int
	AccountID      int
	CardNumber     string
	CardholderName string
	ExpirationDate time.Time
	CvvCode        string
}

// ListCards list Cards
type ListCards []Card
