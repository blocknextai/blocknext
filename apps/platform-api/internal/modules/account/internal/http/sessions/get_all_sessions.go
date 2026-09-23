package sessions

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	sessionsUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/sessions"
)

type GetAllSessionsRequest struct {
	resultPkg.SearchRequest
	resultPkg.PaginationRequest
}

func NewGetAllSessionsHandler(service *sessionsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetAllSessionsRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)
		sessionID := commonHTTP.GetSessionID(c)

		searchRequest := request.SearchRequest.Normalize()
		paginationRequest := request.PaginationRequest.Normalize()

		result, err := service.GetAllSessions(c.RequestCtx(), &sessionsUseCases.GetAllSessionsQuery{
			UserID:     userID,
			SessionID:  sessionID,
			Search:     searchRequest,
			Pagination: paginationRequest,
		})

		if err != nil {
			return err
		}

		return commonHTTP.RespondPaginated(c, result.Items, result.TotalCount, paginationRequest)
	}
}
