package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lukinairina90/banking_backend/internal/domain"
	"github.com/lukinairina90/banking_backend/internal/repository/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Event repository layer struct.
type Event struct {
	db *sqlx.DB
}

// NewEvent constructor for Event repository layer.
func NewEvent(db *sqlx.DB) *Event {
	return &Event{
		db: db,
	}
}

// CreateEvent creates a event in the database by provides event.
func (e Event) CreateEvent(ctx context.Context, event domain.Event) error {
	fields := logrus.Fields{
		"layer":      "repository",
		"repository": "Event",
		"method":     "CreateEvent",
		"event":      event,
	}

	mEvent := models.Event{
		UserID:   event.UserID,
		Type:     event.Type.String(),
		Message:  event.Message,
		Metadata: event.Metadata,
	}

	query := "INSERT INTO event (user_id, type, metadata, time) VALUES ($1, $2, $3, NOW())"

	if _, err := e.db.ExecContext(ctx, query, mEvent.UserID, mEvent.Type, mEvent.Metadata); err != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("execution event insertion query error")

		return errors.Wrap(err, "execution event insertion query error")
	}

	return nil
}

// GetEventsList provides event list.
func (e Event) GetEventsList(ctx context.Context, userID int) ([]domain.Event, error) {
	fields := logrus.Fields{
		"layer":      "repository",
		"repository": "Event",
		"method":     "GetEventsList",
		"user_id":    userID,
	}

	query := "SELECT * FROM event WHERE user_id = $1 ORDER BY time DESC"

	rows, err := e.db.QueryxContext(ctx, query, userID)
	if err != nil {
		logrus.WithError(err).
			WithFields(fields).
			Error("execution getting events list query error")

		return nil, errors.Wrap(err, "execution getting events list query error")
	}

	var eventsList []domain.Event
	for rows.Next() {
		var mEvent models.Event
		if err := rows.StructScan(&mEvent); err != nil {
			logrus.WithError(err).
				WithFields(fields).
				Error("scanning event row error")

			return nil, errors.Wrap(err, "scanning event row error")
		}

		eventType, err := domain.NewEventTypeFromString(mEvent.Type)
		if err != nil {
			logrus.WithError(err).
				WithFields(fields).
				Error("unsupported event type in row")

			return nil, errors.Wrap(err, "unsupported event type in row")
		}

		eventsList = append(eventsList, domain.Event{
			ID:       mEvent.ID,
			UserID:   mEvent.UserID,
			Type:     eventType,
			Message:  mEvent.Message,
			Metadata: mEvent.Metadata,
			DateTime: mEvent.DateTime,
		})
	}

	return eventsList, nil
}
