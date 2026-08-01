package getalltaskexecutions

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/google/uuid"
)

type GetAllTaskExecutionsQuery struct {
	OrganizationID uuid.UUID
	Search         resultPkg.SearchRequest
	Pagination     resultPkg.PaginationRequest
}
