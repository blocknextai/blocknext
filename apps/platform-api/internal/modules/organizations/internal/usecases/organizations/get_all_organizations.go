package organizations

import (
	"context"
	"time"

	"github.com/google/uuid"

	organizationsDomainOrganizations "github.com/blocknextai/platform-api/internal/modules/organizations/internal/domain/organizations"
)

type GetAllOrganizationsQuery struct {
	UserID uuid.UUID
}

type OrganizationResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	IsVerified  bool      `json:"isVerified"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type GetAllOrganizationsResponse = []OrganizationResponse

func (s *Service) GetAllOrganizations(ctx context.Context, request *GetAllOrganizationsQuery) (*GetAllOrganizationsResponse, error) {
	organizations, err := s.organizationRepository.GetAllByUserID(ctx, request.UserID)
	if err != nil {
		return nil, err
	}

	return MapGetAllOrganizationsQueryToGetAllOrganizationsResponse(organizations), nil
}

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
