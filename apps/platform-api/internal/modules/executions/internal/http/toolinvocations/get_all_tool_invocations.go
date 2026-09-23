package toolinvocations

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	toolinvocationsUseCases "github.com/blocknextai/platform-api/internal/modules/executions/internal/usecases/toolinvocations"
)

type GetAllToolInvocationsRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	resultPkg.SearchRequest
	resultPkg.PaginationRequest
}

func NewGetAllToolInvocationsHandler(service *toolinvocationsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetAllToolInvocationsRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		searchRequest := request.SearchRequest.Normalize()
		paginationRequest := request.PaginationRequest.Normalize()

		result, err := service.GetAllToolInvocations(c.RequestCtx(), &toolinvocationsUseCases.GetAllToolInvocationsQuery{
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
