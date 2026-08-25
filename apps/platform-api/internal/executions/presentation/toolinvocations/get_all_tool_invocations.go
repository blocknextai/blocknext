package toolinvocations

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/executions/application/toolinvocations/getalltoolinvocations"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetAllToolInvocationsRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	resultPkg.SearchRequest
	resultPkg.PaginationRequest
}

func NewGetAllToolInvocationsHandler(handler cqrs.Handler[*getalltoolinvocations.GetAllToolInvocationsQuery, *getalltoolinvocations.GetAllToolInvocationsResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetAllToolInvocationsRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		searchRequest := request.SearchRequest.Normalize()
		paginationRequest := request.PaginationRequest.Normalize()

		result, err := handler.Handle(c.RequestCtx(), &getalltoolinvocations.GetAllToolInvocationsQuery{
			OrganizationID: request.OrganizationID,
			Search:         searchRequest,
			Pagination:     paginationRequest,
		})

		if err != nil {
			return err
		}

		return commonHTTP.RespondPaginated(c, result.Items, result.TotalCount, paginationRequest)
	}
}
