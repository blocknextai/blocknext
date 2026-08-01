package apikeys

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/apikeys/application/apikeys/getallapikeys"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
)

type GetAllOrganizationAPIKeysRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	resultPkg.SearchRequest
	resultPkg.PaginationRequest
}

func NewGetAllOrganizationAPIKeysHandler(handler cqrs.Handler[*getallapikeys.GetAllAPIKeysQuery, *getallapikeys.GetAllAPIKeysResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetAllOrganizationAPIKeysRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		searchRequest := request.SearchRequest.Normalize()
		paginationRequest := request.PaginationRequest.Normalize()

		result, err := handler.Handle(c.RequestCtx(), &getallapikeys.GetAllAPIKeysQuery{
			OwnerType:  commonDomain.OwnerTypeOrganization,
			OwnerID:    request.OrganizationID,
			Search:     searchRequest,
			Pagination: paginationRequest,
		})
		if err != nil {
			return err
		}

		return commonHTTP.RespondPaginated(c, result.Items, result.TotalCount, paginationRequest)
	}
}
