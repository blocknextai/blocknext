package taskexecutions

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	accountContract "github.com/blocknextai/platform-api/internal/modules/account/contract"
	executionsDomainNodeExecutions "github.com/blocknextai/platform-api/internal/modules/executions/internal/domain/nodeexecutions"
	executionsDomainTaskexecutions "github.com/blocknextai/platform-api/internal/modules/executions/internal/domain/taskexecutions"
	organizationsContract "github.com/blocknextai/platform-api/internal/modules/organizations/contract"
)

type GetTaskExecutionByIDQuery struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
}

type GetTaskExecutionByIDResponse struct {
	ID              uuid.UUID                    `json:"id,omitempty"`
	Workflow        GetTaskExecutionByIDWorkflow `json:"workflow"`
	TriggeredByUser *TriggeredByUser             `json:"triggeredByUser,omitempty"`
	Status          string                       `json:"status,omitempty"`
	ExecutionType   string                       `json:"executionType,omitempty"`
	ErrorMessage    *string                      `json:"errorMessage,omitempty"`
	StartedAt       *time.Time                   `json:"startedAt,omitempty"`
	CompletedAt     *time.Time                   `json:"completedAt,omitempty"`
	NodeExecutions  []NodeExecution              `json:"nodeExecutions,omitempty"`
}

type GetTaskExecutionByIDWorkflow struct {
	ID    uuid.UUID `json:"id,omitempty"`
	Title string    `json:"title,omitempty"`
}

type LinkedAccount struct {
	AuthProvider string  `json:"authProvider"`
	DisplayName  *string `json:"displayName"`
	IsPrimary    bool    `json:"isPrimary"`
}

type TriggeredByUser struct {
	ID             uuid.UUID       `json:"id,omitempty"`
	Alias          string          `json:"alias,omitempty"`
	IsVerified     bool            `json:"isVerified"`
	LinkedAccounts []LinkedAccount `json:"linkedAccounts"`
}

type NodeExecution struct {
	ID           uuid.UUID        `json:"id,omitempty"`
	NodeType     string           `json:"nodeType,omitempty"`
	Status       string           `json:"status,omitempty"`
	Inputs       []map[string]any `json:"inputs,omitempty"`
	Outputs      []map[string]any `json:"outputs,omitempty"`
	ErrorMessage *string          `json:"errorMessage,omitempty"`
	StartedAt    *time.Time       `json:"startedAt,omitempty"`
	CompletedAt  *time.Time       `json:"completedAt,omitempty"`
}

func (s *Service) GetTaskExecutionByID(ctx context.Context, request *GetTaskExecutionByIDQuery) (*GetTaskExecutionByIDResponse, error) {
	taskExecution, err := s.taskExecutionRepository.GetByIDAndOrganizationID(ctx, request.ID, request.OrganizationID)
	if err != nil {
		return nil, err
	}

	resolved := s.workflowResolver.Resolve(ctx, taskExecution.ExecutionContext, taskExecution.ContextItemID, taskExecution.OrganizationID)
	workflow := GetTaskExecutionByIDWorkflow{ID: resolved.ID, Title: resolved.Title}

	var organizationUser *organizationsContract.OrganizationUser
	if taskExecution.TriggeredByUserID != nil {
		organizationUser, err = s.organizationUserService.GetByOrganizationIDAndUserID(ctx, taskExecution.OrganizationID, *taskExecution.TriggeredByUserID)
		if err != nil {
			slog.WarnContext(ctx, "Failed to get organization user for task execution",
				"component", "gettaskexecutionbyid",
				"organization_id", taskExecution.OrganizationID,
				"user_id", *taskExecution.TriggeredByUserID,
				"error", err)
		}
	}

	var user *accountContract.User
	var linkedAccounts []*accountContract.LinkedAccount
	if organizationUser != nil {
		user, err = s.userService.GetByID(ctx, organizationUser.UserID)
		if err != nil {
			slog.WarnContext(ctx, "Failed to get user for task execution",
				"component", "gettaskexecutionbyid",
				"user_id", organizationUser.UserID,
				"error", err)
		}

		linkedAccounts, err = s.linkedAccountService.GetAllByUserID(ctx, organizationUser.UserID)
		if err != nil {
			slog.WarnContext(ctx, "Failed to get linked accounts for task execution",
				"component", "gettaskexecutionbyid",
				"user_id", organizationUser.UserID,
				"error", err)
			linkedAccounts = []*accountContract.LinkedAccount{}
		}
	}

	nodeExecutions, err := s.nodeExecutionService.GetAllByTaskID(ctx, taskExecution.ID)
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

func MapTaskExecutionToResponse(
	taskExecution *executionsDomainTaskexecutions.TaskExecution,
	workflow GetTaskExecutionByIDWorkflow,
	organizationUser *organizationsContract.OrganizationUser,
	user *accountContract.User,
	linkedAccounts []*accountContract.LinkedAccount,
	nodeExecutions []*executionsDomainNodeExecutions.NodeExecution,
) *GetTaskExecutionByIDResponse {
	return &GetTaskExecutionByIDResponse{
		ID:              taskExecution.ID,
		Workflow:        workflow,
		TriggeredByUser: buildTriggeredByUser(organizationUser, user, linkedAccounts),
		Status:          taskExecution.Status,
		ExecutionType:   taskExecution.ExecutionType.String(),
		ErrorMessage:    taskExecution.ErrorMessage,
		StartedAt:       taskExecution.StartedAt,
		CompletedAt:     taskExecution.CompletedAt,
		NodeExecutions:  mapNodeExecutions(nodeExecutions),
	}
}

func buildTriggeredByUser(
	organizationUser *organizationsContract.OrganizationUser,
	user *accountContract.User,
	linkedAccounts []*accountContract.LinkedAccount,
) *TriggeredByUser {
	if organizationUser == nil {
		return nil
	}

	triggeredByUser := TriggeredByUser{
		ID:             organizationUser.ID,
		Alias:          organizationUser.Alias,
		LinkedAccounts: mapLinkedAccounts(linkedAccounts),
	}

	if user != nil {
		triggeredByUser.IsVerified = user.IsVerified
	}

	return &triggeredByUser
}

func mapLinkedAccounts(linkedAccounts []*accountContract.LinkedAccount) []LinkedAccount {
	result := make([]LinkedAccount, 0, len(linkedAccounts))
	for _, linkedAccount := range linkedAccounts {
		result = append(result, LinkedAccount{
			AuthProvider: linkedAccount.AuthProvider.String(),
			DisplayName:  linkedAccount.DisplayName,
			IsPrimary:    linkedAccount.IsPrimary,
		})
	}
	return result
}

func mapNodeExecutions(nodeExecutions []*executionsDomainNodeExecutions.NodeExecution) []NodeExecution {
	result := make([]NodeExecution, 0, len(nodeExecutions))
	for _, nodeExecution := range nodeExecutions {
		result = append(result, NodeExecution{
			ID:           nodeExecution.ID,
			NodeType:     nodeExecution.NodeType,
			Status:       nodeExecution.Status,
			Inputs:       nodeExecution.Inputs,
			Outputs:      nodeExecution.Outputs,
			ErrorMessage: nodeExecution.ErrorMessage,
			StartedAt:    nodeExecution.StartedAt,
			CompletedAt:  nodeExecution.CompletedAt,
		})
	}
	return result
}
