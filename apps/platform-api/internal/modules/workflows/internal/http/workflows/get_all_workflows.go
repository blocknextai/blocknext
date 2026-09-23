package workflows

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	workflowsUseCases "github.com/blocknextai/platform-api/internal/modules/workflows/internal/usecases/workflows"
)

type GetAllWorkflowsRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	resultPkg.SearchRequest
	resultPkg.PaginationRequest
}

func NewGetAllWorkflowsHandler(service *workflowsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetAllWorkflowsRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		searchRequest := request.SearchRequest.Normalize()
		paginationRequest := request.PaginationRequest.Normalize()

		result, err := service.GetAllWorkflows(c.RequestCtx(), &workflowsUseCases.GetAllWorkflowsQuery{
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
