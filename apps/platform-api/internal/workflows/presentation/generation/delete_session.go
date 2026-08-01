package generation

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/workflows/application/generation/sessions/deletesession"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type DeleteSessionRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	SessionID      uuid.UUID `uri:"sessionId"`
}

func NewDeleteSessionHandler(handler cqrs.Handler[*deletesession.DeleteSessionCommand, *deletesession.DeleteSessionResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(DeleteSessionRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &deletesession.DeleteSessionCommand{
			OrganizationID: request.OrganizationID,
			SessionID:      request.SessionID,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("generation session deleted")))
	}
}
