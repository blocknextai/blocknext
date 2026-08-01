package inboxentries

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type InboxEntry struct {
	HandlerKey  string
	EventID     uuid.UUID
	ProcessedAt time.Time
}

func New(handlerKey string, eventID uuid.UUID, processedAt time.Time) (*InboxEntry, error) {
	entry := &InboxEntry{
		HandlerKey:  handlerKey,
		EventID:     eventID,
		ProcessedAt: processedAt,
	}

	return entry.validateThenReturn()
}

func (e *InboxEntry) validateThenReturn() (*InboxEntry, error) {
	if strings.TrimSpace(e.HandlerKey) == "" {
		return nil, ErrHandlerKeyRequired
	}
	return e, nil
}
