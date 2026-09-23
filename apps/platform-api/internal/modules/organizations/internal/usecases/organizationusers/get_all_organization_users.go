package organizationusers

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/rbac"
	resultPkg "github.com/blocknextai/go-packages/result"
	accountContract "github.com/blocknextai/platform-api/internal/modules/account/contract"
	"github.com/blocknextai/platform-api/internal/modules/organizations/internal/domain/organizationusers"
)

type GetAllOrganizationUsersQuery struct {
	OrganizationID uuid.UUID
	Search         resultPkg.SearchRequest
	Pagination     resultPkg.PaginationRequest
}

type LinkedAccount struct {
	AuthProvider string  `json:"authProvider"`
	DisplayName  *string `json:"displayName"`
	IsPrimary    bool    `json:"isPrimary"`
}

type OrganizationUserResponse struct {
	ID             uuid.UUID       `json:"id"`
	OrganizationID uuid.UUID       `json:"organizationId"`
	UserID         uuid.UUID       `json:"userId"`
	Role           string          `json:"role"`
	Permissions    []string        `json:"permissions"`
	Alias          string          `json:"alias"`
	IsVerified     bool            `json:"isVerified"`
	LinkedAccounts []LinkedAccount `json:"linkedAccounts"`
}

type GetAllOrganizationUsersResponse struct {
	Items      []*OrganizationUserResponse
	TotalCount int64
}

func (s *Service) GetAllOrganizationUsers(ctx context.Context, request *GetAllOrganizationUsersQuery) (*GetAllOrganizationUsersResponse, error) {
	organizationUsers, totalCount, err := s.organizationUserRepository.GetAllByOrganizationID(
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

	users, err := s.userService.GetAllByIDs(ctx, userIDs)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get users for organization users",
			"component", "Handler",
			"organization_id", request.OrganizationID,
			"error", err)
		users = []*accountContract.User{}
	}

	usersByID := make(map[uuid.UUID]*accountContract.User, len(users))
	for _, user := range users {
		usersByID[user.ID] = user
	}

	linkedAccounts, err := s.linkedAccountService.GetAllByUserIDs(ctx, userIDs)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get linked accounts for organization users",
			"component", "Handler",
			"organization_id", request.OrganizationID,
			"error", err)
		linkedAccounts = []*accountContract.LinkedAccount{}
	}

	linkedAccountsByUserID := make(map[uuid.UUID][]*accountContract.LinkedAccount)
	for _, linkedAccount := range linkedAccounts {
		linkedAccountsByUserID[linkedAccount.UserID] = append(linkedAccountsByUserID[linkedAccount.UserID], linkedAccount)
	}

	return &GetAllOrganizationUsersResponse{
		Items:      MapGetAllOrganizationUsersQueryToGetAllOrganizationUsersResponse(organizationUsers, usersByID, linkedAccountsByUserID),
		TotalCount: totalCount,
	}, nil
}

func mapDomainLinkedAccountsToResponse(linkedAccounts []*accountContract.LinkedAccount) []LinkedAccount {
	result := make([]LinkedAccount, 0, len(linkedAccounts))
	for _, linkedAccount := range linkedAccounts {
		result = append(result, LinkedAccount{
			AuthProvider: linkedAccount.AuthProvider.String(),
			DisplayName:  linkedAccount.DisplayName,
			IsPrimary:    linkedAccount.IsPrimary,
		})
	}
	return result
}

func MapGetAllOrganizationUsersQueryToGetAllOrganizationUsersResponse(
	organizationUsers []*organizationusers.OrganizationUser,
	usersByID map[uuid.UUID]*accountContract.User,
	linkedAccountsByUserID map[uuid.UUID][]*accountContract.LinkedAccount,
) []*OrganizationUserResponse {
	organizationUsersResponse := make([]*OrganizationUserResponse, 0, len(organizationUsers))
	for _, organizationUser := range organizationUsers {
		permissions := rbac.OrganizationPermissions(organizationUser.Role)

		isVerified := false
		if user, ok := usersByID[organizationUser.UserID]; ok {
			isVerified = user.IsVerified
		}

		linkedAccountsResponse := []LinkedAccount{}
		if linkedAccounts, ok := linkedAccountsByUserID[organizationUser.UserID]; ok {
			linkedAccountsResponse = mapDomainLinkedAccountsToResponse(linkedAccounts)
		}

		organizationUsersResponse = append(organizationUsersResponse, &OrganizationUserResponse{
			ID:             organizationUser.ID,
			OrganizationID: organizationUser.OrganizationID,
			UserID:         organizationUser.UserID,
			Role:           organizationUser.Role,
			Permissions:    permissions,
			Alias:          organizationUser.Alias,
			IsVerified:     isVerified,
			LinkedAccounts: linkedAccountsResponse,
		})
	}

	return organizationUsersResponse
}
