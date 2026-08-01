package getallcredentials

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
)

type GetAllCredentialsQuery struct {
	OwnerType  commonDomain.OwnerType
	OwnerID    uuid.UUID
	Search     resultPkg.SearchRequest
	Pagination resultPkg.PaginationRequest
}
