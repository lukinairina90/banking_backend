package service

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lukinairina90/banking_backend/internal/domain"
	"github.com/lukinairina90/banking_backend/internal/transport/rest/messages"
	"github.com/pkg/errors"
)

const userRoleName = "user"

// User business logic layer struct.
type User struct {
	userRepo    UsersRepository
	sessionRepo SessionRepository
	roleRepo    RolesRepository
	eventRepo   EventRepository
	hasher      PasswordHasher

	hmacSecret []byte
	tokenTtl   time.Duration
}

// NewUsers constructor for transaction.
func NewUsers(repo UsersRepository, sessionRepo SessionRepository, roleRepo RolesRepository, eventRepo EventRepository, hasher PasswordHasher, hmacSecret []byte, tokenTtl time.Duration) *User {
	return &User{
		userRepo:    repo,
		sessionRepo: sessionRepo,
		roleRepo:    roleRepo,
		eventRepo:   eventRepo,
		hasher:      hasher,
		hmacSecret:  hmacSecret,
		tokenTtl:    tokenTtl,
	}
}

// SignUp checks if the user already exists, hashes the password, sets the user's role, and creates the user.
func (s *User) SignUp(ctx context.Context, inp domain.SignUpInput) error {
	exists, err := s.userRepo.Exists(ctx, inp.Email)
	if err != nil {
		return errors.Wrap(err, "user existence check error")
	}

	if exists {
		return errors.New("user already exists error")
	}

	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return errors.Wrap(err, "password hash error")
	}

	role, err := s.roleRepo.GetByName(ctx, userRoleName)
	if err != nil {
		return errors.Wrap(err, "role getting error")
	}

	user := domain.User{
		Name:     inp.Name,
		Surname:  inp.Surname,
		Email:    inp.Email,
		Password: password,
		RoleId:   role.ID,
		Blocked:  false,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return errors.Wrap(err, "user creation error")
	}

	return nil
}

// SignIn hashes the password, finds the user by the provided email and password, generates a token.
func (s *User) SignIn(ctx context.Context, inp domain.SignInInput) (string, string, error) {
	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return "", "", errors.Wrap(err, "password hash error")
	}

	user, err := s.userRepo.GetByCredentials(ctx, inp.Email, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", errors.Wrap(messages.ErrUserNotFound, "getting user by credential error")
		}
		return "", "", errors.Wrap(err, "getting user by credential error")
	}

	accessToken, refreshToken, err := s.generateTokens(ctx, user)
	if err != nil {
		return "", "", errors.Wrap(err, "generating tokens error")
	}

	return accessToken, refreshToken, nil
}

// ParseToken parses the token and returns the user id and role id.
func (s *User) ParseToken(_ context.Context, token string) (int, int, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Wrapf(nil, "unexpecting signing method %v", token.Header["alg"])
		}
		return s.hmacSecret, nil
	})

	if err != nil {
		return 0, 0, errors.Wrap(err, "jwt parsing error")
	}

	if !t.Valid {
		return 0, 0, errors.Wrap(nil, "invalid token error")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, 0, errors.Wrap(nil, "invalid claims error")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, 0, errors.Wrap(nil, "invalid subject error")
	}

	subjectParts := strings.Split(subject, ":")
	if len(subjectParts) != 2 {
		return 0, 0, errors.Wrap(nil, "token subject content error")
	}

	userID, err := strconv.Atoi(subjectParts[0])
	if err != nil {
		return 0, 0, errors.Wrap(err, "invalid userID in from token subject")
	}

	roleID, err := strconv.Atoi(subjectParts[1])
	if err != nil {
		return 0, 0, errors.Wrap(err, "invalid userID in from token subject error")
	}

	return userID, roleID, nil
}

// RefreshTokens takes a token from the session, checks it for expiration, takes the user by ID and passes it to token generation, returns accessToken and refreshToken.
func (s *User) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	session, err := s.sessionRepo.Get(ctx, refreshToken)
	if err != nil {
		return "", "", errors.Wrap(err, "getting refresh session error")
	}

	if session.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", errors.Wrap(err, "refresh token expired error")
	}

	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return "", "", errors.Wrap(err, "getting user by id error")
	}

	accessToken, refreshToken, err := s.generateTokens(ctx, user)
	if err != nil {
		return "", "", errors.Wrap(err, "generating tokens error")
	}

	return accessToken, refreshToken, nil
}

// generateTokens generates and returns accessToken, refreshToken.
func (s *User) generateTokens(ctx context.Context, user domain.User) (string, string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d:%d", user.ID, user.RoleId),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenTtl)),
	})

	accessToken, err := t.SignedString(s.hmacSecret)
	if err != nil {
		return "", "", errors.Wrap(err, "creating and returning a complete, signed JWT token error")
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", errors.Wrap(err, "creating new refresh token error")
	}

	if err := s.sessionRepo.Create(ctx, domain.RefreshSession{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}); err != nil {
		return "", "", errors.Wrap(err, "creating refresh session error")
	}

	return accessToken, refreshToken, nil
}

// BlockUser checks if the user is blocking himself, blocks user according to the provided user id and creates an event.
func (s *User) BlockUser(ctx context.Context, blockUserID, userID int) error {
	if blockUserID == userID {
		return errors.New("user cannot block himself error")
	}

	err := s.userRepo.BlockUser(ctx, blockUserID)
	if err != nil {
		return errors.Wrap(err, "block user error")
	}

	event := domain.Event{
		UserID:  blockUserID,
		Type:    domain.UserBlockedEvent,
		Message: "user blocked successfully",
	}

	if err := s.eventRepo.CreateEvent(ctx, event); err != nil {
		return errors.Wrap(err, "user blocking event blocking error")
	}

	return nil
}

// UnblockUser unblocks user according to the provided user id and creates an event.
func (s *User) UnblockUser(ctx context.Context, userID int) error {
	err := s.userRepo.UnblockUser(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "unblock user error")
	}

	event := domain.Event{
		UserID:  userID,
		Type:    domain.UserUnblockedEvent,
		Message: "user unblocked successfully",
	}

	if err := s.eventRepo.CreateEvent(ctx, event); err != nil {
		return errors.Wrap(err, "user unblocking event unblocking error")
	}

	return nil
}

// CheckBlockUser checks if the user is blocked out.
func (s *User) CheckBlockUser(ctx context.Context, userID int) (bool, error) {
	checkBlock, err := s.userRepo.CheckBlockUser(ctx, userID)
	if err != nil {
		return false, errors.Wrap(err, "check block user error")
	}

	return checkBlock, nil
}

// newRefreshToken creates a token from random bytes, converts it to a string, and returns.
func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", errors.Wrap(err, "read generates random bytes error")
	}

	return fmt.Sprintf("%x", b), nil
}
