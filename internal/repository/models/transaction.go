package models

import (
	"database/sql"
	"time"

	"github.com/lukinairina90/banking_backend/internal/domain"
)

// Transaction object representation of the database table transaction
type Transaction struct {
	ID          int           `db:"id"`
	FromAccount sql.NullInt64 `db:"from_account"`
	ToAccount   int           `db:"to_account"`
	Amount      float64       `db:"amount"`
	Status      string        `db:"status"`
	DateCreated time.Time     `db:"date_created"`
	DateUpdated sql.NullTime  `db:"date_updated"`
}

// ToDomain converts Transaction to domain.Transaction
func (t Transaction) ToDomain() domain.Transaction {
	var fromAccountID int
	// if false, leave the default fromAccountID and do not enter to if
	// if true means the value in the database is not null and enter if
	if t.FromAccount.Valid {
		// if the value is not nul, we take the field Int64 and set it to a variable fromAccountID
		fromAccountID = int(t.FromAccount.Int64)
	}

	var dateUpdated time.Time
	// if false, leave the default dateUpdated and do not enter to if
	// if true means the value in the database is not null and enter if
	if t.DateUpdated.Valid {
		// if the value is not nul, we take the field Time and set it to a variable dateUpdated
		dateUpdated = t.DateUpdated.Time
	}

	return domain.Transaction{
		ID:          t.ID,
		FromAccount: fromAccountID,
		ToAccount:   t.ToAccount,
		Amount:      t.Amount,
		Status:      t.Status,
		DateCreated: t.DateCreated,
		DateUpdated: dateUpdated,
	}
}
