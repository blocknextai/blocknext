package revokesession

import (
	"github.com/google/uuid"
)

type RevokeSessionCommand struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
}
