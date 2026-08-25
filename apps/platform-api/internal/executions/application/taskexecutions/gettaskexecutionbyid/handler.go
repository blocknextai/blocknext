package gettaskexecutionbyid

import (
	"context"
	"log/slog"

	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/account/application/linkedaccounts"
	accountApplicationUsers "github.com/blocknextai/platform-api/internal/account/application/users"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
	executionsApplicationNodeExecutions "github.com/blocknextai/platform-api/internal/executions/application/nodeexecutions"
	"github.com/blocknextai/platform-api/internal/executions/application/taskexecutions/workflowresolver"
	executionsDomainNodeExecutions "github.com/blocknextai/platform-api/internal/executions/domain/nodeexecutions"
	"github.com/blocknextai/platform-api/internal/executions/domain/taskexecutions"
	organizationsApplicationOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/application/organizationusers"
	organizationsDomainOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
)

type Handler struct {
	taskExecutionRepository taskexecutions.TaskExecutionRepository
	nodeExecutionService    executionsApplicationNodeExecutions.NodeExecutionService
	workflowResolver        workflowresolver.Resolver
	organizationUserService organizationsApplicationOrganizationUsers.OrganizationUserService
	userService             accountApplicationUsers.UserService
	linkedAccountService    accountApplicationLinkedAccounts.LinkedAccountService
}

func New(
	taskExecutionRepository taskexecutions.TaskExecutionRepository,
	nodeExecutionService executionsApplicationNodeExecutions.NodeExecutionService,
	workflowResolver workflowresolver.Resolver,
	organizationUserService organizationsApplicationOrganizationUsers.OrganizationUserService,
	userService accountApplicationUsers.UserService,
	linkedAccountService accountApplicationLinkedAccounts.LinkedAccountService,
) *Handler {
	return &Handler{
		taskExecutionRepository: taskExecutionRepository,
		nodeExecutionService:    nodeExecutionService,
		workflowResolver:        workflowResolver,
		organizationUserService: organizationUserService,
		userService:             userService,
		linkedAccountService:    linkedAccountService,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetTaskExecutionByIDQuery) (*GetTaskExecutionByIDResponse, error) {
	taskExecution, err := h.taskExecutionRepository.GetByIDAndOrganizationID(ctx, request.ID, request.OrganizationID)
	if err != nil {
		return nil, err
	}

	resolved := h.workflowResolver.Resolve(ctx, taskExecution.ExecutionContext, taskExecution.ContextItemID, taskExecution.OrganizationID)
	workflow := Workflow{ID: resolved.ID, Title: resolved.Title}

	var organizationUser *organizationsDomainOrganizationUsers.OrganizationUser
	if taskExecution.TriggeredByUserID != nil {
		organizationUser, err = h.organizationUserService.GetByOrganizationIDAndUserID(ctx, taskExecution.OrganizationID, *taskExecution.TriggeredByUserID)
		if err != nil {
			slog.WarnContext(ctx, "Failed to get organization user for task execution",
				"component", "gettaskexecutionbyid",
				"organization_id", taskExecution.OrganizationID,
				"user_id", *taskExecution.TriggeredByUserID,
				"error", err)
		}
	}

	var user *accountDomainUsers.User
	var linkedAccounts []*accountDomainLinkedAccounts.LinkedAccount
	if organizationUser != nil {
		user, err = h.userService.GetByID(ctx, organizationUser.UserID)
		if err != nil {
			slog.WarnContext(ctx, "Failed to get user for task execution",
				"component", "gettaskexecutionbyid",
				"user_id", organizationUser.UserID,
				"error", err)
		}

		linkedAccounts, err = h.linkedAccountService.GetAllByUserID(ctx, organizationUser.UserID)
		if err != nil {
			slog.WarnContext(ctx, "Failed to get linked accounts for task execution",
				"component", "gettaskexecutionbyid",
				"user_id", organizationUser.UserID,
				"error", err)
			linkedAccounts = []*accountDomainLinkedAccounts.LinkedAccount{}
		}
	}

	nodeExecutions, err := h.nodeExecutionService.GetAllByTaskID(ctx, taskExecution.ID)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get node executions for task execution",
			"component", "gettaskexecutionbyid",
			"task_id", taskExecution.ID,
			"error", err)
		nodeExecutions = []*executionsDomainNodeExecutions.NodeExecution{}
	}

	return MapTaskExecutionToResponse(
		taskExecution,
		workflow,
		organizationUser,
		user,
		linkedAccounts,
		nodeExecutions,
	), nil
}
