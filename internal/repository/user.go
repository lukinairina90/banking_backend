package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lukinairina90/banking_backend/internal/domain"
	"github.com/lukinairina90/banking_backend/internal/repository/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Users repository layer struct.
type Users struct {
	db *sqlx.DB
}

// NewUsers constructor for Users repository layer.
func NewUsers(db *sqlx.DB) *Users {
	return &Users{db: db}
}

// CheckBlockUser checks if the user is blocked.
func (r Users) CheckBlockUser(ctx context.Context, userID int) (bool, error) {
	fields := logrus.Fields{
		"layer":      "repository",
		"repository": "Users",
		"method":     "CheckBlockUser",
		"user_id":    userID,
	}

	var checkBlock bool

	query := "SELECT blocked FROM users WHERE id = $1"

	if err := r.db.QueryRowContext(ctx, query, userID).Scan(&checkBlock); err != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("execution getting blocked from users query error")

		return false, errors.Wrap(err, "execution getting blocked from users query error")
	}

	if checkBlock {
		return true, nil
	} else {
		return false, nil
	}
}

// GetUserNameAndSurnameByID returns username and surname by provided user ID.
func (r Users) GetUserNameAndSurnameByID(ctx context.Context, userID int) (string, error) {
	fields := logrus.Fields{
		"layer":      "repository",
		"repository": "Users",
		"method":     "GetUserNameAndSurnameByID",
		"user_id":    userID,
	}

	query := "select name, surname from users where id = $1"

	rows, err := r.db.QueryxContext(ctx, query, userID)
	if err != nil && rows.Err() != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("execution getting name and surname from query users error")

		return "", errors.Wrap(err, "execution getting name and surname from query users error")
	}

	var nameSurname models.NameSurname
	for rows.Next() {
		if err := rows.StructScan(&nameSurname); err != nil {
			logrus.WithError(err).
				WithFields(fields).
				Error("scanning row into struct error")

			return "", errors.Wrap(err, "scanning row into struct error")
		}
	}

	fullNameStr := nameSurname.Name + " " + nameSurname.Surname

	return fullNameStr, nil
}

// Exists checks if the user already exists.
func (r Users) Exists(ctx context.Context, email string) (bool, error) {
	fields := logrus.Fields{
		"layer":      "repository",
		"repository": "Users",
		"method":     "Exists",
		"email":      email,
	}

	var exists bool

	query := "select exists(select id from users where email = $1)"

	if err := r.db.QueryRowContext(ctx, query, email).Scan(&exists); err != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("execution checking exists user query error")

		return false, errors.Wrap(err, "execution checking exists user query error")
	}

	return exists, nil
}

// Create creates a user in the database according to the data provided by the user.
func (r Users) Create(ctx context.Context, user domain.User) error {
	fields := logrus.Fields{
		"layer":      "repository",
		"repository": "Users",
		"method":     "Create",
		"user":       user,
	}

	query := "insert into users (name, surname, email, password, role_id, blocked, registered_at) values ($1, $2, $3, $4, $5, $6, now())"

	_, err := r.db.ExecContext(ctx, query, user.Name, user.Surname, user.Email, user.Password, user.RoleId, user.Blocked)
	if err != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("execution creating user query error")

		return errors.Wrap(err, "execution creating user query error")
	}

	return nil
}

// GetByCredentials returns the user according to the provided email and password.
func (r Users) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	fields := logrus.Fields{
		"layer":      "repository",
		"repository": "Users",
		"method":     "GetByCredentials",
		"email":      email,
		"password":   password,
	}

	var user models.User

	query := "select id, name, surname, email, password, role_id, blocked, registered_at from users where email = $1 and password = $2"

	err := r.db.QueryRowContext(ctx, query, email, password).
		Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.RoleId, &user.Blocked, &user.RegisteredAt)
	if err != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("execution selecting user by credentials query error")

		return domain.User{}, errors.Wrap(err, "execution selecting user by credentials query error")
	}

	domainUser := user.ToDomain()

	return domainUser, nil
}

// GetByID returns the user according to the provided id.
func (r Users) GetByID(ctx context.Context, id int) (domain.User, error) {
	fields := logrus.Fields{
		"layer":      "repository",
		"repository": "Users",
		"method":     "GetByID",
		"id":         id,
	}

	var user models.User

	query := "select id, name, surname, email, password, role_id, blocked, registered_at from users where id = $1"

	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.RoleId, &user.Blocked, &user.RegisteredAt)
	if err != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("execution getting user by ID query error")

		return domain.User{}, errors.Wrap(err, "execution getting user by ID query error")
	}

	domainUser := user.ToDomain()

	return domainUser, nil

}

// BlockUser blocks user according to the provided user id.
func (r Users) BlockUser(ctx context.Context, userID int) error {
	fields := logrus.Fields{
		"layer":      "repository",
		"repository": "Users",
		"method":     "BlockUser",
		"user_id":    userID,
	}

	query := "update users set blocked = $1 where id = $2"

	if _, err := r.db.ExecContext(ctx, query, block, userID); err != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("execution updating blocked from user by ID query error")

		return errors.Wrap(err, "execution updating blocked from user by ID query error")
	}

	return nil
}

// UnblockUser unblocks user according to the provided user id.
func (r Users) UnblockUser(ctx context.Context, userID int) error {
	fields := logrus.Fields{
		"layer":      "repository",
		"repository": "Users",
		"method":     "UnblockUser",
		"user_id":    userID,
	}

	query := "update users set blocked= $1 where id = $2"

	if _, err := r.db.ExecContext(ctx, query, unblock, userID); err != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("execution updating blocked from user by ID query error")

		return errors.Wrap(err, "execution updating blocked from user by ID query error")
	}

	return nil
}
