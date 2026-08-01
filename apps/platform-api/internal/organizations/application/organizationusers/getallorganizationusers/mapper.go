package getallorganizationusers

import (
	"github.com/blocknextai/go-packages/rbac"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
	"github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
	"github.com/google/uuid"
)

func mapDomainLinkedAccountsToResponse(linkedAccounts []*accountDomainLinkedAccounts.LinkedAccount) []LinkedAccount {
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

func MapGetAllOrganizationUsersQueryToGetAllOrganizationUsersResponse(
	organizationUsers []*organizationusers.OrganizationUser,
	usersByID map[uuid.UUID]*accountDomainUsers.User,
	linkedAccountsByUserID map[uuid.UUID][]*accountDomainLinkedAccounts.LinkedAccount,
) []*OrganizationUserResponse {
	organizationUsersResponse := make([]*OrganizationUserResponse, 0, len(organizationUsers))
	for _, organizationUser := range organizationUsers {
		permissions := rbac.OrganizationPermissions(organizationUser.Role)

		isVerified := false
		if user, ok := usersByID[organizationUser.UserID]; ok {
			isVerified = user.IsVerified
		}

		linkedAccountsResponse := []LinkedAccount{}
		if linkedAccounts, ok := linkedAccountsByUserID[organizationUser.UserID]; ok {
			linkedAccountsResponse = mapDomainLinkedAccountsToResponse(linkedAccounts)
		}

		organizationUsersResponse = append(organizationUsersResponse, &OrganizationUserResponse{
			ID:             organizationUser.ID,
			OrganizationID: organizationUser.OrganizationID,
			UserID:         organizationUser.UserID,
			Role:           organizationUser.Role,
			Permissions:    permissions,
			Alias:          organizationUser.Alias,
			IsVerified:     isVerified,
			LinkedAccounts: linkedAccountsResponse,
		})
	}

	return organizationUsersResponse
}
