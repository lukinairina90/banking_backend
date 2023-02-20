package models

import (
	"time"

	"github.com/lukinairina90/banking_backend/internal/domain"
)

// Card object representation of the database table cards
type Card struct {
	Id             int       `db:"id"`
	AccountID      int       `db:"account_id"`
	CardNumber     string    `db:"card_number"`
	CardholderName string    `db:"cardholder_name"`
	ExpirationDate time.Time `db:"expiration_date"`
	CvvCode        string    `db:"cvv_code"`
}

// ToDomain converts Card to domain.Card
func (c Card) ToDomain() domain.Card {
	return domain.Card{
		Id:             c.Id,
		AccountID:      c.AccountID,
		CardNumber:     c.CardNumber,
		CardholderName: c.CardholderName,
		ExpirationDate: c.ExpirationDate,
		CvvCode:        c.CvvCode,
	}
}
