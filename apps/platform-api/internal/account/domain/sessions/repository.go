package sessions

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	GetActiveByUserID(ctx context.Context, userID uuid.UUID, searchQuery string, offset int, limit int) ([]*Session, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Session, error)
	GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Session, error)
	RevokeByID(ctx context.Context, id uuid.UUID, updatedAt time.Time) error
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID, updatedAt time.Time) ([]uuid.UUID, error)
	RevokeAllByTokenFamily(ctx context.Context, tokenFamily uuid.UUID, updatedAt time.Time) ([]uuid.UUID, error)
	UpdateRefreshToken(ctx context.Context, session *Session, expectedGeneration int) error
}
