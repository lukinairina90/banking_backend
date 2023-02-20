package service

import (
	"context"

	"github.com/lukinairina90/banking_backend/internal/domain"
	"github.com/pkg/errors"
)

// Event business logic layer struct.
type Event struct {
	eventRepository EventRepository
}

// NewEvent constructor for Event.
func NewEvent(eventRepository EventRepository) *Event {
	return &Event{
		eventRepository: eventRepository,
	}
}

// GetEventList returns all events for user.
func (e Event) GetEventList(ctx context.Context, userID int) ([]domain.Event, error) {
	list, err := e.eventRepository.GetEventsList(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "getting events list error")
	}

	return list, nil
}
