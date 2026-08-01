package deleteapikey

import (
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
)

type DeleteAPIKeyCommand struct {
	OwnerType commonDomain.OwnerType
	OwnerID   uuid.UUID
	APIKeyID  uuid.UUID
}
