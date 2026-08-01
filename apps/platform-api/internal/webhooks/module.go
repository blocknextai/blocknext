package webhooks

import (
	taskRunnerApplicationWebhooks "github.com/blocknextai/platform-api/internal/taskrunner/application/webhooks"
	webhooksInfrastructure "github.com/blocknextai/platform-api/internal/webhooks/infrastructure"
	webhooksPresentation "github.com/blocknextai/platform-api/internal/webhooks/presentation"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	TriggerWebhookProcessor taskRunnerApplicationWebhooks.WebhookProcessor
}

type Module struct {
	handlers *webhooksInfrastructure.Handlers
}

func NewModule(deps Dependencies) *Module {
	handlers := webhooksInfrastructure.RegisterInfrastructure(webhooksInfrastructure.RegisterInfrastructureDeps{
		TriggerWebhookProcessor: deps.TriggerWebhookProcessor,
	})
	return &Module{handlers: handlers}
}

func (m *Module) Register(router fiber.Router) {
	webhooksPresentation.RegisterPresentation(router, m.handlers)
}
