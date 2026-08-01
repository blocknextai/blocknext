package addlinkedaccount

import (
	"github.com/blocknextai/platform-api/internal/account/domain"
	"github.com/google/uuid"
)

type AddLinkedAccountCommand struct {
	UserID       uuid.UUID
	AuthProvider domain.AuthProvider
	Payload      map[string]any
}
