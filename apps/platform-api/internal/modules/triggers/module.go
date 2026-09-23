package triggers

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/secretmanager"
	"github.com/blocknextai/platform-api/internal/common/auth"
	triggersApplicationTriggers "github.com/blocknextai/platform-api/internal/modules/triggers/internal/application/triggers"
	triggersApplicationWebhooks "github.com/blocknextai/platform-api/internal/modules/triggers/internal/application/webhooks"
	triggersHTTP "github.com/blocknextai/platform-api/internal/modules/triggers/internal/http"
	triggersPostgresTriggers "github.com/blocknextai/platform-api/internal/modules/triggers/internal/postgres/triggers"
	triggersUseCases "github.com/blocknextai/platform-api/internal/modules/triggers/internal/usecases"
	workflowsContract "github.com/blocknextai/platform-api/internal/modules/workflows/contract"
)

type Dependencies struct {
	DB                 *sql.DB
	TransactionManager database.TransactionManager
	SecretManager      secretmanager.SecretManager

	WorkflowService workflowsContract.WorkflowService
}

type Module struct {
	TriggerService  triggersApplicationTriggers.TriggerService
	WebhookResolver triggersApplicationWebhooks.WebhookResolver

	useCases *triggersUseCases.Services
}

func NewModule(deps Dependencies) *Module {
	repository := triggersPostgresTriggers.NewTriggerRepository(deps.DB)
	service := triggersApplicationTriggers.NewTriggerService(repository)
	webhookResolver := triggersApplicationWebhooks.NewWebhookResolver(repository, deps.SecretManager)
	useCases := triggersUseCases.NewServices(triggersUseCases.ServiceDependencies{
		TransactionManager: deps.TransactionManager,
		SecretManager:      deps.SecretManager,

		TriggerRepository: repository,
		WorkflowService:   deps.WorkflowService,
	})
	return &Module{
		TriggerService:  service,
		WebhookResolver: webhookResolver,
		useCases:        useCases,
	}
}

type Services struct {
	TriggerService triggersApplicationTriggers.TriggerService
}

type ServicesDependencies struct {
	DB *sql.DB
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		TriggerService: triggersApplicationTriggers.NewTriggerService(
			triggersPostgresTriggers.NewTriggerRepository(deps.DB),
		),
	}
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware) {
	triggersHTTP.RegisterRoutes(router, authMiddleware, m.useCases)
}
