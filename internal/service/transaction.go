package service

import (
	"context"

	"github.com/lukinairina90/banking_backend/internal/domain"
	"github.com/pkg/errors"
)

// Transaction business logic layer struct.
type Transaction struct {
	TransactionRepo TransactionRepository
	AccountRepo     AccountRepository
}

// NewTransaction constructor for transaction.
func NewTransaction(transactionRepo TransactionRepository, accountRepo AccountRepository) *Transaction {
	return &Transaction{
		TransactionRepo: transactionRepo,
		AccountRepo:     accountRepo,
	}
}

// GetTransactionList returns a list of transaction for a specific account.
func (s Transaction) GetTransactionList(ctx context.Context, accountID, userID int, ordering domain.Orderings, paginator domain.Paginator) ([]domain.Transaction, error) {
	checkUserID, err := s.AccountRepo.GetUserIDByAccountID(ctx, accountID)
	if err != nil {
		return nil, errors.Wrap(err, "getting userID by accountID error")
	}

	if checkUserID != userID {
		return nil, errors.New("userID matching check error")
	}

	listTransaction, err := s.TransactionRepo.GetTransactionList(ctx, accountID, ordering, paginator)
	if err != nil {
		return nil, errors.Wrap(err, "getting transaction list error")
	}

	return listTransaction, nil
}
