package getallworkflows

import (
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
	organizationsDomainOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
	workflowsDomainWorkflows "github.com/blocknextai/platform-api/internal/workflows/domain/workflows"
	"github.com/google/uuid"
)

func MapGetAllWorkflowsQueryToGetAllWorkflowsResponse(
	workflows []*workflowsDomainWorkflows.Workflow,
	organizationUsersByID map[uuid.UUID]*organizationsDomainOrganizationUsers.OrganizationUser,
	usersByID map[uuid.UUID]*accountDomainUsers.User,
	linkedAccountsByUserID map[uuid.UUID][]*accountDomainLinkedAccounts.LinkedAccount,
) []*WorkflowResponse {
	items := make([]*WorkflowResponse, 0, len(workflows))
	for _, workflow := range workflows {
		item := &WorkflowResponse{
			ID:             workflow.ID,
			OrganizationID: workflow.OrganizationID,
			Owner:          buildOwner(workflow.OwnerID, organizationUsersByID, usersByID, linkedAccountsByUserID),
			Title:          workflow.Title,
			Description:    workflow.Description,
			IsPinned:       workflow.IsPinned,
		}

		items = append(items, item)
	}
	return items
}

func buildOwner(
	ownerID uuid.UUID,
	organizationUsersByID map[uuid.UUID]*organizationsDomainOrganizationUsers.OrganizationUser,
	usersByID map[uuid.UUID]*accountDomainUsers.User,
	linkedAccountsByUserID map[uuid.UUID][]*accountDomainLinkedAccounts.LinkedAccount,
) *Owner {
	organizationUser, ok := organizationUsersByID[ownerID]
	if !ok {
		return nil
	}

	owner := &Owner{
		ID:    organizationUser.ID,
		Alias: organizationUser.Alias,
	}

	if user, ok := usersByID[organizationUser.UserID]; ok {
		owner.IsVerified = user.IsVerified
	}

	if linkedAccounts, ok := linkedAccountsByUserID[organizationUser.UserID]; ok {
		owner.LinkedAccounts = mapLinkedAccounts(linkedAccounts)
	}

	return owner
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
