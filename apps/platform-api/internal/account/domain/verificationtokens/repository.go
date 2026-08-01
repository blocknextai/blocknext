package verificationtokens

import (
	"context"

	"github.com/google/uuid"
)

type VerificationTokenRepository interface {
	GetByTokenHash(ctx context.Context, tokenHash string) (*VerificationToken, error)
	Create(ctx context.Context, token *VerificationToken) error
	Update(ctx context.Context, token *VerificationToken) error
	InvalidateAllForUserAndPurpose(ctx context.Context, userID uuid.UUID, purpose Purpose) error
}
