package getallorganizations

import (
	organizationsDomainOrganizations "github.com/blocknextai/platform-api/internal/organizations/domain/organizations"
)

func MapGetAllOrganizationsQueryToGetAllOrganizationsResponse(
	organizations []*organizationsDomainOrganizations.Organization,
) *GetAllOrganizationsResponse {
	response := make(GetAllOrganizationsResponse, 0, len(organizations))
	for _, organization := range organizations {
		response = append(response, OrganizationResponse{
			ID:          organization.ID,
			Title:       organization.Title,
			Description: organization.Description,
			IsVerified:  organization.IsVerified,
			CreatedAt:   organization.CreatedAt,
			UpdatedAt:   organization.UpdatedAt,
		})
	}
	return &response
}
