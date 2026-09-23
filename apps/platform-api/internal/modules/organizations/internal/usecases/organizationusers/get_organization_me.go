package organizationusers

import (
	"context"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/rbac"
	organizationsDomainOrganizationusers "github.com/blocknextai/platform-api/internal/modules/organizations/internal/domain/organizationusers"
)

type GetOrganizationMeQuery struct {
	OrganizationID uuid.UUID
	UserID         uuid.UUID
}

type GetOrganizationMeResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organizationId"`
	UserID         uuid.UUID `json:"userId"`
	Alias          string    `json:"alias"`
	Role           string    `json:"role"`
	Permissions    []string  `json:"permissions"`
}

func (s *Service) GetOrganizationMe(ctx context.Context, request *GetOrganizationMeQuery) (*GetOrganizationMeResponse, error) {
	organizationUser, err := s.organizationUserRepository.GetByOrganizationIDAndUserID(ctx, request.OrganizationID, request.UserID)
	if err != nil {
		return nil, err
	}

	permissions := rbac.OrganizationPermissions(organizationUser.Role)

	return MapOrganizationUserToResponse(organizationUser, permissions), nil
}

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
