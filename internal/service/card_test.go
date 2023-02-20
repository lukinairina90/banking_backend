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

func TestCard_CreateCard(t *testing.T) {
	controller := gomock.NewController(t)
	cardRepositoryMock := NewMockCardRepository(controller)
	userRepositoryMock := NewMockUsersRepository(controller)
	accountRepositoryMock := NewMockAccountRepository(controller)
	eventRepositoryMock := NewMockEventRepository(controller)
	randomGeneratorMock := NewMockRandomGenerator(controller)

	ctx := context.Background()
	accountID := 2
	userID := 1
	checkUserID := 1
	cardNumber := "7679988891629820"
	cardholderName := "Nata Kuka"
	cvvCode := "782"

	expirationDate := time.Now().Add(time.Minute * 5)

	card := domain.Card{
		Id:             1,
		AccountID:      accountID,
		CardNumber:     cardNumber,
		CardholderName: cardholderName,
		ExpirationDate: expirationDate,
		CvvCode:        cvvCode,
	}

	event := domain.Event{
		UserID:  userID,
		Type:    domain.CardCreatedEvent,
		Message: "new card successfully created",
		Metadata: map[string]any{
			"card_id":         card.Id,
			"account_id":      accountID,
			"cardholder_name": card.CardholderName,
			"expiration_date": card.ExpirationDate,
		},
	}

	testError := errors.New("test error")

	type fields struct {
		cardRepo        CardRepository
		userRepo        UsersRepository
		accountRepo     AccountRepository
		eventRepo       EventRepository
		randomGenerator RandomGenerator
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
		want          domain.Card
		wantErr       bool
	}{
		{
			name: "card_repository_error",
			fields: fields{
				cardRepo:        cardRepositoryMock,
				userRepo:        userRepositoryMock,
				accountRepo:     accountRepositoryMock,
				eventRepo:       eventRepositoryMock,
				randomGenerator: randomGeneratorMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				userID:    userID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetUserIDByAccountID(gomock.Eq(ctx), gomock.Eq(accountID)).Return(0, testError)
			},
			want:    domain.Card{},
			wantErr: true,
		},
		{
			name: "account_userID_mismatching",
			fields: fields{
				cardRepo:        cardRepositoryMock,
				userRepo:        userRepositoryMock,
				accountRepo:     accountRepositoryMock,
				eventRepo:       eventRepositoryMock,
				randomGenerator: randomGeneratorMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				userID:    userID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetUserIDByAccountID(gomock.Eq(ctx), gomock.Eq(accountID)).Return(10, testError)
			},
			want:    domain.Card{},
			wantErr: true,
		},
		{
			name: "user_repository_error",
			fields: fields{
				cardRepo:        cardRepositoryMock,
				userRepo:        userRepositoryMock,
				accountRepo:     accountRepositoryMock,
				eventRepo:       eventRepositoryMock,
				randomGenerator: randomGeneratorMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				userID:    userID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetUserIDByAccountID(gomock.Eq(ctx), gomock.Eq(accountID)).Return(checkUserID, nil)
				randomGeneratorMock.EXPECT().GenerateRandomCardNumber().Return(cardNumber)
				randomGeneratorMock.EXPECT().GenerateRandomCvv().Return(cvvCode)
				userRepositoryMock.EXPECT().GetUserNameAndSurnameByID(gomock.Eq(ctx), gomock.Eq(userID)).Return("", testError)
			},
			want:    domain.Card{},
			wantErr: true,
		},
		{
			name: "card_repository_error",
			fields: fields{
				cardRepo:        cardRepositoryMock,
				userRepo:        userRepositoryMock,
				accountRepo:     accountRepositoryMock,
				eventRepo:       eventRepositoryMock,
				randomGenerator: randomGeneratorMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				userID:    userID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetUserIDByAccountID(gomock.Eq(ctx), gomock.Eq(accountID)).Return(checkUserID, nil)
				randomGeneratorMock.EXPECT().GenerateRandomCardNumber().Return(cardNumber)
				randomGeneratorMock.EXPECT().GenerateRandomCvv().Return(cvvCode)
				userRepositoryMock.EXPECT().GetUserNameAndSurnameByID(gomock.Eq(ctx), gomock.Eq(userID)).Return(cardholderName, nil)
				cardRepositoryMock.EXPECT().CreateCard(gomock.Eq(ctx), gomock.Eq(accountID), gomock.Eq(cardNumber), gomock.Eq(cardholderName), gomock.Eq(cvvCode)).Return(domain.Card{}, testError)
			},
			want:    domain.Card{},
			wantErr: true,
		},
		{
			name: "event_repository_error",
			fields: fields{
				cardRepo:        cardRepositoryMock,
				userRepo:        userRepositoryMock,
				accountRepo:     accountRepositoryMock,
				eventRepo:       eventRepositoryMock,
				randomGenerator: randomGeneratorMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				userID:    userID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetUserIDByAccountID(gomock.Eq(ctx), gomock.Eq(accountID)).Return(checkUserID, nil)
				userRepositoryMock.EXPECT().GetUserNameAndSurnameByID(gomock.Eq(ctx), gomock.Eq(userID)).Return(cardholderName, nil)
				randomGeneratorMock.EXPECT().GenerateRandomCardNumber().Return(cardNumber)
				randomGeneratorMock.EXPECT().GenerateRandomCvv().Return(cvvCode)
				cardRepositoryMock.EXPECT().CreateCard(gomock.Eq(ctx), gomock.Eq(accountID), gomock.Eq(cardNumber), gomock.Eq(cardholderName), gomock.Eq(cvvCode)).Return(card, nil)
				eventRepositoryMock.EXPECT().CreateEvent(gomock.Eq(ctx), gomock.Eq(event)).Return(testError)
			},
			want:    domain.Card{},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				cardRepo:        cardRepositoryMock,
				userRepo:        userRepositoryMock,
				accountRepo:     accountRepositoryMock,
				eventRepo:       eventRepositoryMock,
				randomGenerator: randomGeneratorMock,
			},
			args: args{
				ctx:       ctx,
				accountID: accountID,
				userID:    userID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetUserIDByAccountID(gomock.Eq(ctx), gomock.Eq(accountID)).Return(checkUserID, nil)
				userRepositoryMock.EXPECT().GetUserNameAndSurnameByID(gomock.Eq(ctx), gomock.Eq(userID)).Return(cardholderName, nil)
				randomGeneratorMock.EXPECT().GenerateRandomCardNumber().Return(cardNumber)
				randomGeneratorMock.EXPECT().GenerateRandomCvv().Return(cvvCode)
				cardRepositoryMock.EXPECT().CreateCard(gomock.Eq(ctx), gomock.Eq(accountID), gomock.Eq(cardNumber), gomock.Eq(cardholderName), gomock.Eq(cvvCode)).Return(card, nil)
				eventRepositoryMock.EXPECT().CreateEvent(gomock.Eq(ctx), gomock.Eq(event)).Return(nil)
			},
			want:    card,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configureMock()

			c := NewCard(tt.fields.cardRepo, tt.fields.userRepo, tt.fields.accountRepo, tt.fields.eventRepo, tt.fields.randomGenerator)
			got, err := c.CreateCard(tt.args.ctx, tt.args.accountID, tt.args.userID)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCard_GetCardListUser(t *testing.T) {
	controller := gomock.NewController(t)
	cardRepositoryMock := NewMockCardRepository(controller)

	ctx := context.Background()
	userID := 1
	cardNumber := "2880296171523095"
	cvvCode := "268"
	testError := errors.New("test error")
	expirationDate := time.Now().Add(time.Minute * 5)

	card1 := domain.Card{
		Id:             1,
		AccountID:      1,
		CardNumber:     cardNumber,
		CardholderName: "Mark Twen",
		ExpirationDate: expirationDate,
		CvvCode:        cvvCode,
	}
	card2 := domain.Card{
		Id:             2,
		AccountID:      2,
		CardNumber:     cardNumber,
		CardholderName: "Met Kuk",
		ExpirationDate: expirationDate,
		CvvCode:        cvvCode,
	}

	listCard := []domain.Card{card1, card2}

	type fields struct {
		cardRepo CardRepository
	}
	type args struct {
		ctx    context.Context
		userID int
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		configureMock func()
		want          []domain.Card
		wantErr       bool
	}{
		{
			name: "card_repository_error",
			fields: fields{
				cardRepo: cardRepositoryMock,
			},
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			configureMock: func() {
				cardRepositoryMock.EXPECT().GetCardListUser(gomock.Eq(ctx), gomock.Eq(userID)).Return(nil, testError)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				cardRepo: cardRepositoryMock,
			},
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			configureMock: func() {
				cardRepositoryMock.EXPECT().GetCardListUser(gomock.Eq(ctx), gomock.Eq(userID)).Return(listCard, nil)
			},
			want:    listCard,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configureMock()

			c := NewCard(tt.fields.cardRepo, nil, nil, nil, nil)
			got, err := c.GetCardListUser(tt.args.ctx, tt.args.userID)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCard_GetCardListByAccount(t *testing.T) {
	controller := gomock.NewController(t)
	cardRepositoryMock := NewMockCardRepository(controller)

	ctx := context.Background()
	userID := 1
	accountID := 1
	cardNumber := "2880296171523095"
	cvvCode := "268"
	expirationDate := time.Now().Add(time.Minute * 5)

	testError := errors.New("test error")

	card1 := domain.Card{
		Id:             1,
		AccountID:      accountID,
		CardNumber:     cardNumber,
		CardholderName: "Mark Twen",
		ExpirationDate: expirationDate,
		CvvCode:        cvvCode,
	}
	card2 := domain.Card{
		Id:             2,
		AccountID:      accountID,
		CardNumber:     cardNumber,
		CardholderName: "Met Kuk",
		ExpirationDate: expirationDate,
		CvvCode:        cvvCode,
	}

	listCard := []domain.Card{card1, card2}

	type fields struct {
		cardRepo CardRepository
	}
	type args struct {
		ctx       context.Context
		userID    int
		accountID int
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		configureMock func()
		want          []domain.Card
		wantErr       bool
	}{
		{
			name: "card_repository_error",
			fields: fields{
				cardRepo: cardRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				userID:    userID,
				accountID: accountID,
			},
			configureMock: func() {
				cardRepositoryMock.EXPECT().GetCardListByAccount(gomock.Eq(ctx), gomock.Eq(userID), gomock.Eq(accountID)).Return(nil, testError)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				cardRepo: cardRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				userID:    userID,
				accountID: accountID,
			},
			configureMock: func() {
				cardRepositoryMock.EXPECT().GetCardListByAccount(gomock.Eq(ctx), gomock.Eq(userID), gomock.Eq(accountID)).Return(listCard, nil)
			},
			want:    listCard,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configureMock()

			c := NewCard(tt.fields.cardRepo, nil, nil, nil, nil)
			got, err := c.GetCardListByAccount(tt.args.ctx, tt.args.userID, tt.args.accountID)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCard_GetCard(t *testing.T) {
	controller := gomock.NewController(t)
	cardRepositoryMock := NewMockCardRepository(controller)
	accountRepositoryMock := NewMockAccountRepository(controller)

	ctx := context.Background()
	testError := errors.New("test error")

	cardID := 1
	accountID := 1
	userID := 1
	checkUserID := 1
	expirationDate := time.Now().Add(time.Minute * 5)

	card := domain.Card{
		Id:             cardID,
		AccountID:      accountID,
		CardNumber:     "1234567891234567",
		CardholderName: "Mark Twen",
		ExpirationDate: expirationDate,
		CvvCode:        "123",
	}

	type fields struct {
		cardRepo    CardRepository
		accountRepo AccountRepository
	}
	type args struct {
		ctx       context.Context
		cardID    int
		accountID int
		userID    int
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		configureMock func()
		want          domain.Card
		wantErr       bool
	}{
		{
			name: "account_repository_error",
			fields: fields{
				cardRepo:    cardRepositoryMock,
				accountRepo: accountRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				cardID:    cardID,
				accountID: accountID,
				userID:    userID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetUserIDByAccountID(gomock.Eq(ctx), gomock.Eq(accountID)).Return(0, testError)
			},
			want:    domain.Card{},
			wantErr: true,
		},
		{
			name: "matching_userID_error",
			fields: fields{
				cardRepo:    cardRepositoryMock,
				accountRepo: accountRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				cardID:    cardID,
				accountID: accountID,
				userID:    userID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetUserIDByAccountID(gomock.Eq(ctx), gomock.Eq(accountID)).Return(33, nil)
			},
			want:    domain.Card{},
			wantErr: true,
		},
		{
			name: "card_repository_error",
			fields: fields{
				cardRepo:    cardRepositoryMock,
				accountRepo: accountRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				cardID:    cardID,
				accountID: accountID,
				userID:    userID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetUserIDByAccountID(gomock.Eq(ctx), gomock.Eq(accountID)).Return(checkUserID, nil)
				cardRepositoryMock.EXPECT().GetCard(gomock.Eq(ctx), gomock.Eq(cardID), gomock.Eq(accountID)).Return(domain.Card{}, testError)
			},
			want:    domain.Card{},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				cardRepo:    cardRepositoryMock,
				accountRepo: accountRepositoryMock,
			},
			args: args{
				ctx:       ctx,
				cardID:    cardID,
				accountID: accountID,
				userID:    userID,
			},
			configureMock: func() {
				accountRepositoryMock.EXPECT().GetUserIDByAccountID(gomock.Eq(ctx), gomock.Eq(accountID)).Return(checkUserID, nil)
				cardRepositoryMock.EXPECT().GetCard(gomock.Eq(ctx), gomock.Eq(cardID), gomock.Eq(accountID)).Return(card, nil)
			},
			want:    card,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configureMock()

			c := NewCard(tt.fields.cardRepo, nil, tt.fields.accountRepo, nil, nil)
			got, err := c.GetCard(tt.args.ctx, tt.args.cardID, tt.args.accountID, tt.args.userID)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
