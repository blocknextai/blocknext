package getorganizationbyid

import (
	organizationsDomainOrganizations "github.com/blocknextai/platform-api/internal/organizations/domain/organizations"
)

func MapGetOrganizationByIDQueryToGetOrganizationByIDResponse(
	organization *organizationsDomainOrganizations.Organization,
) *GetOrganizationByIDResponse {
	return &GetOrganizationByIDResponse{
		ID:          organization.ID,
		Title:       organization.Title,
		Description: organization.Description,
		IsVerified:  organization.IsVerified,
		CreatedAt:   organization.CreatedAt,
		UpdatedAt:   organization.UpdatedAt,
	}
}
