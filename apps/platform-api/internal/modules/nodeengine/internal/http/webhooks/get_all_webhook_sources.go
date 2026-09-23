package webhooks

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	adaptersUseCases "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/usecases/adapters"
)

func NewGetAllWebhookSourcesHandler(service *adaptersUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := service.GetAllWebhookSources(c.RequestCtx(), &adaptersUseCases.GetAllWebhookSourcesQuery{})
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
