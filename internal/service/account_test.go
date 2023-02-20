package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lukinairina90/banking_backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestAccount_Create(t *testing.T) {
	controller := gomock.NewController(t)
	accountRepositoryMock := NewMockAccountRepository(controller)
	eventRepositoryMock := NewMockEventRepository(controller)
	ibanGeneratorMock := NewMockRandomGenerator(controller)

	testIban := "UA031234560000039096125330468"
	testError := errors.New("test error")

	ctx := context.Background()
	userID := 1
	currencyID := 1

	account := domain.Account{
		ID:         1,
		Iban:       testIban,
		UserID:     userID,
		CurrencyID: currencyID,
		Blocked:    false,
		Amount:     0,
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

	type fields struct {
		accountRepo   AccountRepository
		eventRepo     EventRepository
		ibanGenerator RandomGenerator
	}
	type args struct {
		ctx        context.Context
		userID     int
		currencyID int
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		configureMock func()
		want          domain.Account
		wantErr       bool
	}{
		{
			name: "account_repository_error",
			fields: fields{
				accountRepo:   accountRepositoryMock,
				eventRepo:     eventRepositoryMock,
				ibanGenerator: ibanGeneratorMock,
			},
			args: args{
				ctx:        ctx,
				userID:     userID,
				currencyID: currencyID,
			},
			configureMock: func() {
				ibanGeneratorMock.EXPECT().GenerateRandomIban().Return(testIban)
				accountRepositoryMock.EXPECT().Create(gomock.Eq(ctx), gomock.Eq(userID), gomock.Eq(currencyID), gomock.Eq(testIban)).Return(domain.Account{}, testError)
			},
			want:    domain.Account{},
			wantErr: true,
		},
		{
			name: "event_repository_error",
			fields: fields{
				accountRepo:   accountRepositoryMock,
				eventRepo:     eventRepositoryMock,
				ibanGenerator: ibanGeneratorMock,
			},
			args: args{
				ctx:        ctx,
				userID:     userID,
				currencyID: currencyID,
			},
			configureMock: func() {
				ibanGeneratorMock.EXPECT().GenerateRandomIban().Return(testIban)
				accountRepositoryMock.EXPECT().Create(gomock.Eq(ctx), gomock.Eq(userID), gomock.Eq(currencyID), gomock.Eq(testIban)).Return(account, nil)
				eventRepositoryMock.EXPECT().CreateEvent(gomock.Eq(ctx), gomock.Eq(event)).Return(testError)
			},
			want:    domain.Account{},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				accountRepo:   accountRepositoryMock,
				eventRepo:     eventRepositoryMock,
				ibanGenerator: ibanGeneratorMock,
			},
			args: args{
				ctx:        ctx,
				userID:     userID,
				currencyID: currencyID,
			},
			configureMock: func() {
				ibanGeneratorMock.EXPECT().GenerateRandomIban().Return(testIban)
				accountRepositoryMock.EXPECT().Create(gomock.Eq(ctx), gomock.Eq(userID), gomock.Eq(currencyID), gomock.Eq(testIban)).Return(account, nil)
				eventRepositoryMock.EXPECT().CreateEvent(gomock.Eq(ctx), gomock.Eq(event)).Return(nil)
			},
			want:    account,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configureMock()

			a := NewAccount(tt.fields.accountRepo, nil, tt.fields.eventRepo, tt.fields.ibanGenerator)
			got, err := a.Create(tt.args.ctx, tt.args.userID, tt.args.currencyID)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAccount_GetAccountsList(t *testing.T) {
	controller := gomock.NewController(t)
	accountRepositoryMock := NewMockAccountRepository(controller)

	testIban := "UA031234560000039096125330468"
	testError := errors.New("test error")

	ctx := context.Background()
	userID := 1
	currencyID := 1

	account1 := domain.Account{
		ID:         1,
		Iban:       testIban,
		UserID:     userID,
		CurrencyID: currencyID,
		Blocked:    false,
		Amount:     0,
	}

	account2 := domain.Account{
		ID:         2,
		Iban:       testIban,
		UserID:     userID,
		CurrencyID: currencyID,
		Blocked:    false,
		Amount:     0,
	}

	accountList := []domain.Account{account1, account2}

	paginator := domain.Paginator{
		Page:    1,
		PerPage: 5,
	}

	orderiings := domain.Orderings{
		"id": "desc",
	}

	type fields struct {
		accountRepo     AccountRepository
		transactionRepo TransactionRepository
		eventRepo       EventRepository
		ibanGenerator   RandomGenerator
	}
	type args struct {
		ctx       context.Context
		userID    int
		paginator domain.Paginator
		ordering  domain.Orderings
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		configureMock func()
		want          []domain.Account
		wantErr       bool
	}{
		{
			name: "account_repository_error",
			fields: fields{
				accountRepo: accountRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				userID:    userID,
				paginator: paginator,
				ordering:  orderiings,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetAccountsList(gomock.Eq(ctx), gomock.Eq(userID), gomock.Eq(paginator), gomock.Eq(orderiings)).Return(nil, testError)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				accountRepo: accountRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				userID:    userID,
				paginator: paginator,
				ordering:  orderiings,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetAccountsList(gomock.Eq(ctx), gomock.Eq(userID), gomock.Eq(paginator), gomock.Eq(orderiings)).Return(accountList, nil)
			},
			want:    accountList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configureMock()

			a := NewAccount(tt.fields.accountRepo, nil, nil, nil)
			got, err := a.GetAccountsList(tt.args.ctx, tt.args.userID, tt.args.paginator, tt.args.ordering)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAccount_TransferAccount(t *testing.T) {
	controller := gomock.NewController(t)
	accountRepositoryMock := NewMockAccountRepository(controller)
	eventRepositoryMock := NewMockEventRepository(controller)
	transactionRepositoryMock := NewMockTransactionRepository(controller)

	ctx := context.Background()
	fromAccountID := 1
	toAccountID := 2
	userID := 1
	amount := float64(60)
	toAccountIban := "UA031234560000039096125330468"

	testError := errors.New("test error")
	fromAccountCurrencyID := 1
	toAccountCurrencyID := 1
	transaction := domain.Transaction{
		ID:          1,
		FromAccount: fromAccountID,
		ToAccount:   toAccountID,
		Amount:      amount,
		Status:      "PREPARED",
		DateCreated: time.Now(),
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

	type fields struct {
		accountRepo     AccountRepository
		transactionRepo TransactionRepository
		eventRepo       EventRepository
	}
	type args struct {
		ctx           context.Context
		fromAccountID int
		userID        int
		amount        float64
		toAccountIban string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		configureMock func()
		wantErr       bool
	}{
		{
			name: "getting_account_currency_by_id_error",
			fields: fields{
				accountRepo:     accountRepositoryMock,
				transactionRepo: transactionRepositoryMock,
				eventRepo:       eventRepositoryMock,
			},
			args: args{
				ctx:           ctx,
				fromAccountID: fromAccountID,
				userID:        userID,
				amount:        amount,
				toAccountIban: toAccountIban,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByID(gomock.Eq(ctx), gomock.Eq(fromAccountID)).Return(0, testError)
			},
			wantErr: true,
		},
		{
			name: "getting_account_currency_by_iban_error",
			fields: fields{
				accountRepo:     accountRepositoryMock,
				transactionRepo: transactionRepositoryMock,
				eventRepo:       eventRepositoryMock,
			},
			args: args{
				ctx:           ctx,
				fromAccountID: fromAccountID,
				userID:        userID,
				amount:        amount,
				toAccountIban: toAccountIban,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByID(gomock.Eq(ctx), gomock.Eq(fromAccountID)).Return(fromAccountCurrencyID, nil)
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByIban(gomock.Eq(ctx), gomock.Eq(toAccountIban)).Return(0, testError)
			},
			wantErr: true,
		},
		{
			name: "account_currencies_mismatching",
			fields: fields{
				accountRepo:     accountRepositoryMock,
				transactionRepo: transactionRepositoryMock,
				eventRepo:       eventRepositoryMock,
			},
			args: args{
				ctx:           ctx,
				fromAccountID: fromAccountID,
				userID:        userID,
				amount:        amount,
				toAccountIban: toAccountIban,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByID(gomock.Eq(ctx), gomock.Eq(fromAccountID)).Return(fromAccountCurrencyID, nil)
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByIban(gomock.Eq(ctx), gomock.Eq(toAccountIban)).Return(2, nil)
			},
			wantErr: true,
		},
		{
			name: "getting_from_account_amount_error",
			fields: fields{
				accountRepo:     accountRepositoryMock,
				transactionRepo: transactionRepositoryMock,
				eventRepo:       eventRepositoryMock,
			},
			args: args{
				ctx:           ctx,
				fromAccountID: fromAccountID,
				userID:        userID,
				amount:        amount,
				toAccountIban: toAccountIban,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByID(gomock.Eq(ctx), gomock.Eq(fromAccountID)).Return(fromAccountCurrencyID, nil)
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByIban(gomock.Eq(ctx), gomock.Eq(toAccountIban)).Return(toAccountCurrencyID, nil)
				accountRepositoryMock.EXPECT().GetAccountAmount(gomock.Eq(ctx), gomock.Eq(fromAccountCurrencyID), gomock.Eq(userID)).Return(float64(0), testError)
			},
			wantErr: true,
		},
		{
			name: "not_enough_money_error",
			fields: fields{
				accountRepo:     accountRepositoryMock,
				transactionRepo: transactionRepositoryMock,
				eventRepo:       eventRepositoryMock,
			},
			args: args{
				ctx:           ctx,
				fromAccountID: fromAccountID,
				userID:        userID,
				amount:        amount,
				toAccountIban: toAccountIban,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByID(gomock.Eq(ctx), gomock.Eq(fromAccountID)).Return(fromAccountCurrencyID, nil)
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByIban(gomock.Eq(ctx), gomock.Eq(toAccountIban)).Return(toAccountCurrencyID, nil)
				accountRepositoryMock.EXPECT().GetAccountAmount(gomock.Eq(ctx), gomock.Eq(fromAccountCurrencyID), gomock.Eq(userID)).Return(float64(50), nil)
			},
			wantErr: true,
		},
		{
			name: "getting_account_id_by_iban_error",
			fields: fields{
				accountRepo:     accountRepositoryMock,
				transactionRepo: transactionRepositoryMock,
				eventRepo:       eventRepositoryMock,
			},
			args: args{
				ctx:           ctx,
				fromAccountID: fromAccountID,
				userID:        userID,
				amount:        amount,
				toAccountIban: toAccountIban,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByID(gomock.Eq(ctx), gomock.Eq(fromAccountID)).Return(fromAccountCurrencyID, nil)
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByIban(gomock.Eq(ctx), gomock.Eq(toAccountIban)).Return(toAccountCurrencyID, nil)
				accountRepositoryMock.EXPECT().GetAccountAmount(gomock.Eq(ctx), gomock.Eq(fromAccountCurrencyID), gomock.Eq(userID)).Return(float64(100), nil)
				accountRepositoryMock.EXPECT().GetAccountIDByIban(gomock.Eq(ctx), gomock.Eq(toAccountIban)).Return(0, testError)
			},
			wantErr: true,
		},
		{
			name: "creating_transaction_error",
			fields: fields{
				accountRepo:     accountRepositoryMock,
				transactionRepo: transactionRepositoryMock,
				eventRepo:       eventRepositoryMock,
			},
			args: args{
				ctx:           ctx,
				fromAccountID: fromAccountID,
				userID:        userID,
				amount:        amount,
				toAccountIban: toAccountIban,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByID(gomock.Eq(ctx), gomock.Eq(fromAccountID)).Return(fromAccountCurrencyID, nil)
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByIban(gomock.Eq(ctx), gomock.Eq(toAccountIban)).Return(toAccountCurrencyID, nil)
				accountRepositoryMock.EXPECT().GetAccountAmount(gomock.Eq(ctx), gomock.Eq(fromAccountCurrencyID), gomock.Eq(userID)).Return(float64(100), nil)
				accountRepositoryMock.EXPECT().GetAccountIDByIban(gomock.Eq(ctx), gomock.Eq(toAccountIban)).Return(toAccountID, nil)
				transactionRepositoryMock.EXPECT().CreateTransaction(gomock.Eq(ctx), gomock.Eq(fromAccountID), gomock.Eq(toAccountID), gomock.Eq(amount)).Return(domain.Transaction{}, testError)
			},
			wantErr: true,
		},
		{
			name: "creating_transaction_error",
			fields: fields{
				accountRepo:     accountRepositoryMock,
				transactionRepo: transactionRepositoryMock,
				eventRepo:       eventRepositoryMock,
			},
			args: args{
				ctx:           ctx,
				fromAccountID: fromAccountID,
				userID:        userID,
				amount:        amount,
				toAccountIban: toAccountIban,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByID(gomock.Eq(ctx), gomock.Eq(fromAccountID)).Return(fromAccountCurrencyID, nil)
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByIban(gomock.Eq(ctx), gomock.Eq(toAccountIban)).Return(toAccountCurrencyID, nil)
				accountRepositoryMock.EXPECT().GetAccountAmount(gomock.Eq(ctx), gomock.Eq(fromAccountCurrencyID), gomock.Eq(userID)).Return(float64(100), nil)
				accountRepositoryMock.EXPECT().GetAccountIDByIban(gomock.Eq(ctx), gomock.Eq(toAccountIban)).Return(toAccountID, nil)
				transactionRepositoryMock.EXPECT().CreateTransaction(gomock.Eq(ctx), gomock.Eq(fromAccountID), gomock.Eq(toAccountID), gomock.Eq(amount)).Return(transaction, nil)
				accountRepositoryMock.EXPECT().TransferAccount(gomock.Eq(ctx), gomock.Eq(fromAccountID), gomock.Eq(userID), gomock.Eq(amount), gomock.Eq(toAccountIban)).Return(testError)
			},
			wantErr: true,
		},
		{
			name: "updating_transaction_status_error",
			fields: fields{
				accountRepo:     accountRepositoryMock,
				transactionRepo: transactionRepositoryMock,
				eventRepo:       eventRepositoryMock,
			},
			args: args{
				ctx:           ctx,
				fromAccountID: fromAccountID,
				userID:        userID,
				amount:        amount,
				toAccountIban: toAccountIban,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByID(gomock.Eq(ctx), gomock.Eq(fromAccountID)).Return(fromAccountCurrencyID, nil)
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByIban(gomock.Eq(ctx), gomock.Eq(toAccountIban)).Return(toAccountCurrencyID, nil)
				accountRepositoryMock.EXPECT().GetAccountAmount(gomock.Eq(ctx), gomock.Eq(fromAccountCurrencyID), gomock.Eq(userID)).Return(float64(100), nil)
				accountRepositoryMock.EXPECT().GetAccountIDByIban(gomock.Eq(ctx), gomock.Eq(toAccountIban)).Return(toAccountID, nil)
				transactionRepositoryMock.EXPECT().CreateTransaction(gomock.Eq(ctx), gomock.Eq(fromAccountID), gomock.Eq(toAccountID), gomock.Eq(amount)).Return(transaction, nil)
				accountRepositoryMock.EXPECT().TransferAccount(gomock.Eq(ctx), gomock.Eq(fromAccountID), gomock.Eq(userID), gomock.Eq(amount), gomock.Eq(toAccountIban)).Return(nil)
				transactionRepositoryMock.EXPECT().SetTransactionStatusToSent(gomock.Eq(ctx), gomock.Eq(transaction.ID)).Return(testError)
			},
			wantErr: true,
		},
		{
			name: "creating_event_error",
			fields: fields{
				accountRepo:     accountRepositoryMock,
				transactionRepo: transactionRepositoryMock,
				eventRepo:       eventRepositoryMock,
			},
			args: args{
				ctx:           ctx,
				fromAccountID: fromAccountID,
				userID:        userID,
				amount:        amount,
				toAccountIban: toAccountIban,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByID(gomock.Eq(ctx), gomock.Eq(fromAccountID)).Return(fromAccountCurrencyID, nil)
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByIban(gomock.Eq(ctx), gomock.Eq(toAccountIban)).Return(toAccountCurrencyID, nil)
				accountRepositoryMock.EXPECT().GetAccountAmount(gomock.Eq(ctx), gomock.Eq(fromAccountCurrencyID), gomock.Eq(userID)).Return(float64(100), nil)
				accountRepositoryMock.EXPECT().GetAccountIDByIban(gomock.Eq(ctx), gomock.Eq(toAccountIban)).Return(toAccountID, nil)
				transactionRepositoryMock.EXPECT().CreateTransaction(gomock.Eq(ctx), gomock.Eq(fromAccountID), gomock.Eq(toAccountID), gomock.Eq(amount)).Return(transaction, nil)
				accountRepositoryMock.EXPECT().TransferAccount(gomock.Eq(ctx), gomock.Eq(fromAccountID), gomock.Eq(userID), gomock.Eq(amount), gomock.Eq(toAccountIban)).Return(nil)
				transactionRepositoryMock.EXPECT().SetTransactionStatusToSent(gomock.Eq(ctx), gomock.Eq(transaction.ID)).Return(nil)
				eventRepositoryMock.EXPECT().CreateEvent(gomock.Eq(ctx), gomock.Eq(event)).Return(testError)
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				accountRepo:     accountRepositoryMock,
				transactionRepo: transactionRepositoryMock,
				eventRepo:       eventRepositoryMock,
			},
			args: args{
				ctx:           ctx,
				fromAccountID: fromAccountID,
				userID:        userID,
				amount:        amount,
				toAccountIban: toAccountIban,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByID(gomock.Eq(ctx), gomock.Eq(fromAccountID)).Return(fromAccountCurrencyID, nil)
				accountRepositoryMock.EXPECT().GetAccountCurrencyIDByIban(gomock.Eq(ctx), gomock.Eq(toAccountIban)).Return(toAccountCurrencyID, nil)
				accountRepositoryMock.EXPECT().GetAccountAmount(gomock.Eq(ctx), gomock.Eq(fromAccountCurrencyID), gomock.Eq(userID)).Return(float64(100), nil)
				accountRepositoryMock.EXPECT().GetAccountIDByIban(gomock.Eq(ctx), gomock.Eq(toAccountIban)).Return(toAccountID, nil)
				transactionRepositoryMock.EXPECT().CreateTransaction(gomock.Eq(ctx), gomock.Eq(fromAccountID), gomock.Eq(toAccountID), gomock.Eq(amount)).Return(transaction, nil)
				accountRepositoryMock.EXPECT().TransferAccount(gomock.Eq(ctx), gomock.Eq(fromAccountID), gomock.Eq(userID), gomock.Eq(amount), gomock.Eq(toAccountIban)).Return(nil)
				transactionRepositoryMock.EXPECT().SetTransactionStatusToSent(gomock.Eq(ctx), gomock.Eq(transaction.ID)).Return(nil)
				eventRepositoryMock.EXPECT().CreateEvent(gomock.Eq(ctx), gomock.Eq(event)).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configureMock()

			a := NewAccount(accountRepositoryMock, transactionRepositoryMock, eventRepositoryMock, nil)
			err := a.TransferAccount(tt.args.ctx, tt.args.fromAccountID, tt.args.userID, tt.args.amount, tt.args.toAccountIban)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestAccount_GetAccount(t *testing.T) {
	controller := gomock.NewController(t)
	accountRepositoryMock := NewMockAccountRepository(controller)

	ctx := context.Background()

	accountID := 1
	userID := 1

	account := domain.Account{
		ID:         1,
		Iban:       "UA251234560000025017606879007",
		UserID:     1,
		CurrencyID: 1,
		Blocked:    false,
		Amount:     0,
	}

	testError := errors.New("test error")

	type fields struct {
		accountRepo AccountRepository
	}

	type args struct {
		ctx       context.Context
		accountID int
		userID    int
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		configureMock func()
		want          domain.Account
		wantErr       bool
	}{
		{
			name: "account_repository_error",
			fields: fields{
				accountRepo: accountRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				userID:    userID,
				accountID: accountID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetAccount(gomock.Eq(ctx), gomock.Eq(accountID), gomock.Eq(userID)).Return(domain.Account{}, testError)
			},
			want:    domain.Account{},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				accountRepo: accountRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				userID:    userID,
				accountID: accountID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetAccount(gomock.Eq(ctx), gomock.Eq(accountID), gomock.Eq(userID)).Return(account, nil)
			},
			want:    account,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configureMock()

			a := NewAccount(tt.fields.accountRepo, nil, nil, nil)
			got, err := a.GetAccount(tt.args.ctx, tt.args.accountID, tt.args.userID)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAccount_DeleteAccount(t *testing.T) {
	controller := gomock.NewController(t)
	accountRepositoryMock := NewMockAccountRepository(controller)
	eventRepositoryMock := NewMockEventRepository(controller)

	ctx := context.Background()

	accountID := 1

	event := domain.Event{
		UserID:  accountID,
		Type:    domain.AccountDeletedEvent,
		Message: "account successfully deleted",
	}

	testError := errors.New("test error")

	type fields struct {
		accountRepo AccountRepository
		eventRepo   EventRepository
	}
	type args struct {
		ctx       context.Context
		accountID int
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		configureMock func()
		wantErr       bool
	}{
		{
			name: "account_repository_error",
			fields: fields{
				accountRepo: accountRepositoryMock,
				eventRepo:   eventRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().DeleteAccount(gomock.Eq(ctx), gomock.Eq(accountID)).Return(testError)
			},
			wantErr: true,
		},
		{
			name: "event_repository_error",
			fields: fields{
				accountRepo: accountRepositoryMock,
				eventRepo:   eventRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().DeleteAccount(gomock.Eq(ctx), gomock.Eq(accountID)).Return(nil)
				eventRepositoryMock.EXPECT().CreateEvent(gomock.Eq(ctx), gomock.Eq(event)).Return(testError)
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				accountRepo: accountRepositoryMock,
				eventRepo:   eventRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().DeleteAccount(gomock.Eq(ctx), gomock.Eq(accountID)).Return(nil)
				eventRepositoryMock.EXPECT().CreateEvent(gomock.Eq(ctx), gomock.Eq(event)).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configureMock()

			a := NewAccount(tt.fields.accountRepo, nil, tt.fields.eventRepo, nil)
			err := a.DeleteAccount(tt.args.ctx, tt.args.accountID)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestAccount_DepositAccount(t *testing.T) {
	controller := gomock.NewController(t)
	accountRepositoryMock := NewMockAccountRepository(controller)
	eventRepositoryMock := NewMockEventRepository(controller)
	transactionRepositoryMock := NewMockTransactionRepository(controller)

	ctx := context.Background()

	accountID := 1
	amount := float64(100)
	const fromAccountIDForDeposit = 0

	transaction := domain.Transaction{
		ID:          1,
		FromAccount: fromAccountIDForDeposit,
		ToAccount:   accountID,
		Amount:      amount,
		Status:      "PREPARED",
		DateCreated: time.Now(),
	}

	event := domain.Event{
		Type:    domain.DepositEvent,
		Message: "deposit account successful",
		Metadata: map[string]any{
			"account_id": accountID,
			"amount":     amount,
		},
	}

	testError := errors.New("test error")

	type fields struct {
		accountRepo     AccountRepository
		eventRepo       EventRepository
		transactionRepo TransactionRepository
	}
	type args struct {
		ctx       context.Context
		accountID int
		amount    float64
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		configureMock func()
		wantErr       bool
	}{
		{
			name: "checking_existence_account_repository_error",
			fields: fields{
				accountRepo:     accountRepositoryMock,
				eventRepo:       eventRepositoryMock,
				transactionRepo: transactionRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				amount:    amount,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().ExistsAccount(gomock.Eq(ctx), gomock.Eq(accountID)).Return(false, testError)
			},
			wantErr: true,
		},
		{
			name: "account_not_exists_error",
			fields: fields{
				accountRepo:     accountRepositoryMock,
				eventRepo:       eventRepositoryMock,
				transactionRepo: transactionRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				amount:    amount,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().ExistsAccount(gomock.Eq(ctx), gomock.Eq(accountID)).Return(false, nil)
			},
			wantErr: true,
		},
		{
			name: "create_transaction_repository_error",
			fields: fields{
				accountRepo:     accountRepositoryMock,
				eventRepo:       eventRepositoryMock,
				transactionRepo: transactionRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				amount:    amount,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().ExistsAccount(gomock.Eq(ctx), gomock.Eq(accountID)).Return(true, nil)
				transactionRepositoryMock.EXPECT().CreateTransaction(gomock.Eq(ctx), gomock.Eq(fromAccountIDForDeposit), gomock.Eq(accountID), gomock.Eq(amount)).Return(domain.Transaction{}, testError)
			},
			wantErr: true,
		},
		{
			name: "deposit_account_repository_error",
			fields: fields{
				accountRepo:     accountRepositoryMock,
				eventRepo:       eventRepositoryMock,
				transactionRepo: transactionRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				amount:    amount,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().ExistsAccount(gomock.Eq(ctx), gomock.Eq(accountID)).Return(true, nil)
				transactionRepositoryMock.EXPECT().CreateTransaction(gomock.Eq(ctx), gomock.Eq(fromAccountIDForDeposit), gomock.Eq(accountID), gomock.Eq(amount)).Return(transaction, nil)
				accountRepositoryMock.EXPECT().DepositAccount(gomock.Eq(ctx), gomock.Eq(accountID), gomock.Eq(amount)).Return(testError)
			},
			wantErr: true,
		},
		{
			name: "setting_transaction_status_transaction_error",
			fields: fields{
				accountRepo:     accountRepositoryMock,
				eventRepo:       eventRepositoryMock,
				transactionRepo: transactionRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				amount:    amount,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().ExistsAccount(gomock.Eq(ctx), gomock.Eq(accountID)).Return(true, nil)
				transactionRepositoryMock.EXPECT().CreateTransaction(gomock.Eq(ctx), gomock.Eq(fromAccountIDForDeposit), gomock.Eq(accountID), gomock.Eq(amount)).Return(transaction, nil)
				accountRepositoryMock.EXPECT().DepositAccount(gomock.Eq(ctx), gomock.Eq(accountID), gomock.Eq(amount)).Return(nil)
				transactionRepositoryMock.EXPECT().SetTransactionStatusToSent(gomock.Eq(ctx), gomock.Eq(transaction.ID)).Return(testError)
			},
			wantErr: true,
		},
		{
			name: "event_repository_error",
			fields: fields{
				accountRepo:     accountRepositoryMock,
				eventRepo:       eventRepositoryMock,
				transactionRepo: transactionRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				amount:    amount,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().ExistsAccount(gomock.Eq(ctx), gomock.Eq(accountID)).Return(true, nil)
				transactionRepositoryMock.EXPECT().CreateTransaction(gomock.Eq(ctx), gomock.Eq(fromAccountIDForDeposit), gomock.Eq(accountID), gomock.Eq(amount)).Return(transaction, nil)
				accountRepositoryMock.EXPECT().DepositAccount(gomock.Eq(ctx), gomock.Eq(accountID), gomock.Eq(amount)).Return(nil)
				transactionRepositoryMock.EXPECT().SetTransactionStatusToSent(gomock.Eq(ctx), gomock.Eq(transaction.ID)).Return(nil)
				eventRepositoryMock.EXPECT().CreateEvent(gomock.Eq(ctx), gomock.Eq(event)).Return(testError)

			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				accountRepo:     accountRepositoryMock,
				eventRepo:       eventRepositoryMock,
				transactionRepo: transactionRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				amount:    amount,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().ExistsAccount(gomock.Eq(ctx), gomock.Eq(accountID)).Return(true, nil)
				transactionRepositoryMock.EXPECT().CreateTransaction(gomock.Eq(ctx), gomock.Eq(fromAccountIDForDeposit), gomock.Eq(accountID), gomock.Eq(amount)).Return(transaction, nil)
				accountRepositoryMock.EXPECT().DepositAccount(gomock.Eq(ctx), gomock.Eq(accountID), gomock.Eq(amount)).Return(nil)
				transactionRepositoryMock.EXPECT().SetTransactionStatusToSent(gomock.Eq(ctx), gomock.Eq(transaction.ID)).Return(nil)
				eventRepositoryMock.EXPECT().CreateEvent(gomock.Eq(ctx), gomock.Eq(event)).Return(nil)

			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configureMock()

			a := NewAccount(tt.fields.accountRepo, tt.fields.transactionRepo, tt.fields.eventRepo, nil)
			err := a.DepositAccount(tt.args.ctx, tt.args.accountID, tt.args.amount)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestAccount_BlockAccount(t *testing.T) {
	controller := gomock.NewController(t)
	accountRepositoryMock := NewMockAccountRepository(controller)
	eventRepositoryMock := NewMockEventRepository(controller)

	ctx := context.Background()
	accountID := 1
	userID := 1
	testError := errors.New("test error")

	event := domain.Event{
		UserID:  userID,
		Type:    domain.AccountBlockedEvent,
		Message: "account blocked successfully",
		Metadata: map[string]any{
			"account_id": accountID,
		},
	}

	type fields struct {
		accountRepo AccountRepository
		eventRepo   EventRepository
	}
	type args struct {
		ctx       context.Context
		accountID int
		userID    int
	}
	tests := []struct {
		name          string
		fields        fields
		configureMock func()
		args          args
		wantErr       bool
	}{
		{
			name: "account_repository_error",
			fields: fields{
				accountRepo: accountRepositoryMock,
				eventRepo:   eventRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				userID:    userID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().BlockAccount(gomock.Eq(ctx), gomock.Eq(accountID), gomock.Eq(userID)).Return(testError)
			},
			wantErr: true,
		},
		{
			name: "event_repository_error",
			fields: fields{
				accountRepo: accountRepositoryMock,
				eventRepo:   eventRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				userID:    userID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().BlockAccount(gomock.Eq(ctx), gomock.Eq(accountID), gomock.Eq(userID)).Return(nil)
				eventRepositoryMock.EXPECT().CreateEvent(gomock.Eq(ctx), gomock.Eq(event)).Return(testError)
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				accountRepo: accountRepositoryMock,
				eventRepo:   eventRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				userID:    userID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().BlockAccount(gomock.Eq(ctx), gomock.Eq(accountID), gomock.Eq(userID)).Return(nil)
				eventRepositoryMock.EXPECT().CreateEvent(gomock.Eq(ctx), gomock.Eq(event)).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configureMock()

			a := NewAccount(tt.fields.accountRepo, nil, tt.fields.eventRepo, nil)
			err := a.BlockAccount(tt.args.ctx, tt.args.accountID, tt.args.userID)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestAccount_UnblockAccount(t *testing.T) {
	controller := gomock.NewController(t)
	accountRepositoryMock := NewMockAccountRepository(controller)
	eventRepositoryMock := NewMockEventRepository(controller)

	ctx := context.Background()
	accountID := 1
	userID := 1
	testError := errors.New("test error")

	event := domain.Event{
		UserID:  userID,
		Type:    domain.AccountUnblockedEvent,
		Message: "account unblocked successfully",
		Metadata: map[string]any{
			"account_id": accountID,
		},
	}

	type fields struct {
		accountRepo   AccountRepository
		eventRepo     EventRepository
		ibanGenerator RandomGenerator
	}
	type args struct {
		ctx       context.Context
		accountID int
		userID    int
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		configureMock func()
		wantErr       bool
	}{
		{
			name: "account_repository_error",
			fields: fields{
				accountRepo: accountRepositoryMock,
				eventRepo:   eventRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				userID:    userID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().UnblockAccount(gomock.Eq(ctx), gomock.Eq(accountID), gomock.Eq(userID)).Return(testError)
			},
			wantErr: true,
		},
		{
			name: "event_repository_error",
			fields: fields{
				accountRepo: accountRepositoryMock,
				eventRepo:   eventRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				userID:    userID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().UnblockAccount(gomock.Eq(ctx), gomock.Eq(accountID), gomock.Eq(userID)).Return(nil)
				eventRepositoryMock.EXPECT().CreateEvent(gomock.Eq(ctx), event).Return(testError)
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				accountRepo: accountRepositoryMock,
				eventRepo:   eventRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				userID:    userID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().UnblockAccount(gomock.Eq(ctx), gomock.Eq(accountID), gomock.Eq(userID)).Return(nil)
				eventRepositoryMock.EXPECT().CreateEvent(gomock.Eq(ctx), event).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configureMock()

			a := NewAccount(tt.fields.accountRepo, nil, tt.fields.eventRepo, nil)
			err := a.UnblockAccount(tt.args.ctx, tt.args.accountID, tt.args.userID)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
