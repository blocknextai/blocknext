package triggers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/modules/triggers/internal/domain/triggers"
	triggersUseCases "github.com/blocknextai/platform-api/internal/modules/triggers/internal/usecases/triggers"
)

type UpdateTriggerRequest struct {
	OrganizationID uuid.UUID                             `uri:"organizationId"`
	TriggerID      uuid.UUID                             `uri:"triggerId"`
	IsActive       *bool                                 `json:"isActive,omitempty"`
	CronPattern    *string                               `json:"cronPattern,omitempty"`
	Timezone       *string                               `json:"timezone,omitempty"`
	WebhookSecret  *string                               `json:"webhookSecret,omitempty"`
	RuntimeConfig  *triggersDomainTriggers.RuntimeConfig `json:"runtimeConfig,omitempty"`
}

func NewUpdateTriggerHandler(service *triggersUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(UpdateTriggerRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.UpdateTrigger(c.RequestCtx(), &triggersUseCases.UpdateTriggerCommand{
			OrganizationID: request.OrganizationID,
			TriggerID:      request.TriggerID,
			IsActive:       request.IsActive,
			CronPattern:    request.CronPattern,
			Timezone:       request.Timezone,
			WebhookSecret:  request.WebhookSecret,
			RuntimeConfig:  request.RuntimeConfig,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("trigger updated")))
	}
}
