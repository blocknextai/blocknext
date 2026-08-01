package getorganizationme

import (
	organizationsDomainOrganizationusers "github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
)

func MapOrganizationUserToResponse(organizationUser *organizationsDomainOrganizationusers.OrganizationUser, permissions []string) *GetOrganizationMeResponse {
	return &GetOrganizationMeResponse{
		ID:             organizationUser.ID,
		OrganizationID: organizationUser.OrganizationID,
		UserID:         organizationUser.UserID,
		Alias:          organizationUser.Alias,
		Role:           organizationUser.Role,
		Permissions:    permissions,
	}
}
