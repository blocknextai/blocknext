package messages

import (
	"context"

	"github.com/google/uuid"
)

type MessageRepository interface {
	GetAllBySessionID(ctx context.Context, sessionID uuid.UUID, offset int, limit int) ([]*GenerationMessage, int64, error)
	Create(ctx context.Context, message *GenerationMessage) error
	DeleteBySessionID(ctx context.Context, sessionID uuid.UUID) error
}
