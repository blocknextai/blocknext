package organizations

import (
	"context"
	"time"

	"github.com/google/uuid"

	organizationsDomainOrganizations "github.com/blocknextai/platform-api/internal/modules/organizations/internal/domain/organizations"
)

type GetOrganizationByIDQuery struct {
	UserID         uuid.UUID
	OrganizationID uuid.UUID
}

type GetOrganizationByIDResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	IsVerified  bool      `json:"isVerified"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (s *Service) GetOrganizationByID(ctx context.Context, request *GetOrganizationByIDQuery) (*GetOrganizationByIDResponse, error) {
	if _, err := s.organizationUserRepository.GetByOrganizationIDAndUserID(ctx, request.OrganizationID, request.UserID); err != nil {
		return nil, organizationsDomainOrganizations.ErrOrganizationNotFound
	}

	organization, err := s.organizationRepository.GetByID(ctx, request.OrganizationID)
	if err != nil {
		return nil, err
	}

	return MapGetOrganizationByIDQueryToGetOrganizationByIDResponse(organization), nil
}

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
