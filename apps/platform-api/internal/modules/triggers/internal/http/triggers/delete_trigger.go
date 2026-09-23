package triggers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	triggersUseCases "github.com/blocknextai/platform-api/internal/modules/triggers/internal/usecases/triggers"
)

type DeleteTriggerRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	TriggerID      uuid.UUID `uri:"triggerId"`
}

func NewDeleteTriggerHandler(service *triggersUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(DeleteTriggerRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.DeleteTrigger(c.RequestCtx(), &triggersUseCases.DeleteTriggerCommand{
			OrganizationID: request.OrganizationID,
			TriggerID:      request.TriggerID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("trigger deleted")))
	}
}
