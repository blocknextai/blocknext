package getcredentialsfornodes

import (
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
)

type GetCredentialsForNodesQuery struct {
	OwnerType commonDomain.OwnerType
	OwnerID   uuid.UUID
	NodeIDs   []string
}
