package triggers

import (
	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/secretmanager"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/modules/triggers/internal/domain/triggers"
	workflowsContract "github.com/blocknextai/platform-api/internal/modules/workflows/contract"
)

type Service struct {
	triggerRepository  triggersDomainTriggers.TriggerRepository
	transactionManager database.TransactionManager
	workflowService    workflowsContract.WorkflowService
	secretManager      secretmanager.SecretManager
}

func NewService(
	triggerRepository triggersDomainTriggers.TriggerRepository,
	transactionManager database.TransactionManager,
	workflowService workflowsContract.WorkflowService,
	secretManager secretmanager.SecretManager,
) *Service {
	return &Service{
		triggerRepository:  triggerRepository,
		transactionManager: transactionManager,
		workflowService:    workflowService,
		secretManager:      secretManager,
	}
}
