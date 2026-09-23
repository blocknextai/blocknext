package oauth2

import (
	"github.com/google/uuid"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
)

type State struct {
	CredentialID uuid.UUID              `json:"credentialId"`
	OwnerType    commonDomain.OwnerType `json:"ownerType"`
	OwnerID      uuid.UUID              `json:"ownerId"`
	CodeVerifier string                 `json:"codeVerifier"`
}
