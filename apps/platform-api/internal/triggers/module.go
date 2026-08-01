package triggers

import (
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/secretmanager"
	"github.com/blocknextai/platform-api/internal/common/presentation/auth"
	triggersApplicationTriggers "github.com/blocknextai/platform-api/internal/triggers/application/triggers"
	triggersApplicationWebhooks "github.com/blocknextai/platform-api/internal/triggers/application/webhooks"
	triggersInfrastructure "github.com/blocknextai/platform-api/internal/triggers/infrastructure"
	triggersInfrastructureTriggers "github.com/blocknextai/platform-api/internal/triggers/infrastructure/triggers"
	triggersPresentation "github.com/blocknextai/platform-api/internal/triggers/presentation"
	workflowsApplicationWorkflows "github.com/blocknextai/platform-api/internal/workflows/application/workflows"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	DB                 *sql.DB
	TransactionManager database.TransactionManager
	SecretManager      secretmanager.SecretManager

	WorkflowService workflowsApplicationWorkflows.WorkflowService
}

type Module struct {
	TriggerService  triggersApplicationTriggers.TriggerService
	WebhookResolver triggersApplicationWebhooks.WebhookResolver

	handlers *triggersInfrastructure.Handlers
}

func NewModule(deps Dependencies) *Module {
	repository := triggersInfrastructureTriggers.NewTriggerRepository(deps.DB)
	service := triggersApplicationTriggers.NewTriggerService(repository)
	webhookResolver := triggersApplicationWebhooks.NewWebhookResolver(repository, deps.SecretManager)
	handlers := triggersInfrastructure.RegisterInfrastructure(triggersInfrastructure.RegisterInfrastructureDeps{
		TransactionManager: deps.TransactionManager,
		SecretManager:      deps.SecretManager,

		TriggerRepository: repository,
		WorkflowService:   deps.WorkflowService,
	})
	return &Module{
		TriggerService:  service,
		WebhookResolver: webhookResolver,
		handlers:        handlers,
	}
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware) {
	triggersPresentation.RegisterPresentation(router, authMiddleware, m.handlers)
}
