package notifications

import (
	"strings"
	"time"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	"github.com/google/uuid"
)

type Notification struct {
	database.BaseEntity

	Type         string
	Level        Level
	AudienceType AudienceType
	AudienceID   uuid.UUID
	Title        string
	Body         *string
	Data         map[string]any
	ActionURL    *string
}

func New(
	notificationType string,
	level Level,
	audienceType AudienceType,
	audienceID uuid.UUID,
	title string,
	body *string,
	data map[string]any,
	actionURL *string,
) (*Notification, error) {
	utcNow := time.Now().UTC()

	notification := &Notification{
		BaseEntity: database.BaseEntity{
			ID:        bnuuid.NewV7(),
			CreatedAt: utcNow,
			UpdatedAt: utcNow,
			DeletedAt: nil,
		},
		Type:         notificationType,
		Level:        level,
		AudienceType: audienceType,
		AudienceID:   audienceID,
		Title:        title,
		Body:         body,
		Data:         data,
		ActionURL:    actionURL,
	}

	return notification.validateThenReturn()
}

func (n *Notification) validateThenReturn() (*Notification, error) {
	if strings.TrimSpace(n.Type) == "" {
		return nil, ErrInvalidType
	}

	if !n.Level.IsValid() {
		return nil, ErrInvalidLevel
	}

	if !n.AudienceType.IsValid() {
		return nil, ErrInvalidAudienceType
	}

	if n.AudienceID == uuid.Nil {
		return nil, ErrInvalidAudienceID
	}

	if strings.TrimSpace(n.Title) == "" {
		return nil, ErrInvalidTitle
	}

	return n, nil
}
