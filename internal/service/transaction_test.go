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

func TestTransaction_GetTransactionList(t *testing.T) {
	controller := gomock.NewController(t)
	transactionRepositoryMock := NewMockTransactionRepository(controller)
	accountRepositoryMock := NewMockAccountRepository(controller)

	ctx := context.Background()
	accountID := 1
	userID := 1
	checkUserID := 1
	orderings := domain.Orderings{
		"id": "desc",
	}
	paginator := domain.Paginator{
		Page:    1,
		PerPage: 5,
	}
	testError := errors.New("test error")

	transaction1 := domain.Transaction{
		ID:          1,
		FromAccount: fromAccountIDForDeposit,
		ToAccount:   accountID,
		Amount:      float64(100),
		Status:      "PREPARED",
		DateCreated: time.Now(),
	}
	transaction2 := domain.Transaction{
		ID:          1,
		FromAccount: 1,
		ToAccount:   accountID,
		Amount:      float64(60),
		Status:      "PREPARED",
		DateCreated: time.Now(),
	}

	listTransaction := []domain.Transaction{transaction1, transaction2}

	type fields struct {
		TransactionRepo TransactionRepository
		AccountRepo     AccountRepository
	}
	type args struct {
		ctx       context.Context
		accountID int
		userID    int
		ordering  domain.Orderings
		paginator domain.Paginator
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		configureMock func()
		want          []domain.Transaction
		wantErr       bool
	}{
		{
			name: "account_repository_error",
			fields: fields{
				TransactionRepo: transactionRepositoryMock,
				AccountRepo:     accountRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				userID:    userID,
				ordering:  orderings,
				paginator: paginator,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetUserIDByAccountID(gomock.Eq(ctx), gomock.Eq(accountID)).Return(0, testError)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "matching_userID_error",
			fields: fields{
				TransactionRepo: transactionRepositoryMock,
				AccountRepo:     accountRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				userID:    userID,
				ordering:  orderings,
				paginator: paginator,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetUserIDByAccountID(gomock.Eq(ctx), gomock.Eq(accountID)).Return(33, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "transaction_repository_error",
			fields: fields{
				TransactionRepo: transactionRepositoryMock,
				AccountRepo:     accountRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				userID:    userID,
				ordering:  orderings,
				paginator: paginator,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetUserIDByAccountID(gomock.Eq(ctx), gomock.Eq(accountID)).Return(checkUserID, nil)
				transactionRepositoryMock.EXPECT().GetTransactionList(gomock.Eq(ctx), gomock.Eq(accountID), gomock.Eq(orderings), gomock.Eq(paginator)).Return(nil, testError)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				TransactionRepo: transactionRepositoryMock,
				AccountRepo:     accountRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				userID:    userID,
				ordering:  orderings,
				paginator: paginator,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetUserIDByAccountID(gomock.Eq(ctx), gomock.Eq(accountID)).Return(checkUserID, nil)
				transactionRepositoryMock.EXPECT().GetTransactionList(gomock.Eq(ctx), gomock.Eq(accountID), gomock.Eq(orderings), gomock.Eq(paginator)).Return(listTransaction, nil)
			},
			want:    listTransaction,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configureMock()

			tr := NewTransaction(tt.fields.TransactionRepo, tt.fields.AccountRepo)
			got, err := tr.GetTransactionList(tt.args.ctx, tt.args.accountID, tt.args.userID, tt.args.ordering, tt.args.paginator)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
