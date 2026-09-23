package generation

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	generationUseCases "github.com/blocknextai/platform-api/internal/modules/workflows/internal/usecases/generation"
)

type UpdateSessionRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	SessionID      uuid.UUID `uri:"sessionId"`
	Title          string    `json:"title"`
}

func NewUpdateSessionHandler(service *generationUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(UpdateSessionRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.UpdateSession(c.RequestCtx(), &generationUseCases.UpdateSessionCommand{
			OrganizationID: request.OrganizationID,
			SessionID:      request.SessionID,
			Title:          request.Title,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("generation session updated")))
	}
}
