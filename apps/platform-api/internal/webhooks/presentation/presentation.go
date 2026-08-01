package presentation

import (
	webhooksInfrastructure "github.com/blocknextai/platform-api/internal/webhooks/infrastructure"
	webhooksPresentationTriggers "github.com/blocknextai/platform-api/internal/webhooks/presentation/triggers"
	"github.com/gofiber/fiber/v3"
)

func RegisterPresentation(
	router fiber.Router,
	handlers *webhooksInfrastructure.Handlers,
) {
	triggerPath := "/triggers/:source/:token"
	triggerHandler := webhooksPresentationTriggers.NewWebhookHandler(handlers.TriggerProcessor) // TODO: is it correct location??
	router.Get(triggerPath, triggerHandler)
	router.Post(triggerPath, triggerHandler)
}
