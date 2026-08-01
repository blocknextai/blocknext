package generation

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/workflows/application/generation/sessions/getallsessionmessages"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetAllSessionMessagesRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	SessionID      uuid.UUID `uri:"sessionId"`
	resultPkg.PaginationRequest
}

func NewGetAllSessionMessagesHandler(handler cqrs.Handler[*getallsessionmessages.GetAllSessionMessagesQuery, *getallsessionmessages.GetAllSessionMessagesResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetAllSessionMessagesRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		paginationRequest := request.Normalize()

		result, err := handler.Handle(c.RequestCtx(), &getallsessionmessages.GetAllSessionMessagesQuery{
			OrganizationID: request.OrganizationID,
			SessionID:      request.SessionID,
			Pagination:     paginationRequest,
		})
		if err != nil {
			return err
		}

		return commonHTTP.RespondPaginated(c, result.Items, result.TotalCount, paginationRequest)
	}
}
