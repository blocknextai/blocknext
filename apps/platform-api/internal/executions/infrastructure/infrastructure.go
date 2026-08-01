package infrastructure

import (
	"github.com/blocknextai/go-packages/database"
	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/account/application/linkedaccounts"
	accountApplicationUsers "github.com/blocknextai/platform-api/internal/account/application/users"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	executionsApplicationNodeExecutions "github.com/blocknextai/platform-api/internal/executions/application/nodeexecutions"
	"github.com/blocknextai/platform-api/internal/executions/application/taskexecutions/bulkdeletetaskexecutions"
	"github.com/blocknextai/platform-api/internal/executions/application/taskexecutions/deletetaskexecution"
	"github.com/blocknextai/platform-api/internal/executions/application/taskexecutions/getalltaskexecutions"
	"github.com/blocknextai/platform-api/internal/executions/application/taskexecutions/gettaskexecutionbyid"
	"github.com/blocknextai/platform-api/internal/executions/domain/taskexecutions"
	organizationsApplicationOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/application/organizationusers"
	workflowsApplicationWorkflows "github.com/blocknextai/platform-api/internal/workflows/application/workflows"
)

type Handlers struct {
	GetAllTaskExecutions     cqrs.Handler[*getalltaskexecutions.GetAllTaskExecutionsQuery, *getalltaskexecutions.GetAllTaskExecutionsResponse]
	GetTaskExecutionByID     cqrs.Handler[*gettaskexecutionbyid.GetTaskExecutionByIDQuery, *gettaskexecutionbyid.GetTaskExecutionByIDResponse]
	DeleteTaskExecution      cqrs.Handler[*deletetaskexecution.DeleteTaskExecutionCommand, *deletetaskexecution.DeleteTaskExecutionResponse]
	BulkDeleteTaskExecutions cqrs.Handler[*bulkdeletetaskexecutions.BulkDeleteTaskExecutionsCommand, *bulkdeletetaskexecutions.BulkDeleteTaskExecutionsResponse]
}

type RegisterInfrastructureDeps struct {
	TransactionManager database.TransactionManager

	TaskExecutionRepository taskexecutions.TaskExecutionRepository
	NodeExecutionService    executionsApplicationNodeExecutions.NodeExecutionService
	WorkflowService         workflowsApplicationWorkflows.WorkflowService
	OrganizationUserService organizationsApplicationOrganizationUsers.OrganizationUserService
	UserService             accountApplicationUsers.UserService
	LinkedAccountService    accountApplicationLinkedAccounts.LinkedAccountService
}

func RegisterInfrastructure(deps RegisterInfrastructureDeps) *Handlers {
	return &Handlers{
		GetAllTaskExecutions:     cqrs.ValidationBehavior(getalltaskexecutions.New(deps.TaskExecutionRepository, deps.WorkflowService)),
		GetTaskExecutionByID:     cqrs.ValidationBehavior(gettaskexecutionbyid.New(deps.TaskExecutionRepository, deps.NodeExecutionService, deps.WorkflowService, deps.OrganizationUserService, deps.UserService, deps.LinkedAccountService)),
		DeleteTaskExecution:      cqrs.ValidationBehavior(deletetaskexecution.New(deps.TaskExecutionRepository, deps.TransactionManager)),
		BulkDeleteTaskExecutions: cqrs.ValidationBehavior(bulkdeletetaskexecutions.New(deps.TaskExecutionRepository, deps.TransactionManager)),
	}
}
