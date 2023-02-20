package rest

import (
	"context"

	"github.com/lukinairina90/banking_backend/internal/domain"
)

type AccountService interface {
	Create(ctx context.Context, userID, currencyID int) (domain.Account, error)
	GetAccountsList(ctx context.Context, userID int, paginator domain.Paginator, ordering domain.Orderings) ([]domain.Account, error)
	GetAccount(ctx context.Context, accountID, userID int) (domain.Account, error)
	DeleteAccount(ctx context.Context, accountID int) error
	DepositAccount(ctx context.Context, accountID int, amount float64) error
	TransferAccount(ctx context.Context, fromAccountID, userID int, amount float64, toAccountIban string) error
	BlockAccount(ctx context.Context, accountID, userID int) error
	UnblockAccount(ctx context.Context, accountID, userID int) error
}

type UserService interface {
	SignUp(ctx context.Context, inp domain.SignUpInput) error
	SignIn(ctx context.Context, inp domain.SignInInput) (string, string, error)
	ParseToken(ctx context.Context, token string) (int, int, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
	BlockUser(ctx context.Context, blockUserID, userID int) error
	UnblockUser(ctx context.Context, userID int) error
	CheckBlockUser(ctx context.Context, userID int) (bool, error)
}

type CardService interface {
	CreateCard(ctx context.Context, accountID, UserID int) (domain.Card, error)
	GetCardListUser(ctx context.Context, userID int) ([]domain.Card, error)
	GetCardListByAccount(ctx context.Context, userID, accountID int) ([]domain.Card, error)
	GetCard(ctx context.Context, cardID, accountID, userID int) (domain.Card, error)
}

type EventService interface {
	GetEventList(ctx context.Context, userID int) ([]domain.Event, error)
}

type RoleRepository interface {
	GetByID(ctx context.Context, id int) (domain.Role, error)
}

type TransactionService interface {
	GetTransactionList(ctx context.Context, accountID, userID int, ordering domain.Orderings, paginator domain.Paginator) ([]domain.Transaction, error)
}
