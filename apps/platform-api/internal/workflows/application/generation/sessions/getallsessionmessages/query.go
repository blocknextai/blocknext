package getallsessionmessages

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/google/uuid"
)

type GetAllSessionMessagesQuery struct {
	OrganizationID uuid.UUID
	SessionID      uuid.UUID
	Pagination     resultPkg.PaginationRequest
}
