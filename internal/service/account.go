package service

import (
	"context"

	"github.com/lukinairina90/banking_backend/internal/domain"
	"github.com/pkg/errors"
)

const fromAccountIDForDeposit = 0

// Account business logic layer struct.
type Account struct {
	accountRepo     AccountRepository
	transactionRepo TransactionRepository
	eventRepo       EventRepository
	ibanGenerator   RandomGenerator
}

// NewAccount constructor for Account.
func NewAccount(accountRepo AccountRepository, transactionRepo TransactionRepository, eventRepository EventRepository, ibanGenerator RandomGenerator) *Account {
	return &Account{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		eventRepo:       eventRepository,
		ibanGenerator:   ibanGenerator,
	}
}

// Create creates user account for user by provided currency ID.
func (s Account) Create(ctx context.Context, userID, currencyID int) (domain.Account, error) {
	iban := s.ibanGenerator.GenerateRandomIban()

	account, err := s.accountRepo.Create(ctx, userID, currencyID, iban)
	if err != nil {
		return domain.Account{}, errors.Wrap(err, "account creation error")
	}

	event := domain.Event{
		UserID:  userID,
		Type:    domain.AccountCreatedEvent,
		Message: "new account successfully created",
		Metadata: map[string]any{
			"account_id":  account.ID,
			"currency_id": account.CurrencyID,
			"iban":        account.Iban,
		},
	}

	if err := s.eventRepo.CreateEvent(ctx, event); err != nil {
		return domain.Account{}, errors.Wrap(err, "account created event creation error")
	}

	return account, nil
}

// GetAccountsList provides user accounts list.
func (s Account) GetAccountsList(ctx context.Context, userID int, paginator domain.Paginator, ordering domain.Orderings) ([]domain.Account, error) {
	list, err := s.accountRepo.GetAccountsList(ctx, userID, paginator, ordering)
	if err != nil {
		return nil, errors.Wrap(err, "getting account list error")
	}

	return list, nil
}

// GetAccount returns the user account as provided account ID and user ID.
func (s Account) GetAccount(ctx context.Context, accountID, userID int) (domain.Account, error) {
	account, err := s.accountRepo.GetAccount(ctx, accountID, userID)
	if err != nil {
		return domain.Account{}, errors.Wrap(err, "getting account error")
	}

	return account, nil
}

// DeleteAccount removes the user account as provided account ID and creates an event..
func (s Account) DeleteAccount(ctx context.Context, accountID int) error {
	err := s.accountRepo.DeleteAccount(ctx, accountID)
	if err != nil {
		return errors.Wrap(err, "account deleting error")
	}

	event := domain.Event{
		UserID:  accountID,
		Type:    domain.AccountDeletedEvent,
		Message: "account successfully deleted",
	}

	if err := s.eventRepo.CreateEvent(ctx, event); err != nil {
		return errors.Wrap(err, "account deleted event deletion error")
	}

	return nil
}

// DepositAccount replenishes the user's account by account id for the provided amount.
func (s Account) DepositAccount(ctx context.Context, accountID int, amount float64) error {
	exists, err := s.accountRepo.ExistsAccount(ctx, accountID)
	if err != nil {
		return errors.Wrap(err, "account existence check error")
	}

	if !exists {
		return errors.New("account does not exist error")
	}

	transaction, err := s.transactionRepo.CreateTransaction(ctx, fromAccountIDForDeposit, accountID, amount)
	if err != nil {
		return errors.Wrap(err, "transaction created error")
	}

	err = s.accountRepo.DepositAccount(ctx, accountID, amount)
	if err != nil {
		return errors.Wrap(err, "deposit account error")
	}

	if err := s.transactionRepo.SetTransactionStatusToSent(ctx, transaction.ID); err != nil {
		return errors.Wrap(err, "set transaction status to sent error")
	}

	event := domain.Event{
		Type:    domain.DepositEvent,
		Message: "deposit account successful",
		Metadata: map[string]any{
			"account_id": accountID,
			"amount":     amount,
		},
	}

	if err := s.eventRepo.CreateEvent(ctx, event); err != nil {
		return errors.Wrap(err, "account deposit event deposit error")
	}

	return nil
}

// TransferAccount transfer money from the user's account to another account with the same currency ID.
func (s Account) TransferAccount(ctx context.Context, fromAccountID, userID int, amount float64, toAccountIban string) error {
	fromCurrency, err := s.accountRepo.GetAccountCurrencyIDByID(ctx, fromAccountID)
	if err != nil {
		return errors.Wrap(err, "getting account currencyID by ID error")
	}

	toCurrency, err := s.accountRepo.GetAccountCurrencyIDByIban(ctx, toAccountIban)
	if err != nil {
		return errors.Wrap(err, "getting account currencyID by iban error")
	}

	if fromCurrency != toCurrency {
		return errors.New("currency matching check error")
	}

	enough, err := s.accountRepo.GetAccountAmount(ctx, fromAccountID, userID)
	if err != nil {
		return errors.Wrap(err, "getting account amount error")
	}

	if enough < amount {
		return errors.New("checking enough money in the account for the transfer error")
	}

	toAccountID, err := s.accountRepo.GetAccountIDByIban(ctx, toAccountIban)
	if err != nil {
		return errors.Wrap(err, "getting accountID by iban error")
	}

	transaction, err := s.transactionRepo.CreateTransaction(ctx, fromAccountID, toAccountID, amount)
	if err != nil {
		return errors.Wrap(err, "transaction created error")
	}

	if err := s.accountRepo.TransferAccount(ctx, fromAccountID, userID, amount, toAccountIban); err != nil {
		return errors.Wrap(err, "transfer to account error")
	}

	if err := s.transactionRepo.SetTransactionStatusToSent(ctx, transaction.ID); err != nil {
		return errors.Wrap(err, "set transaction status to sent error")
	}

	event := domain.Event{
		UserID:  userID,
		Type:    domain.WithdrawalEvent,
		Message: "money transfer successful",
		Metadata: map[string]any{
			"from_account_id": fromAccountID,
			"to_account_id":   toAccountID,
			"amount":          amount,
		},
	}

	if err := s.eventRepo.CreateEvent(ctx, event); err != nil {
		return errors.Wrap(err, "account transfer event withdrawal error")
	}

	return nil
}

// BlockAccount blocks the user account to the provided account ID and user ID.
func (s Account) BlockAccount(ctx context.Context, accountID, userID int) error {
	err := s.accountRepo.BlockAccount(ctx, accountID, userID)
	if err != nil {
		return errors.Wrap(err, "blocking account error")
	}

	event := domain.Event{
		UserID:  userID,
		Type:    domain.AccountBlockedEvent,
		Message: "account blocked successfully",
		Metadata: map[string]any{
			"account_id": accountID,
		},
	}

	if err = s.eventRepo.CreateEvent(ctx, event); err != nil {
		return errors.Wrap(err, "account blocking event blocking error")
	}

	return nil
}

// UnblockAccount unblocks the user account to the provided account ID and user ID.
func (s Account) UnblockAccount(ctx context.Context, accountID, userID int) error {
	err := s.accountRepo.UnblockAccount(ctx, accountID, userID)
	if err != nil {
		return errors.Wrap(err, "unblocking account error")
	}

	event := domain.Event{
		UserID:  userID,
		Type:    domain.AccountUnblockedEvent,
		Message: "account unblocked successfully",
		Metadata: map[string]any{
			"account_id": accountID,
		},
	}

	if err := s.eventRepo.CreateEvent(ctx, event); err != nil {
		return errors.Wrap(err, "account unblocking event unblocking error")
	}

	return nil
}
