package getallorganizations

import (
	"context"

	organizationsDomainOrganizations "github.com/blocknextai/platform-api/internal/organizations/domain/organizations"
)

type Handler struct {
	organizationRepository organizationsDomainOrganizations.OrganizationRepository
}

func New(
	organizationRepository organizationsDomainOrganizations.OrganizationRepository,
) *Handler {
	return &Handler{
		organizationRepository: organizationRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetAllOrganizationsQuery) (*GetAllOrganizationsResponse, error) {
	organizations, err := h.organizationRepository.GetAllByUserID(ctx, request.UserID)
	if err != nil {
		return nil, err
	}

	return MapGetAllOrganizationsQueryToGetAllOrganizationsResponse(organizations), nil
}
