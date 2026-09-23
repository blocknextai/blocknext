package usecases

import (
	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/secretmanager"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/modules/triggers/internal/domain/triggers"
	triggersUseCases "github.com/blocknextai/platform-api/internal/modules/triggers/internal/usecases/triggers"
	workflowsContract "github.com/blocknextai/platform-api/internal/modules/workflows/contract"
)

type Services struct {
	Triggers *triggersUseCases.Service
}

type ServiceDependencies struct {
	TransactionManager database.TransactionManager
	SecretManager      secretmanager.SecretManager

	TriggerRepository triggersDomainTriggers.TriggerRepository
	WorkflowService   workflowsContract.WorkflowService
}

func NewServices(deps ServiceDependencies) *Services {
	return &Services{
		Triggers: triggersUseCases.NewService(deps.TriggerRepository, deps.TransactionManager, deps.WorkflowService, deps.SecretManager),
	}
}
