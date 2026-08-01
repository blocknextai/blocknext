package getallworkflows

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/google/uuid"
)

type GetAllWorkflowsQuery struct {
	OrganizationID uuid.UUID
	Search         resultPkg.SearchRequest
	Pagination     resultPkg.PaginationRequest
}
