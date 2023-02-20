package service

import (
	"context"

	"github.com/lukinairina90/banking_backend/internal/domain"
)

//go:generate mockgen -destination=./mock_test.go -package=service -source=./contracts.go

// RolesRepository contract for roles repository.
type RolesRepository interface {
	GetByName(ctx context.Context, name string) (domain.Role, error)
}

// SessionRepository contract for refresh session repository.
type SessionRepository interface {
	Create(ctx context.Context, token domain.RefreshSession) error
	Get(ctx context.Context, token string) (domain.RefreshSession, error)
}

// RandomGenerator contract for generator.
type RandomGenerator interface {
	GenerateRandomIban() string
	GenerateRandomCardNumber() string
	GenerateRandomCvv() string
}

// AccountRepository contract for account repository.
type AccountRepository interface {
	GetUserIDByAccountID(ctx context.Context, accountID int) (int, error)
	GetAccountIDByIban(ctx context.Context, iban string) (int, error)
	GetAccountCurrencyIDByID(ctx context.Context, accountID int) (int, error)
	GetAccountCurrencyIDByIban(ctx context.Context, iban string) (int, error)
	GetAccountAmount(ctx context.Context, accountID, userID int) (float64, error)
	ExistsAccount(ctx context.Context, accountID int) (bool, error)
	Create(ctx context.Context, userID, currencyID int, iban string) (domain.Account, error)
	GetAccountsList(ctx context.Context, userID int, paginator domain.Paginator, ordering domain.Orderings) ([]domain.Account, error)
	GetAccount(ctx context.Context, accountID, userID int) (domain.Account, error)
	DeleteAccount(ctx context.Context, accountID int) error
	DepositAccount(ctx context.Context, accountID int, amount float64) error
	TransferAccount(ctx context.Context, fromAccountID, userID int, amount float64, toAccountIban string) error
	BlockAccount(ctx context.Context, accountID, userID int) error
	UnblockAccount(ctx context.Context, accountID, userID int) error
}

// TransactionRepository contract for transaction repository.
type TransactionRepository interface {
	CreateTransaction(ctx context.Context, fromAccountID, toAccountID int, amount float64) (domain.Transaction, error)
	SetTransactionStatusToSent(ctx context.Context, transactionID int) error
	GetTransactionList(ctx context.Context, accountID int, ordering domain.Orderings, paginator domain.Paginator) ([]domain.Transaction, error)
}

// CardRepository contract for card repository.
type CardRepository interface {
	CreateCard(ctx context.Context, accountID int, cardNumber string, cardholderName string, cvvCode string) (domain.Card, error)
	GetCardListUser(ctx context.Context, userID int) ([]domain.Card, error)
	GetCardListByAccount(ctx context.Context, userID, accountID int) ([]domain.Card, error)
	GetCard(ctx context.Context, id, accountID int) (domain.Card, error)
}

// EventRepository contract for event repository.
type EventRepository interface {
	CreateEvent(ctx context.Context, event domain.Event) error
	GetEventsList(ctx context.Context, userID int) ([]domain.Event, error)
}

// PasswordHasher contract for hash.
type PasswordHasher interface {
	Hash(password string) (string, error)
}

// UsersRepository contract for user repository.
type UsersRepository interface {
	GetUserNameAndSurnameByID(ctx context.Context, userID int) (string, error)
	Exists(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	GetByID(ctx context.Context, id int) (domain.User, error)
	BlockUser(ctx context.Context, userID int) error
	UnblockUser(ctx context.Context, userID int) error
	CheckBlockUser(ctx context.Context, userID int) (bool, error)
}
