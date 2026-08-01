package webhooks

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/adapters/getallwebhooksources"
	"github.com/gofiber/fiber/v3"
)

func NewGetAllWebhookSourcesHandler(handler cqrs.Handler[*getallwebhooksources.GetAllWebhookSourcesQuery, *getallwebhooksources.GetAllWebhookSourcesResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		result, err := handler.Handle(c.RequestCtx(), &getallwebhooksources.GetAllWebhookSourcesQuery{})
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
