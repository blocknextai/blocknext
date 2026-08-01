package getorganizationuserbyuserid

import (
	"context"
	"log/slog"

	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/account/application/linkedaccounts"
	accountApplicationUsers "github.com/blocknextai/platform-api/internal/account/application/users"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
	organizationsDomainOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
)

type Handler struct {
	organizationUserRepository organizationsDomainOrganizationUsers.OrganizationUserRepository
	userService                accountApplicationUsers.UserService
	linkedAccountService       accountApplicationLinkedAccounts.LinkedAccountService
}

func New(
	organizationUserRepository organizationsDomainOrganizationUsers.OrganizationUserRepository,
	userService accountApplicationUsers.UserService,
	linkedAccountService accountApplicationLinkedAccounts.LinkedAccountService,
) *Handler {
	return &Handler{
		organizationUserRepository: organizationUserRepository,
		userService:                userService,
		linkedAccountService:       linkedAccountService,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetOrganizationUserByUserIDQuery) (*GetOrganizationUserByUserIDResponse, error) {
	organizationUser, err := h.organizationUserRepository.GetByOrganizationIDAndUserID(ctx, request.OrganizationID, request.UserID)
	if err != nil {
		return nil, err
	}

	var user *accountDomainUsers.User
	var linkedAccounts []*accountDomainLinkedAccounts.LinkedAccount

	user, err = h.userService.GetByID(ctx, organizationUser.UserID)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get user for organization user",
			"component", "Handler",
			"user_id", organizationUser.UserID,
			"error", err)
	}

	linkedAccounts, err = h.linkedAccountService.GetAllByUserID(ctx, organizationUser.UserID)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get linked accounts for organization user",
			"component", "Handler",
			"user_id", organizationUser.UserID,
			"error", err)
	}

	return MapGetOrganizationUserByUserIDQueryToGetOrganizationUserByUserIDResponse(organizationUser, user, linkedAccounts), nil
}
