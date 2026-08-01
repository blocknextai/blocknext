package getorganizationme

import (
	"context"

	"github.com/blocknextai/go-packages/rbac"
	"github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
)

type Handler struct {
	organizationUserRepository organizationusers.OrganizationUserRepository
}

func New(
	organizationUserRepository organizationusers.OrganizationUserRepository,
) *Handler {
	return &Handler{
		organizationUserRepository: organizationUserRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetOrganizationMeQuery) (*GetOrganizationMeResponse, error) {
	organizationUser, err := h.organizationUserRepository.GetByOrganizationIDAndUserID(ctx, request.OrganizationID, request.UserID)
	if err != nil {
		return nil, err
	}

	permissions := rbac.OrganizationPermissions(organizationUser.Role)

	return MapOrganizationUserToResponse(organizationUser, permissions), nil
}
