package domain

import (
	"errors"
	"time"
)

// eventType alias string
type eventType string

// String stringer interface implementation
func (e eventType) String() string {
	return string(e)
}

// NewEventTypeFromString converting events to a string
func NewEventTypeFromString(et string) (eventType, error) {
	switch et {
	case AccountCreatedEvent.String():
		return AccountCreatedEvent, nil
	case AccountDeletedEvent.String():
		return AccountDeletedEvent, nil
	case AccountBlockedEvent.String():
		return AccountBlockedEvent, nil
	case AccountUnblockedEvent.String():
		return AccountUnblockedEvent, nil
	case CardCreatedEvent.String():
		return CardCreatedEvent, nil
	case UserBlockedEvent.String():
		return UserBlockedEvent, nil
	case UserUnblockedEvent.String():
		return UserUnblockedEvent, nil
	case WithdrawalEvent.String():
		return WithdrawalEvent, nil
	case DepositEvent.String():
		return DepositEvent, nil
	default:
		return "", errors.New("unsupported event type")
	}
}

// constants for events
const (
	AccountCreatedEvent   eventType = "ACCOUNT_CREATED"
	AccountDeletedEvent   eventType = "ACCOUNT_DELETED"
	AccountBlockedEvent   eventType = "ACCOUNT_BLOCKED"
	AccountUnblockedEvent eventType = "ACCOUNT_UNBLOCKED"
	CardCreatedEvent      eventType = "CARD_CREATED"
	UserBlockedEvent      eventType = "USER_BLOCKED"
	UserUnblockedEvent    eventType = "USER_UNBLOCKED"
	WithdrawalEvent       eventType = "WITHDRAWAL"
	DepositEvent          eventType = "DEPOSIT"
)

// Event business layer event definition
type Event struct {
	ID       int
	UserID   int
	Type     eventType
	Message  string
	Metadata map[string]any
	DateTime time.Time
}
