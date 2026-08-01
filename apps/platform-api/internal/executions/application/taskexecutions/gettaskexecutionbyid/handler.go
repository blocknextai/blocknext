package gettaskexecutionbyid

import (
	"context"
	"log/slog"

	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/account/application/linkedaccounts"
	accountApplicationUsers "github.com/blocknextai/platform-api/internal/account/application/users"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	executionsApplicationNodeExecutions "github.com/blocknextai/platform-api/internal/executions/application/nodeexecutions"
	executionsDomainNodeExecutions "github.com/blocknextai/platform-api/internal/executions/domain/nodeexecutions"
	"github.com/blocknextai/platform-api/internal/executions/domain/taskexecutions"
	organizationsApplicationOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/application/organizationusers"
	organizationsDomainOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
	workflowsApplicationWorkflows "github.com/blocknextai/platform-api/internal/workflows/application/workflows"
	"github.com/google/uuid"
)

type Handler struct {
	taskExecutionRepository taskexecutions.TaskExecutionRepository
	nodeExecutionService    executionsApplicationNodeExecutions.NodeExecutionService
	workflowService         workflowsApplicationWorkflows.WorkflowService
	organizationUserService organizationsApplicationOrganizationUsers.OrganizationUserService
	userService             accountApplicationUsers.UserService
	linkedAccountService    accountApplicationLinkedAccounts.LinkedAccountService
}

func New(
	taskExecutionRepository taskexecutions.TaskExecutionRepository,
	nodeExecutionService executionsApplicationNodeExecutions.NodeExecutionService,
	workflowService workflowsApplicationWorkflows.WorkflowService,
	organizationUserService organizationsApplicationOrganizationUsers.OrganizationUserService,
	userService accountApplicationUsers.UserService,
	linkedAccountService accountApplicationLinkedAccounts.LinkedAccountService,
) *Handler {
	return &Handler{
		taskExecutionRepository: taskExecutionRepository,
		nodeExecutionService:    nodeExecutionService,
		workflowService:         workflowService,
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

	workflow := h.resolveWorkflow(ctx, taskExecution.ExecutionContext, taskExecution.ContextItemID, taskExecution.OrganizationID)

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

func (h *Handler) resolveWorkflow(
	ctx context.Context,
	executionContext commonDomain.ExecutionContext,
	contextItemID uuid.UUID,
	organizationID uuid.UUID,
) Workflow {
	switch executionContext {
	case commonDomain.ExecutionContextWorkflow:
		workflow, err := h.workflowService.GetWorkflow(ctx, organizationID, contextItemID)
		if err != nil {
			slog.WarnContext(ctx, "Failed to resolve workflow for task execution",
				"component", "gettaskexecutionbyid",
				"organization_id", organizationID,
				"context_item_id", contextItemID,
				"error", err)
			return Workflow{Title: "[deleted]"}
		}
		return Workflow{
			ID:    workflow.ID,
			Title: workflow.Title,
		}
	default:
		return Workflow{Title: "[deleted]"}
	}
}
