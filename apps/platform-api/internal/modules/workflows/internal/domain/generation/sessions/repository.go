package sessions

import (
	"context"

	"github.com/google/uuid"
)

type SessionRepository interface {
	GetByIDAndOrganizationID(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) (*GenerationSession, error)
	GetAllByOrganizationID(ctx context.Context, organizationID uuid.UUID, searchQuery string, offset int, limit int) ([]*GenerationSession, int64, error)
	Create(ctx context.Context, session *GenerationSession) error
	Update(ctx context.Context, session *GenerationSession) error
	Delete(ctx context.Context, session *GenerationSession) error
}
