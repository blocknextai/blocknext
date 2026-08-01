package createcredential

import (
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/credentials/domain/credentials"
	"github.com/google/uuid"
)

type CreateCredentialCommand struct {
	OwnerType  commonDomain.OwnerType
	OwnerID    uuid.UUID
	SourceType credentialsDomainCredentials.SourceType
	Key        string
	Title      string
	Data       map[string]any
}
