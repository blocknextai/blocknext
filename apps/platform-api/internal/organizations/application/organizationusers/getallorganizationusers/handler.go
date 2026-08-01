package getallorganizationusers

import (
	"context"
	"log/slog"

	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/account/application/linkedaccounts"
	accountApplicationUsers "github.com/blocknextai/platform-api/internal/account/application/users"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
	organizationsDomainOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
	"github.com/google/uuid"
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

func (h *Handler) Handle(ctx context.Context, request *GetAllOrganizationUsersQuery) (*GetAllOrganizationUsersResponse, error) {
	organizationUsers, totalCount, err := h.organizationUserRepository.GetAllByOrganizationID(
		ctx,
		request.OrganizationID,
		request.Search.Query,
		request.Pagination.Offset,
		request.Pagination.Limit,
	)
	if err != nil {
		return nil, err
	}

	userIDs := make([]uuid.UUID, 0, len(organizationUsers))
	for _, organizationUser := range organizationUsers {
		userIDs = append(userIDs, organizationUser.UserID)
	}

	users, err := h.userService.GetAllByIDs(ctx, userIDs)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get users for organization users",
			"component", "Handler",
			"organization_id", request.OrganizationID,
			"error", err)
		users = []*accountDomainUsers.User{}
	}

	usersByID := make(map[uuid.UUID]*accountDomainUsers.User, len(users))
	for _, user := range users {
		usersByID[user.ID] = user
	}

	linkedAccounts, err := h.linkedAccountService.GetAllByUserIDs(ctx, userIDs)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get linked accounts for organization users",
			"component", "Handler",
			"organization_id", request.OrganizationID,
			"error", err)
		linkedAccounts = []*accountDomainLinkedAccounts.LinkedAccount{}
	}

	linkedAccountsByUserID := make(map[uuid.UUID][]*accountDomainLinkedAccounts.LinkedAccount)
	for _, linkedAccount := range linkedAccounts {
		linkedAccountsByUserID[linkedAccount.UserID] = append(linkedAccountsByUserID[linkedAccount.UserID], linkedAccount)
	}

	return &GetAllOrganizationUsersResponse{
		Items:      MapGetAllOrganizationUsersQueryToGetAllOrganizationUsersResponse(organizationUsers, usersByID, linkedAccountsByUserID),
		TotalCount: totalCount,
	}, nil
}
