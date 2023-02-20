package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lukinairina90/banking_backend/internal/domain"
	"github.com/lukinairina90/banking_backend/internal/repository/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const defaultCardExpirationPeriodYears = 3

// Card repository layer struct.
type Card struct {
	db *sqlx.DB
}

// NewCard constructor for Card repository layer.
func NewCard(db *sqlx.DB) *Card {
	return &Card{db: db}
}

// CreateCard creates a card in the database by provides account ID, cardNumber, cardholderName, cvvCode.
func (r Card) CreateCard(ctx context.Context, accountID int, cardNumber string, cardholderName string, cvvCode string) (domain.Card, error) {
	expirationDate := time.Now().AddDate(defaultCardExpirationPeriodYears, 0, 0)

	fields := logrus.Fields{
		"layer":           "repository",
		"repository":      "Card",
		"method":          "CreateCard",
		"account_id":      accountID,
		"card_number":     cardNumber,
		"cardholder_name": cardholderName,
		"expiration_date": expirationDate,
	}

	query := "insert into cards (account_id, card_number, cardholder_name, expiration_date, cvv_code) values ($1, $2, $3, $4, $5) RETURNING *"

	row := r.db.QueryRowxContext(ctx, query, accountID, cardNumber, cardholderName, expirationDate, cvvCode)
	if row.Err() != nil {
		logrus.WithError(row.Err()).
			WithFields(fields).
			Error("execution inserting into cards query error")

		return domain.Card{}, errors.Wrap(row.Err(), "execution inserting into cards query error")
	}

	card := models.Card{}
	if err := row.StructScan(&card); err != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("scanning row into struct error")

		return domain.Card{}, errors.Wrap(err, "scanning row into struct error")
	}

	domainCard := card.ToDomain()

	return domainCard, nil
}

// GetCard returns the card as provided account ID and card ID.
func (r Card) GetCard(ctx context.Context, id, accountID int) (domain.Card, error) {
	fields := logrus.Fields{
		"layer":      "repository",
		"repository": "Card",
		"method":     "GetCard",
		"id":         id,
		"account_id": accountID,
	}

	query := "SELECT * FROM  cards WHERE id=$1 AND account_id=$2"

	var card models.Card

	row := r.db.QueryRowxContext(ctx, query, id, accountID)
	if err := row.StructScan(&card); err != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("scanning row into struct error")

		return domain.Card{}, errors.Wrap(err, "scanning row into struct error")
	}

	domainCard := card.ToDomain()

	return domainCard, nil
}

// GetCardListUser returns all cards for all user accounts.
func (r Card) GetCardListUser(ctx context.Context, userID int) ([]domain.Card, error) {
	fields := logrus.Fields{
		"layer":      "repository",
		"repository": "Card",
		"method":     "GetCardListUser",
		"user_id":    userID,
	}

	query := "SELECT c.* FROM cards c INNER JOIN accounts a on a.id = c.account_id WHERE a.user_id = $1 ORDER BY a.currency_id DESC"

	rows, err := r.db.QueryxContext(ctx, query, userID)
	if err != nil && rows.Err() != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("execution select cards list query error")

		return nil, errors.Wrap(rows.Err(), "execution select list cards query error")
	}

	domainListCards := make([]domain.Card, 0)
	for rows.Next() {
		var card models.Card
		if err := rows.StructScan(&card); err != nil {
			logrus.WithError(err).
				WithFields(fields).
				Error("scanning row into struct error")

			return nil, errors.Wrap(err, "scanning row into struct error")
		}
		domainListCards = append(domainListCards, card.ToDomain())
	}

	return domainListCards, nil
}

// GetCardListByAccount returns a list of cards for a specific account.
func (r Card) GetCardListByAccount(ctx context.Context, userID, accountID int) ([]domain.Card, error) {
	fields := logrus.Fields{
		"layer":      "repository",
		"repository": "Card",
		"method":     "GetCardListByAccount",
		"user_id":    userID,
		"account_id": accountID,
	}

	query := "SELECT c.* FROM cards c INNER JOIN accounts a on a.id = c.account_id WHERE a.user_id = $1 AND a.id = $2 ORDER BY a.currency_id DESC"

	rows, err := r.db.QueryxContext(ctx, query, userID, accountID)
	if err != nil && rows.Err() != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("execution select list cards query error")

		return nil, errors.Wrap(rows.Err(), "execution select list cards query error")
	}

	domainListCards := make([]domain.Card, 0)
	for rows.Next() {
		var card models.Card
		if err := rows.StructScan(&card); err != nil {
			logrus.WithError(err).
				WithFields(fields).
				Error("scanning row into struct error")

			return nil, errors.Wrap(err, "scanning row into struct error")
		}
		domainListCards = append(domainListCards, card.ToDomain())
	}

	return domainListCards, nil
}
