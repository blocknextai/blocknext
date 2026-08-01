package revokeallsessions

import (
	"github.com/google/uuid"
)

type RevokeAllSessionsCommand struct {
	UserID uuid.UUID
}
