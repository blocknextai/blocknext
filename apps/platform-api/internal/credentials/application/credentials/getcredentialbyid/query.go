package getcredentialbyid

import (
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
)

type GetCredentialByIDQuery struct {
	OwnerType    commonDomain.OwnerType
	OwnerID      uuid.UUID
	CredentialID uuid.UUID
}
