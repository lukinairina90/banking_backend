package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lukinairina90/banking_backend/internal/domain"
	"github.com/lukinairina90/banking_backend/internal/repository/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	transactionPreparedStatus = "PREPARED"
	transactionSentStatus     = "SENT"

	ingoingTransactionType  = "ingoing"
	outgoingTransactionType = "outgoing"
)

// Transactions repository layer struct.
type Transactions struct {
	db *sqlx.DB
}

// NewTransactions constructor for Transactions repository layer.
func NewTransactions(db *sqlx.DB) *Transactions {
	return &Transactions{db: db}
}

// CreateTransaction creates a transaction in the database by provides from account ID, to account ID and amount.
func (r Transactions) CreateTransaction(ctx context.Context, fromAccountID, toAccountID int, amount float64) (domain.Transaction, error) {
	fields := logrus.Fields{
		"layer":           "repository",
		"repository":      "Transaction",
		"method":          "CreateTransaction",
		"from_account_id": fromAccountID,
		"to_account_id":   toAccountID,
		"amount":          amount,
	}

	nullableFromAccountID := sql.NullInt64{}
	if fromAccountID != 0 {
		nullableFromAccountID.Int64 = int64(fromAccountID)
		nullableFromAccountID.Valid = true
	}

	query := "INSERT INTO transactions (from_account, to_account, amount, status, date_created) VALUES ($1, $2, $3, $4, NOW()) RETURNING *"

	row := r.db.QueryRowContext(ctx, query, nullableFromAccountID, toAccountID, amount, transactionPreparedStatus)
	if row.Err() != nil {
		logrus.WithError(row.Err()).
			WithFields(fields).
			Error("execution inserting into transactions query error")

		return domain.Transaction{}, errors.Wrap(row.Err(), "execution inserting into transactions query error")
	}

	transaction := models.Transaction{}
	if err := row.Scan(
		&transaction.ID,
		&transaction.FromAccount,
		&transaction.ToAccount,
		&transaction.Amount,
		&transaction.Status,
		&transaction.DateCreated,
		&transaction.DateUpdated,
	); err != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("scanning row into struct error")

		return domain.Transaction{}, errors.Wrap(err, "scanning row into struct error")
	}

	domainTransaction := transaction.ToDomain()

	return domainTransaction, nil
}

// SetTransactionStatusToSent set transaction status to sent
func (r Transactions) SetTransactionStatusToSent(ctx context.Context, transactionID int) error {
	fields := logrus.Fields{
		"layer":          "repository",
		"repository":     "Transaction",
		"method":         "CreateTransaction",
		"transaction_id": transactionID,
	}

	query := "UPDATE transactions SET status=$1, date_updated=NOW() WHERE id = $2"

	if _, err := r.db.ExecContext(ctx, query, transactionSentStatus, transactionID); err != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("execution updating status into transactions query error")

		return errors.Wrap(err, "execution updating status into transactions query error")
	}

	return nil
}

// GetTransactionList provides transactions list.
func (r Transactions) GetTransactionList(ctx context.Context, accountID int, ordering domain.Orderings, paginator domain.Paginator) ([]domain.Transaction, error) {
	fields := logrus.Fields{
		"layer":      "repository",
		"repository": "Transaction",
		"method":     "CreateTransaction",
		"account_id": accountID,
	}

	listTransaction := make([]models.Transaction, 0, paginator.PerPage)

	//query := "SELECT * FROM transactions WHERE from_account = $1 OR to_account = $2"
	qb := squirrel.Select("*").
		From("transactions").
		Where("from_account = ? OR to_account = ?", accountID, accountID)

	if ordering != nil {
		var parts []string
		for field, direction := range ordering {
			parts = append(parts, fmt.Sprintf("%s %s", field, strings.ToUpper(direction)))
		}

		qb = qb.OrderBy(parts...)
	} else {
		qb = qb.OrderBy("id ASC")
	}

	qb = qb.Limit(uint64(paginator.PerPage)).
		Offset(uint64((paginator.Page - 1) * paginator.PerPage)).
		PlaceholderFormat(squirrel.Dollar)

	query, params, err := qb.ToSql()
	if err != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("building a ToSql query into a SQL string error")

		return nil, errors.Wrap(err, "building a ToSql query into a SQL string error")
	}

	rows, err := r.db.QueryxContext(ctx, query, params...)
	if err != nil && rows.Err() != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("execution getting transaction list from transaction query error")

		return nil, errors.Wrap(err, "execution getting transaction list from transaction query error")
	}

	for rows.Next() {
		var transaction models.Transaction
		if err := rows.StructScan(&transaction); err != nil {
			logrus.WithError(err).
				WithFields(fields).
				Error("scanning rows into struct error")

			return nil, errors.Wrap(err, "scanning rows into struct error")
		}
		listTransaction = append(listTransaction, transaction)
	}

	domainListTransaction := make([]domain.Transaction, 0)
	for _, transaction := range listTransaction {
		domainTransaction := transaction.ToDomain()

		domainTransaction.Type = outgoingTransactionType
		if domainTransaction.ToAccount == accountID {
			domainTransaction.Type = ingoingTransactionType
		}

		domainListTransaction = append(domainListTransaction, domainTransaction)
	}

	return domainListTransaction, nil
}
