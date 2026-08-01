package messages

import (
	"time"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	"github.com/google/uuid"
)

type GenerationMessage struct {
	database.BaseEntity

	SessionID uuid.UUID
	Role      string
	Content   string
	Metadata  map[string]any
}

func New(
	sessionID uuid.UUID,
	role string,
	content string,
	metadata map[string]any,
) (*GenerationMessage, error) {
	utcNow := time.Now().UTC()

	message := &GenerationMessage{
		BaseEntity: database.BaseEntity{
			ID:        bnuuid.NewV7(),
			CreatedAt: utcNow,
			UpdatedAt: utcNow,
			DeletedAt: nil,
		},
		SessionID: sessionID,
		Role:      role,
		Content:   content,
		Metadata:  metadata,
	}

	return message.validateThenReturn()
}

func (m *GenerationMessage) Update() (*GenerationMessage, error) {
	m.UpdatedAt = time.Now().UTC()

	return m.validateThenReturn()
}

func (m *GenerationMessage) Delete() (*GenerationMessage, error) {
	utcNow := time.Now().UTC()

	m.UpdatedAt = utcNow
	m.DeletedAt = new(utcNow)
	return m.validateThenReturn()
}

func (m *GenerationMessage) validateThenReturn() (*GenerationMessage, error) {
	if m.SessionID == uuid.Nil {
		return nil, ErrSessionIDIsRequired
	}
	return m, nil
}
