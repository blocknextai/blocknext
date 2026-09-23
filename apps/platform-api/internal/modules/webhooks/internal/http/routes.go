package http

import (
	"github.com/gofiber/fiber/v3"

	taskrunnerContract "github.com/blocknextai/platform-api/internal/modules/taskrunner/contract"
	webhooksHTTPTriggers "github.com/blocknextai/platform-api/internal/modules/webhooks/internal/http/triggers"
)

func RegisterRoutes(router fiber.Router, triggerProcessor taskrunnerContract.WebhookProcessor) {
	triggerPath := "/triggers/:source/:token"
	triggerHandler := webhooksHTTPTriggers.NewWebhookHandler(triggerProcessor)
	router.Get(triggerPath, triggerHandler)
	router.Post(triggerPath, triggerHandler)
}
