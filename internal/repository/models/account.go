package models

import "github.com/lukinairina90/banking_backend/internal/domain"

// Account object representation of the database table accounts
type Account struct {
	ID         int     `db:"id"`
	Iban       string  `db:"iban"`
	UserID     int     `db:"user_id"`
	CurrencyID int     `db:"currency_id"`
	Blocked    bool    `db:"blocked"`
	Amount     float64 `db:"amount"`
}

// ToDomain converts Account to domain.Account
func (a Account) ToDomain() domain.Account {
	return domain.Account{
		ID:         a.ID,
		Iban:       a.Iban,
		UserID:     a.UserID,
		CurrencyID: a.CurrencyID,
		Blocked:    a.Blocked,
		Amount:     a.Amount,
	}
}
