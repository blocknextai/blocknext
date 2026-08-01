package gettaskexecutionbyid

import (
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
	executionsDomainNodeExecutions "github.com/blocknextai/platform-api/internal/executions/domain/nodeexecutions"
	executionsDomainTaskexecutions "github.com/blocknextai/platform-api/internal/executions/domain/taskexecutions"
	organizationsDomainOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
)

func MapTaskExecutionToResponse(
	taskExecution *executionsDomainTaskexecutions.TaskExecution,
	workflow Workflow,
	organizationUser *organizationsDomainOrganizationUsers.OrganizationUser,
	user *accountDomainUsers.User,
	linkedAccounts []*accountDomainLinkedAccounts.LinkedAccount,
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
	organizationUser *organizationsDomainOrganizationUsers.OrganizationUser,
	user *accountDomainUsers.User,
	linkedAccounts []*accountDomainLinkedAccounts.LinkedAccount,
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

func mapLinkedAccounts(linkedAccounts []*accountDomainLinkedAccounts.LinkedAccount) []LinkedAccount {
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
