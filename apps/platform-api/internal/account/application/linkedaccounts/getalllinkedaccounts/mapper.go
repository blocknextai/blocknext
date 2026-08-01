package getalllinkedaccounts

import (
	accountDomainLinkedaccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
)

func MapLinkedAccountsToResponse(linkedAccounts []*accountDomainLinkedaccounts.LinkedAccount) *GetAllLinkedAccountsResponse {
	result := make(GetAllLinkedAccountsResponse, 0, len(linkedAccounts))
	for _, linkedAccount := range linkedAccounts {
		result = append(result, LinkedAccount{
			ID:           linkedAccount.ID,
			AuthProvider: linkedAccount.AuthProvider.String(),
			Identifier:   linkedAccount.Identifier,
			DisplayName:  linkedAccount.DisplayName,
			IsPrimary:    linkedAccount.IsPrimary,
			IsVerified:   linkedAccount.IsVerified,
		})
	}
	return &result
}
