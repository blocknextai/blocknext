package generation

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	generationUseCases "github.com/blocknextai/platform-api/internal/modules/workflows/internal/usecases/generation"
)

type GetAllSessionMessagesRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	SessionID      uuid.UUID `uri:"sessionId"`
	resultPkg.PaginationRequest
}

func NewGetAllSessionMessagesHandler(service *generationUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetAllSessionMessagesRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		paginationRequest := request.Normalize()

		result, err := service.GetAllSessionMessages(c.RequestCtx(), &generationUseCases.GetAllSessionMessagesQuery{
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
