package triggers

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/triggers/application/triggers/regeneratewebhooktoken"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type RegenerateWebhookTokenRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	TriggerID      uuid.UUID `uri:"triggerId"`
}

func NewRegenerateWebhookTokenHandler(handler cqrs.Handler[*regeneratewebhooktoken.RegenerateWebhookTokenCommand, *regeneratewebhooktoken.RegenerateWebhookTokenResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(RegenerateWebhookTokenRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &regeneratewebhooktoken.RegenerateWebhookTokenCommand{
			OrganizationID: request.OrganizationID,
			TriggerID:      request.TriggerID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("webhook token regenerated")))
	}
}
