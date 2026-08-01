package createapikey

import (
	apiKeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/apikeys/domain/apikeys"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
)

type CreateAPIKeyCommand struct {
	OwnerType commonDomain.OwnerType
	OwnerID   uuid.UUID
	Name      string
	Scopes    apiKeysDomainAPIKeys.Scopes
}
