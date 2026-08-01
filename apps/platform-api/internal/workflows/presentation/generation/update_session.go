package generation

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/workflows/application/generation/sessions/updatesession"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UpdateSessionRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	SessionID      uuid.UUID `uri:"sessionId"`
	Title          string    `json:"title"`
}

func NewUpdateSessionHandler(handler cqrs.Handler[*updatesession.UpdateSessionCommand, *updatesession.UpdateSessionResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(UpdateSessionRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &updatesession.UpdateSessionCommand{
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
