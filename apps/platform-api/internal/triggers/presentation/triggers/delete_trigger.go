package triggers

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/triggers/application/triggers/deletetrigger"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type DeleteTriggerRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	TriggerID      uuid.UUID `uri:"triggerId"`
}

func NewDeleteTriggerHandler(handler cqrs.Handler[*deletetrigger.DeleteTriggerCommand, *deletetrigger.DeleteTriggerResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(DeleteTriggerRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &deletetrigger.DeleteTriggerCommand{
			OrganizationID: request.OrganizationID,
			TriggerID:      request.TriggerID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("trigger deleted")))
	}
}
