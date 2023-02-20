package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lukinairina90/banking_backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestEvent_GetEventList(t *testing.T) {
	controller := gomock.NewController(t)
	eventRepositoryMock := NewMockEventRepository(controller)

	ctx := context.Background()
	userID := 1

	testError := errors.New("test error")

	event1 := domain.Event{
		UserID:  userID,
		Type:    domain.CardCreatedEvent,
		Message: "new card successfully created",
		Metadata: map[string]any{
			"user_id": userID,
		},
	}
	event2 := domain.Event{
		UserID:  userID,
		Type:    domain.DepositEvent,
		Message: "new card successfully created",
		Metadata: map[string]any{
			"user_id": userID,
		},
	}

	listEvent := []domain.Event{event1, event2}

	type fields struct {
		eventRepository EventRepository
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
		want          []domain.Event
		wantErr       bool
	}{
		{
			name: "card_repository_error",
			fields: fields{
				eventRepository: eventRepositoryMock,
			},
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			configureMock: func() {
				eventRepositoryMock.EXPECT().GetEventsList(gomock.Eq(ctx), gomock.Eq(userID)).Return(nil, testError)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				eventRepository: eventRepositoryMock,
			},
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			configureMock: func() {
				eventRepositoryMock.EXPECT().GetEventsList(gomock.Eq(ctx), gomock.Eq(userID)).Return(listEvent, nil)
			},
			want:    listEvent,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configureMock()

			e := NewEvent(tt.fields.eventRepository)
			got, err := e.GetEventList(tt.args.ctx, tt.args.userID)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
