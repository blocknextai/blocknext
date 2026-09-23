package webhooks

import (
	"github.com/gofiber/fiber/v3"

	taskrunnerContract "github.com/blocknextai/platform-api/internal/modules/taskrunner/contract"
	webhooksHTTP "github.com/blocknextai/platform-api/internal/modules/webhooks/internal/http"
)

type Dependencies struct {
	TriggerWebhookProcessor taskrunnerContract.WebhookProcessor
}

type Module struct {
	deps Dependencies
}

func NewModule(deps Dependencies) *Module {
	return &Module{deps: deps}
}

func (m *Module) Register(router fiber.Router) {
	webhooksHTTP.RegisterRoutes(router, m.deps.TriggerWebhookProcessor)
}
