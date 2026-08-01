package triggers

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/triggers/application/triggers/getalltriggers"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetAllTriggersRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	resultPkg.SearchRequest
	resultPkg.PaginationRequest
}

func NewGetAllTriggersHandler(handler cqrs.Handler[*getalltriggers.GetAllTriggersQuery, *getalltriggers.GetAllTriggersResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetAllTriggersRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		searchRequest := request.SearchRequest.Normalize()
		paginationRequest := request.PaginationRequest.Normalize()

		result, err := handler.Handle(c.RequestCtx(), &getalltriggers.GetAllTriggersQuery{
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
