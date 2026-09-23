package refreshtokens

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	GetByTokenHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	Create(ctx context.Context, token *RefreshToken) error
	Update(ctx context.Context, token *RefreshToken) error
	RevokeAllByGrantID(ctx context.Context, grantID uuid.UUID, revokedAt time.Time) error
}
