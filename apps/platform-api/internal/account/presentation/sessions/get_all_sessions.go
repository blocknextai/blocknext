package sessions

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/sessions/getallsessions"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

type GetAllSessionsRequest struct {
	resultPkg.SearchRequest
	resultPkg.PaginationRequest
}

func NewGetAllSessionsHandler(handler cqrs.Handler[*getallsessions.GetAllSessionsQuery, *getallsessions.GetAllSessionsResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetAllSessionsRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)
		sessionID := commonHTTP.GetSessionID(c)

		searchRequest := request.SearchRequest.Normalize()
		paginationRequest := request.PaginationRequest.Normalize()

		result, err := handler.Handle(c.RequestCtx(), &getallsessions.GetAllSessionsQuery{
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
