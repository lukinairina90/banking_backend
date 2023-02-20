package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lukinairina90/banking_backend/internal/domain"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestUser_SignUp(t *testing.T) {
	controller := gomock.NewController(t)
	userRepositoryMock := NewMockUsersRepository(controller)
	sessionRepositoryMock := NewMockSessionRepository(controller)
	roleRepositoryMock := NewMockRolesRepository(controller)
	eventRepositoryMock := NewMockEventRepository(controller)
	passwordHasher := NewMockPasswordHasher(controller)

	hmacSecret := []byte("secret")
	tokenTtl := time.Duration(12) * time.Hour

	const userRoleName = "user"
	ctx := context.Background()
	inp := domain.SignUpInput{
		Name:     "Bob",
		Surname:  "Marly",
		Email:    "example@gmail.com",
		Password: "marly1234",
	}
	role := domain.Role{
		ID:   2,
		Name: "user",
	}
	password := "73616c74d033e22ae348aeb5660fc2140aec35850c4da997"
	user := domain.User{
		Name:     inp.Name,
		Surname:  inp.Surname,
		Email:    inp.Email,
		Password: password,
		RoleId:   role.ID,
		Blocked:  false,
	}
	testError := errors.New("test error")

	type fields struct {
		userRepo    UsersRepository
		sessionRepo SessionRepository
		roleRepo    RolesRepository
		eventRepo   EventRepository
		hasher      PasswordHasher
		hmacSecret  []byte
		tokenTtl    time.Duration
	}
	type args struct {
		ctx context.Context
		inp domain.SignUpInput
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		configureMock func()
		wantErr       bool
	}{
		{
			name: "user_repository_error",
			fields: fields{
				userRepo:    userRepositoryMock,
				sessionRepo: sessionRepositoryMock,
				roleRepo:    roleRepositoryMock,
				eventRepo:   eventRepositoryMock,
				hasher:      passwordHasher,
				hmacSecret:  hmacSecret,
				tokenTtl:    tokenTtl,
			},
			args: args{
				ctx: ctx,
				inp: inp,
			},
			configureMock: func() {
				userRepositoryMock.EXPECT().Exists(gomock.Eq(ctx), gomock.Eq(inp.Email)).Return(false, testError)
			},
			wantErr: true,
		},
		{
			name: "checking_user_exists_error",
			fields: fields{
				userRepo:    userRepositoryMock,
				sessionRepo: sessionRepositoryMock,
				roleRepo:    roleRepositoryMock,
				eventRepo:   eventRepositoryMock,
				hasher:      passwordHasher,
				hmacSecret:  hmacSecret,
				tokenTtl:    tokenTtl,
			},
			args: args{
				ctx: ctx,
				inp: inp,
			},
			configureMock: func() {
				userRepositoryMock.EXPECT().Exists(gomock.Eq(ctx), gomock.Eq(inp.Email)).Return(true, nil)
			},
			wantErr: true,
		},
		{
			name: "password_hasher_error",
			fields: fields{
				userRepo:    userRepositoryMock,
				sessionRepo: sessionRepositoryMock,
				roleRepo:    roleRepositoryMock,
				eventRepo:   eventRepositoryMock,
				hasher:      passwordHasher,
				hmacSecret:  hmacSecret,
				tokenTtl:    tokenTtl,
			},
			args: args{
				ctx: ctx,
				inp: inp,
			},
			configureMock: func() {
				userRepositoryMock.EXPECT().Exists(gomock.Eq(ctx), gomock.Eq(inp.Email)).Return(false, nil)
				passwordHasher.EXPECT().Hash(gomock.Eq(inp.Password)).Return("", testError)
			},
			wantErr: true,
		},
		{
			name: "role_repository_error",
			fields: fields{
				userRepo:    userRepositoryMock,
				sessionRepo: sessionRepositoryMock,
				roleRepo:    roleRepositoryMock,
				eventRepo:   eventRepositoryMock,
				hasher:      passwordHasher,
				hmacSecret:  hmacSecret,
				tokenTtl:    tokenTtl,
			},
			args: args{
				ctx: ctx,
				inp: inp,
			},
			configureMock: func() {
				userRepositoryMock.EXPECT().Exists(gomock.Eq(ctx), gomock.Eq(inp.Email)).Return(false, nil)
				passwordHasher.EXPECT().Hash(gomock.Eq(inp.Password)).Return(password, nil)
				roleRepositoryMock.EXPECT().GetByName(gomock.Eq(ctx), gomock.Eq(userRoleName)).Return(domain.Role{}, testError)
			},
			wantErr: true,
		},
		{
			name: "user_repository_error",
			fields: fields{
				userRepo:    userRepositoryMock,
				sessionRepo: sessionRepositoryMock,
				roleRepo:    roleRepositoryMock,
				eventRepo:   eventRepositoryMock,
				hasher:      passwordHasher,
				hmacSecret:  hmacSecret,
				tokenTtl:    tokenTtl,
			},
			args: args{
				ctx: ctx,
				inp: inp,
			},
			configureMock: func() {
				userRepositoryMock.EXPECT().Exists(gomock.Eq(ctx), gomock.Eq(inp.Email)).Return(false, nil)
				passwordHasher.EXPECT().Hash(gomock.Eq(inp.Password)).Return(password, nil)
				roleRepositoryMock.EXPECT().GetByName(gomock.Eq(ctx), gomock.Eq(userRoleName)).Return(role, nil)
				userRepositoryMock.EXPECT().Create(gomock.Eq(ctx), gomock.Eq(user)).Return(testError)
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				userRepo:    userRepositoryMock,
				sessionRepo: sessionRepositoryMock,
				roleRepo:    roleRepositoryMock,
				eventRepo:   eventRepositoryMock,
				hasher:      passwordHasher,
				hmacSecret:  hmacSecret,
				tokenTtl:    tokenTtl,
			},
			args: args{
				ctx: ctx,
				inp: inp,
			},
			configureMock: func() {
				userRepositoryMock.EXPECT().Exists(gomock.Eq(ctx), gomock.Eq(inp.Email)).Return(false, nil)
				passwordHasher.EXPECT().Hash(gomock.Eq(inp.Password)).Return(password, nil)
				roleRepositoryMock.EXPECT().GetByName(gomock.Eq(ctx), gomock.Eq(userRoleName)).Return(role, nil)
				userRepositoryMock.EXPECT().Create(gomock.Eq(ctx), gomock.Eq(user)).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configureMock()

			u := NewUsers(tt.fields.userRepo, tt.fields.sessionRepo, tt.fields.roleRepo, tt.fields.eventRepo, tt.fields.hasher, tt.fields.hmacSecret, tt.fields.tokenTtl)
			err := u.SignUp(tt.args.ctx, tt.args.inp)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestUser_SignIn(t *testing.T) {
	controller := gomock.NewController(t)
	userRepositoryMock := NewMockUsersRepository(controller)
	passwordHasher := NewMockPasswordHasher(controller)
	sessionRepositoryMock := NewMockSessionRepository(controller)

	hmacSecret := []byte("secret")
	tokenTtl := time.Duration(12) * time.Hour

	ctx := context.Background()
	inp := domain.SignInInput{
		Email:    "example@gmail.com",
		Password: "marly1234",
	}
	password := "73616c74d033e22ae348aeb5660fc2140aec35850c4da997"
	testError := errors.New("test error")

	user := domain.User{
		ID:       1,
		Name:     "Bob",
		Surname:  "Marly",
		Email:    inp.Email,
		Password: password,
		RoleId:   2,
		Blocked:  false,
	}

	type fields struct {
		userRepo   UsersRepository
		hasher     PasswordHasher
		tokenTtl   time.Duration
		hmacSecret []byte
	}
	type args struct {
		ctx context.Context
		inp domain.SignInInput
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		configureMock func()
		want          bool
		want1         bool
		wantErr       bool
	}{
		{
			name: "password_hasher_error",
			fields: fields{
				userRepo:   userRepositoryMock,
				hasher:     passwordHasher,
				hmacSecret: hmacSecret,
				tokenTtl:   tokenTtl,
			},
			args: args{
				ctx: ctx,
				inp: inp,
			},
			configureMock: func() {
				passwordHasher.EXPECT().Hash(inp.Password).Return("", testError)
			},
			want:    false,
			want1:   false,
			wantErr: true,
		},
		{
			name: "user_repository_error",
			fields: fields{
				userRepo:   userRepositoryMock,
				hasher:     passwordHasher,
				hmacSecret: hmacSecret,
				tokenTtl:   tokenTtl,
			},
			args: args{
				ctx: ctx,
				inp: inp,
			},
			configureMock: func() {
				passwordHasher.EXPECT().Hash(inp.Password).Return(password, nil)
				userRepositoryMock.EXPECT().GetByCredentials(gomock.Eq(ctx), gomock.Eq(inp.Email), gomock.Eq(password)).Return(domain.User{}, testError)
			},
			want:    false,
			want1:   false,
			wantErr: true,
		},
		{
			name: "session_creation_error",
			fields: fields{
				userRepo:   userRepositoryMock,
				hasher:     passwordHasher,
				hmacSecret: hmacSecret,
				tokenTtl:   tokenTtl,
			},
			args: args{
				ctx: ctx,
				inp: inp,
			},
			configureMock: func() {
				passwordHasher.EXPECT().Hash(inp.Password).Return(password, nil)
				userRepositoryMock.EXPECT().GetByCredentials(gomock.Eq(ctx), gomock.Eq(inp.Email), gomock.Eq(password)).Return(user, nil)
				sessionRepositoryMock.EXPECT().Create(gomock.Eq(ctx), gomock.Any()).Return(testError)
			},
			want:    false,
			want1:   false,
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				userRepo:   userRepositoryMock,
				hasher:     passwordHasher,
				hmacSecret: hmacSecret,
				tokenTtl:   tokenTtl,
			},
			args: args{
				ctx: ctx,
				inp: inp,
			},
			configureMock: func() {
				passwordHasher.EXPECT().Hash(inp.Password).Return(password, nil)
				userRepositoryMock.EXPECT().GetByCredentials(gomock.Eq(ctx), gomock.Eq(inp.Email), gomock.Eq(password)).Return(user, nil)
				sessionRepositoryMock.EXPECT().Create(gomock.Eq(ctx), gomock.Any()).Return(nil)
			},
			want:    true,
			want1:   true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configureMock()

			u := NewUsers(tt.fields.userRepo, sessionRepositoryMock, nil, nil, tt.fields.hasher, tt.fields.hmacSecret, tt.fields.tokenTtl)
			got1, got2, err := u.SignIn(tt.args.ctx, tt.args.inp)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, len(got1) != 0)
			assert.Equal(t, tt.want, len(got2) != 0)
		})
	}
}

func TestUser_BlockUser(t *testing.T) {
	controller := gomock.NewController(t)
	userRepositoryMock := NewMockUsersRepository(controller)
	eventRepositoryMock := NewMockEventRepository(controller)

	ctx := context.Background()
	testError := errors.New("test error")

	blockUserID := 1
	userID := 2

	event := domain.Event{
		UserID:  blockUserID,
		Type:    domain.UserBlockedEvent,
		Message: "user blocked successfully",
	}

	type fields struct {
		userRepo  UsersRepository
		eventRepo EventRepository
	}
	type args struct {
		ctx         context.Context
		blockUserID int
		userID      int
	}
	tests := []struct {
		name          string
		fields        fields
		configureMock func()
		args          args
		wantErr       bool
	}{
		{
			name: "user_repository_error",
			fields: fields{
				userRepo:  userRepositoryMock,
				eventRepo: eventRepositoryMock,
			},
			args: args{
				ctx:         ctx,
				blockUserID: blockUserID,
				userID:      userID,
			},
			configureMock: func() {
				userRepositoryMock.EXPECT().BlockUser(gomock.Eq(ctx), gomock.Eq(blockUserID)).Return(testError)
			},
			wantErr: true,
		},
		{
			name: "event_repository_error",
			fields: fields{
				userRepo:  userRepositoryMock,
				eventRepo: eventRepositoryMock,
			},
			args: args{
				ctx:         ctx,
				blockUserID: blockUserID,
				userID:      userID,
			},
			configureMock: func() {
				userRepositoryMock.EXPECT().BlockUser(gomock.Eq(ctx), gomock.Eq(blockUserID)).Return(nil)
				eventRepositoryMock.EXPECT().CreateEvent(gomock.Eq(ctx), gomock.Eq(event)).Return(testError)
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				userRepo:  userRepositoryMock,
				eventRepo: eventRepositoryMock,
			},
			args: args{
				ctx:         ctx,
				blockUserID: blockUserID,
				userID:      userID,
			},
			configureMock: func() {
				userRepositoryMock.EXPECT().BlockUser(gomock.Eq(ctx), gomock.Eq(blockUserID)).Return(nil)
				eventRepositoryMock.EXPECT().CreateEvent(gomock.Eq(ctx), gomock.Eq(event)).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt.configureMock()

		u := NewUsers(tt.fields.userRepo, nil, nil, tt.fields.eventRepo, nil, nil, time.Duration(0))
		err := u.BlockUser(tt.args.ctx, tt.args.blockUserID, tt.args.userID)
		assert.Equal(t, tt.wantErr, err != nil)
	}
}

func TestUser_UnblockUser(t *testing.T) {
	controller := gomock.NewController(t)
	userRepositoryMock := NewMockUsersRepository(controller)
	eventRepositoryMock := NewMockEventRepository(controller)

	ctx := context.Background()
	testError := errors.New("test error")

	userID := 1

	event := domain.Event{
		UserID:  userID,
		Type:    domain.UserUnblockedEvent,
		Message: "user unblocked successfully",
	}

	type fields struct {
		userRepo  UsersRepository
		eventRepo EventRepository
	}
	type args struct {
		ctx    context.Context
		userID int
	}
	tests := []struct {
		name          string
		fields        fields
		configureMock func()
		args          args
		wantErr       bool
	}{
		{
			name: "user_repository_error",
			fields: fields{
				userRepo:  userRepositoryMock,
				eventRepo: eventRepositoryMock,
			},
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			configureMock: func() {
				userRepositoryMock.EXPECT().UnblockUser(gomock.Eq(ctx), gomock.Eq(userID)).Return(testError)
			},
			wantErr: true,
		},
		{
			name: "event_repository_error",
			fields: fields{
				userRepo:  userRepositoryMock,
				eventRepo: eventRepositoryMock,
			},
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			configureMock: func() {
				userRepositoryMock.EXPECT().UnblockUser(gomock.Eq(ctx), gomock.Eq(userID)).Return(nil)
				eventRepositoryMock.EXPECT().CreateEvent(gomock.Eq(ctx), gomock.Eq(event)).Return(testError)
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				userRepo:  userRepositoryMock,
				eventRepo: eventRepositoryMock,
			},
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			configureMock: func() {
				userRepositoryMock.EXPECT().UnblockUser(gomock.Eq(ctx), gomock.Eq(userID)).Return(nil)
				eventRepositoryMock.EXPECT().CreateEvent(gomock.Eq(ctx), gomock.Eq(event)).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configureMock()

			u := NewUsers(tt.fields.userRepo, nil, nil, tt.fields.eventRepo, nil, nil, time.Duration(0))
			err := u.UnblockUser(tt.args.ctx, tt.args.userID)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestUser_CheckBlockUser(t *testing.T) {
	controller := gomock.NewController(t)
	userRepositoryMock := NewMockUsersRepository(controller)

	ctx := context.Background()
	testError := errors.New("test error")
	userID := 1

	type fields struct {
		userRepo UsersRepository
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
		want          bool
		wantErr       bool
	}{
		{
			name: "user_repository_error",
			fields: fields{
				userRepo: userRepositoryMock,
			},
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			configureMock: func() {
				userRepositoryMock.EXPECT().CheckBlockUser(gomock.Eq(ctx), gomock.Eq(userID)).Return(false, testError)
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				userRepo: userRepositoryMock,
			},
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			configureMock: func() {
				userRepositoryMock.EXPECT().CheckBlockUser(gomock.Eq(ctx), gomock.Eq(userID)).Return(false, nil)
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configureMock()

			u := NewUsers(tt.fields.userRepo, nil, nil, nil, nil, nil, time.Duration(0))
			got, err := u.CheckBlockUser(tt.args.ctx, tt.args.userID)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
