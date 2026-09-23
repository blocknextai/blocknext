package usecases

import (
	"github.com/blocknextai/go-packages/database"
	accountContract "github.com/blocknextai/platform-api/internal/modules/account/contract"
	apikeysContract "github.com/blocknextai/platform-api/internal/modules/apikeys/contract"
	executionsApplicationNodeExecutions "github.com/blocknextai/platform-api/internal/modules/executions/internal/application/nodeexecutions"
	"github.com/blocknextai/platform-api/internal/modules/executions/internal/application/taskexecutions/workflowresolver"
	"github.com/blocknextai/platform-api/internal/modules/executions/internal/domain/taskexecutions"
	"github.com/blocknextai/platform-api/internal/modules/executions/internal/domain/toolinvocations"
	taskexecutionsUseCases "github.com/blocknextai/platform-api/internal/modules/executions/internal/usecases/taskexecutions"
	toolinvocationsUseCases "github.com/blocknextai/platform-api/internal/modules/executions/internal/usecases/toolinvocations"
	organizationsContract "github.com/blocknextai/platform-api/internal/modules/organizations/contract"
	workflowsContract "github.com/blocknextai/platform-api/internal/modules/workflows/contract"
)

type Services struct {
	TaskExecutions  *taskexecutionsUseCases.Service
	Toolinvocations *toolinvocationsUseCases.Service
}

type ServiceDependencies struct {
	TransactionManager database.TransactionManager

	TaskExecutionRepository  taskexecutions.TaskExecutionRepository
	ToolInvocationRepository toolinvocations.ToolInvocationRepository
	APIKeyService            apikeysContract.APIKeyService
	NodeExecutionService     executionsApplicationNodeExecutions.NodeExecutionService
	WorkflowService          workflowsContract.WorkflowService
	OrganizationUserService  organizationsContract.OrganizationUserService
	UserService              accountContract.UserService
	LinkedAccountService     accountContract.LinkedAccountService
}

func NewServices(deps ServiceDependencies) *Services {
	workflowResolver := workflowresolver.New(deps.WorkflowService)

	return &Services{
		TaskExecutions:  taskexecutionsUseCases.NewService(deps.TaskExecutionRepository, deps.TransactionManager, workflowResolver, deps.NodeExecutionService, deps.OrganizationUserService, deps.UserService, deps.LinkedAccountService),
		Toolinvocations: toolinvocationsUseCases.NewService(deps.ToolInvocationRepository, deps.APIKeyService),
	}
}
