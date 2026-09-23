package webhooks

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	adaptersUseCases "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/usecases/adapters"
)

func NewGetAllTriggerVariablesHandler(service *adaptersUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := service.GetAllTriggerVariables(c.RequestCtx(), &adaptersUseCases.GetAllTriggerVariablesQuery{})
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
