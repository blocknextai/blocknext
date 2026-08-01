package webhooks

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/adapters/getalltriggervariables"
	"github.com/gofiber/fiber/v3"
)

func NewGetAllTriggerVariablesHandler(handler cqrs.Handler[*getalltriggervariables.GetAllTriggerVariablesQuery, *getalltriggervariables.GetAllTriggerVariablesResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := handler.Handle(c.RequestCtx(), &getalltriggervariables.GetAllTriggerVariablesQuery{})
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
