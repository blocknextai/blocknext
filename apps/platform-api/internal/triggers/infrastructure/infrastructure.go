package infrastructure

import (
	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/secretmanager"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	"github.com/blocknextai/platform-api/internal/triggers/application/triggers/deletetrigger"
	"github.com/blocknextai/platform-api/internal/triggers/application/triggers/getalltriggers"
	"github.com/blocknextai/platform-api/internal/triggers/application/triggers/regeneratewebhooktoken"
	"github.com/blocknextai/platform-api/internal/triggers/application/triggers/updatetrigger"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/triggers/domain/triggers"
	workflowsApplicationWorkflows "github.com/blocknextai/platform-api/internal/workflows/application/workflows"
)

type Handlers struct {
	UpdateTrigger          cqrs.Handler[*updatetrigger.UpdateTriggerCommand, *updatetrigger.UpdateTriggerResponse]
	DeleteTrigger          cqrs.Handler[*deletetrigger.DeleteTriggerCommand, *deletetrigger.DeleteTriggerResponse]
	GetAllTriggers         cqrs.Handler[*getalltriggers.GetAllTriggersQuery, *getalltriggers.GetAllTriggersResponse]
	RegenerateWebhookToken cqrs.Handler[*regeneratewebhooktoken.RegenerateWebhookTokenCommand, *regeneratewebhooktoken.RegenerateWebhookTokenResponse]
}

type RegisterInfrastructureDeps struct {
	TransactionManager database.TransactionManager
	SecretManager      secretmanager.SecretManager

	TriggerRepository triggersDomainTriggers.TriggerRepository
	WorkflowService   workflowsApplicationWorkflows.WorkflowService
}

func RegisterInfrastructure(deps RegisterInfrastructureDeps) *Handlers {
	return &Handlers{
		UpdateTrigger:          cqrs.ValidationBehavior(updatetrigger.New(deps.TriggerRepository, deps.SecretManager, deps.TransactionManager)),
		DeleteTrigger:          cqrs.ValidationBehavior(deletetrigger.New(deps.TriggerRepository, deps.TransactionManager)),
		GetAllTriggers:         cqrs.ValidationBehavior(getalltriggers.New(deps.TriggerRepository, deps.WorkflowService)),
		RegenerateWebhookToken: cqrs.ValidationBehavior(regeneratewebhooktoken.New(deps.TriggerRepository, deps.TransactionManager)),
	}
}
