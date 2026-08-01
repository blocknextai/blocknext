package getorganizationbyid

import (
	"context"

	organizationsDomainOrganizations "github.com/blocknextai/platform-api/internal/organizations/domain/organizations"
	organizationsDomainOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
)

type Handler struct {
	organizationRepository     organizationsDomainOrganizations.OrganizationRepository
	organizationUserRepository organizationsDomainOrganizationUsers.OrganizationUserRepository
}

func New(
	organizationRepository organizationsDomainOrganizations.OrganizationRepository,
	organizationUserRepository organizationsDomainOrganizationUsers.OrganizationUserRepository,
) *Handler {
	return &Handler{
		organizationRepository:     organizationRepository,
		organizationUserRepository: organizationUserRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetOrganizationByIDQuery) (*GetOrganizationByIDResponse, error) {
	if _, err := h.organizationUserRepository.GetByOrganizationIDAndUserID(ctx, request.OrganizationID, request.UserID); err != nil {
		return nil, organizationsDomainOrganizations.ErrOrganizationNotFound
	}

	organization, err := h.organizationRepository.GetByID(ctx, request.OrganizationID)
	if err != nil {
		return nil, err
	}

	return MapGetOrganizationByIDQueryToGetOrganizationByIDResponse(organization), nil
}
