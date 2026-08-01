package outboxmessages

import (
	"strings"
	"time"

	bnuuid "github.com/blocknextai/go-packages/uuid"
	"github.com/google/uuid"
)

type OutboxMessage struct {
	ID            uuid.UUID
	EventName     string
	Payload       []byte
	OccurredAt    time.Time
	Status        Status
	Attempts      int
	NextAttemptAt time.Time
	LockedAt      *time.Time
	LastError     *string
	ProcessedAt   *time.Time
}

func New(eventName string, payload []byte, occurredAt time.Time) (*OutboxMessage, error) {
	message := &OutboxMessage{
		ID:            bnuuid.NewV7(),
		EventName:     eventName,
		Payload:       payload,
		OccurredAt:    occurredAt,
		Status:        StatusPending,
		Attempts:      0,
		NextAttemptAt: occurredAt,
	}

	return message.validateThenReturn()
}

func (m *OutboxMessage) MarkProcessed(now time.Time) (*OutboxMessage, error) {
	if m.Status != StatusProcessing {
		return nil, ErrInvalidStatusTransition
	}

	m.Status = StatusProcessed
	m.ProcessedAt = new(now)
	m.LockedAt = nil

	return m.validateThenReturn()
}

func (m *OutboxMessage) Reschedule(nextAttemptAt time.Time, lastError string) (*OutboxMessage, error) {
	if m.Status != StatusProcessing {
		return nil, ErrInvalidStatusTransition
	}

	m.Status = StatusPending
	m.NextAttemptAt = nextAttemptAt
	m.LastError = new(lastError)
	m.LockedAt = nil

	return m.validateThenReturn()
}

func (m *OutboxMessage) MarkDead(lastError string) (*OutboxMessage, error) {
	if m.Status != StatusProcessing {
		return nil, ErrInvalidStatusTransition
	}

	m.Status = StatusDead
	m.LastError = new(lastError)
	m.LockedAt = nil

	return m.validateThenReturn()
}

func (m *OutboxMessage) validateThenReturn() (*OutboxMessage, error) {
	if strings.TrimSpace(m.EventName) == "" {
		return nil, ErrEventNameRequired
	}
	if len(m.Payload) == 0 {
		return nil, ErrPayloadRequired
	}
	return m, nil
}
