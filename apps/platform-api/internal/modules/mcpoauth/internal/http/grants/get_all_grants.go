package grants

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	grantsUseCases "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/usecases/grants"
)

type GetAllUserGrantsRequest struct {
	resultPkg.SearchRequest
	resultPkg.PaginationRequest
}

func NewGetAllUserGrantsHandler(service *grantsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetAllUserGrantsRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		searchRequest := request.SearchRequest.Normalize()
		paginationRequest := request.PaginationRequest.Normalize()

		result, err := service.GetAllGrants(c.RequestCtx(), &grantsUseCases.GetAllGrantsQuery{
			UserID:     commonHTTP.GetUserID(c),
			Search:     searchRequest,
			Pagination: paginationRequest,
		})
		if err != nil {
			return err
		}

		return commonHTTP.RespondPaginated(c, result.Items, result.TotalCount, paginationRequest)
	}
}
