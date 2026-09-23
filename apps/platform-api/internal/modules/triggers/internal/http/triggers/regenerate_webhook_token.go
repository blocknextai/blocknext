package triggers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	triggersUseCases "github.com/blocknextai/platform-api/internal/modules/triggers/internal/usecases/triggers"
)

type RegenerateWebhookTokenRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	TriggerID      uuid.UUID `uri:"triggerId"`
}

func NewRegenerateWebhookTokenHandler(service *triggersUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(RegenerateWebhookTokenRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.RegenerateWebhookToken(c.RequestCtx(), &triggersUseCases.RegenerateWebhookTokenCommand{
			OrganizationID: request.OrganizationID,
			TriggerID:      request.TriggerID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("webhook token regenerated")))
	}
}
