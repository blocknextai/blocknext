package auth

import (
	"context"

	"github.com/google/uuid"
)

type SessionRevocationChecker interface {
	IsSessionRevoked(ctx context.Context, sessionID uuid.UUID) (bool, error)
}
