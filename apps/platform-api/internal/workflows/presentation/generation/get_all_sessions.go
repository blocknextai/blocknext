package generation

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/workflows/application/generation/sessions/getallsessions"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetAllSessionsRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	resultPkg.SearchRequest
	resultPkg.PaginationRequest
}

func NewGetAllSessionsHandler(handler cqrs.Handler[*getallsessions.GetAllSessionsQuery, *getallsessions.GetAllSessionsResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetAllSessionsRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		searchRequest := request.SearchRequest.Normalize()
		paginationRequest := request.PaginationRequest.Normalize()

		result, err := handler.Handle(c.RequestCtx(), &getallsessions.GetAllSessionsQuery{
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
