package getallnotifications

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/google/uuid"
)

type GetAllNotificationsQuery struct {
	UserID         uuid.UUID
	OrganizationID *uuid.UUID
	Search         resultPkg.SearchRequest
	Pagination     resultPkg.PaginationRequest
}
