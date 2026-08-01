package oauth2

import (
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
)

type State struct {
	CredentialID uuid.UUID              `json:"credentialId"`
	OwnerType    commonDomain.OwnerType `json:"ownerType"`
	OwnerID      uuid.UUID              `json:"ownerId"`
	CodeVerifier string                 `json:"codeVerifier"`
}
