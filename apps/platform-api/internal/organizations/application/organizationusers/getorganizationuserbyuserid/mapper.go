package getorganizationuserbyuserid

import (
	"github.com/blocknextai/go-packages/rbac"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
	"github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
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

func MapGetOrganizationUserByUserIDQueryToGetOrganizationUserByUserIDResponse(
	organizationUser *organizationusers.OrganizationUser,
	user *accountDomainUsers.User,
	linkedAccounts []*accountDomainLinkedAccounts.LinkedAccount,
) *GetOrganizationUserByUserIDResponse {
	permissions := rbac.OrganizationPermissions(organizationUser.Role)

	isVerified := false
	if user != nil {
		isVerified = user.IsVerified
	}

	linkedAccountsResponse := []LinkedAccount{}
	if linkedAccounts != nil {
		linkedAccountsResponse = mapDomainLinkedAccountsToResponse(linkedAccounts)
	}

	return &GetOrganizationUserByUserIDResponse{
		ID:             organizationUser.ID,
		OrganizationID: organizationUser.OrganizationID,
		UserID:         organizationUser.UserID,
		Role:           organizationUser.Role,
		Permissions:    permissions,
		Alias:          organizationUser.Alias,
		IsVerified:     isVerified,
		LinkedAccounts: linkedAccountsResponse,
	}
}
