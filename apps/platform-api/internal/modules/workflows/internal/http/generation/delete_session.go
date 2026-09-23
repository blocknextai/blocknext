package generation

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	generationUseCases "github.com/blocknextai/platform-api/internal/modules/workflows/internal/usecases/generation"
)

type DeleteSessionRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	SessionID      uuid.UUID `uri:"sessionId"`
}

func NewDeleteSessionHandler(service *generationUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(DeleteSessionRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.DeleteSession(c.RequestCtx(), &generationUseCases.DeleteSessionCommand{
			OrganizationID: request.OrganizationID,
			SessionID:      request.SessionID,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("generation session deleted")))
	}
}
