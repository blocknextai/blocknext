package taskexecutions

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	taskexecutionsUseCases "github.com/blocknextai/platform-api/internal/modules/executions/internal/usecases/taskexecutions"
)

type GetAllTaskExecutionsRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	resultPkg.SearchRequest
	resultPkg.PaginationRequest
}

func NewGetAllExecutionsHandler(service *taskexecutionsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetAllTaskExecutionsRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		searchRequest := request.SearchRequest.Normalize()
		paginationRequest := request.PaginationRequest.Normalize()

		result, err := service.GetAllTaskExecutions(c.RequestCtx(), &taskexecutionsUseCases.GetAllTaskExecutionsQuery{
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
