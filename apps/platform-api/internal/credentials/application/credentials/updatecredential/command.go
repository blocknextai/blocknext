package updatecredential

import (
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
)

type UpdateCredentialCommand struct {
	ID        uuid.UUID
	OwnerType commonDomain.OwnerType
	OwnerID   uuid.UUID
	Key       string
	Title     string
	Data      map[string]any
}
