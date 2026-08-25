package infrastructure

import (
	"github.com/blocknextai/go-packages/database"
	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/account/application/linkedaccounts"
	accountApplicationUsers "github.com/blocknextai/platform-api/internal/account/application/users"
	apiKeysApplicationAPIKeys "github.com/blocknextai/platform-api/internal/apikeys/application/apikeys"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	executionsApplicationNodeExecutions "github.com/blocknextai/platform-api/internal/executions/application/nodeexecutions"
	"github.com/blocknextai/platform-api/internal/executions/application/taskexecutions/bulkdeletetaskexecutions"
	"github.com/blocknextai/platform-api/internal/executions/application/taskexecutions/deletetaskexecution"
	"github.com/blocknextai/platform-api/internal/executions/application/taskexecutions/getalltaskexecutions"
	"github.com/blocknextai/platform-api/internal/executions/application/taskexecutions/gettaskexecutionbyid"
	"github.com/blocknextai/platform-api/internal/executions/application/taskexecutions/workflowresolver"
	"github.com/blocknextai/platform-api/internal/executions/application/toolinvocations/getalltoolinvocations"
	"github.com/blocknextai/platform-api/internal/executions/application/toolinvocations/gettoolinvocationbyid"
	"github.com/blocknextai/platform-api/internal/executions/domain/taskexecutions"
	"github.com/blocknextai/platform-api/internal/executions/domain/toolinvocations"
	organizationsApplicationOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/application/organizationusers"
	workflowsApplicationWorkflows "github.com/blocknextai/platform-api/internal/workflows/application/workflows"
)

type Handlers struct {
	GetAllTaskExecutions     cqrs.Handler[*getalltaskexecutions.GetAllTaskExecutionsQuery, *getalltaskexecutions.GetAllTaskExecutionsResponse]
	GetTaskExecutionByID     cqrs.Handler[*gettaskexecutionbyid.GetTaskExecutionByIDQuery, *gettaskexecutionbyid.GetTaskExecutionByIDResponse]
	DeleteTaskExecution      cqrs.Handler[*deletetaskexecution.DeleteTaskExecutionCommand, *deletetaskexecution.DeleteTaskExecutionResponse]
	BulkDeleteTaskExecutions cqrs.Handler[*bulkdeletetaskexecutions.BulkDeleteTaskExecutionsCommand, *bulkdeletetaskexecutions.BulkDeleteTaskExecutionsResponse]
	GetAllToolInvocations    cqrs.Handler[*getalltoolinvocations.GetAllToolInvocationsQuery, *getalltoolinvocations.GetAllToolInvocationsResponse]
	GetToolInvocationByID    cqrs.Handler[*gettoolinvocationbyid.GetToolInvocationByIDQuery, *gettoolinvocationbyid.GetToolInvocationByIDResponse]
}

type RegisterInfrastructureDeps struct {
	TransactionManager database.TransactionManager

	TaskExecutionRepository  taskexecutions.TaskExecutionRepository
	ToolInvocationRepository toolinvocations.ToolInvocationRepository
	APIKeyService            apiKeysApplicationAPIKeys.APIKeyService
	NodeExecutionService     executionsApplicationNodeExecutions.NodeExecutionService
	WorkflowService          workflowsApplicationWorkflows.WorkflowService
	OrganizationUserService  organizationsApplicationOrganizationUsers.OrganizationUserService
	UserService              accountApplicationUsers.UserService
	LinkedAccountService     accountApplicationLinkedAccounts.LinkedAccountService
}

func RegisterInfrastructure(deps RegisterInfrastructureDeps) *Handlers {
	workflowResolver := workflowresolver.New(deps.WorkflowService)

	return &Handlers{
		GetAllTaskExecutions:     cqrs.ValidationBehavior(getalltaskexecutions.New(deps.TaskExecutionRepository, workflowResolver)),
		GetTaskExecutionByID:     cqrs.ValidationBehavior(gettaskexecutionbyid.New(deps.TaskExecutionRepository, deps.NodeExecutionService, workflowResolver, deps.OrganizationUserService, deps.UserService, deps.LinkedAccountService)),
		DeleteTaskExecution:      cqrs.ValidationBehavior(deletetaskexecution.New(deps.TaskExecutionRepository, deps.TransactionManager)),
		BulkDeleteTaskExecutions: cqrs.ValidationBehavior(bulkdeletetaskexecutions.New(deps.TaskExecutionRepository, deps.TransactionManager)),
		GetAllToolInvocations:    cqrs.ValidationBehavior(getalltoolinvocations.New(deps.ToolInvocationRepository)),
		GetToolInvocationByID:    cqrs.ValidationBehavior(gettoolinvocationbyid.New(deps.ToolInvocationRepository, deps.APIKeyService)),
	}
}
