package organizationusers

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/rbac"
	accountContract "github.com/blocknextai/platform-api/internal/modules/account/contract"
	"github.com/blocknextai/platform-api/internal/modules/organizations/internal/domain/organizationusers"
)

type GetOrganizationUserByUserIDQuery struct {
	OrganizationID uuid.UUID
	UserID         uuid.UUID
}

type GetOrganizationUserByUserIDLinkedAccount struct {
	AuthProvider string  `json:"authProvider"`
	DisplayName  *string `json:"displayName"`
	IsPrimary    bool    `json:"isPrimary"`
}

type GetOrganizationUserByUserIDResponse struct {
	ID             uuid.UUID                                  `json:"id"`
	OrganizationID uuid.UUID                                  `json:"organizationId"`
	UserID         uuid.UUID                                  `json:"userId"`
	Role           string                                     `json:"role"`
	Permissions    []string                                   `json:"permissions"`
	Alias          string                                     `json:"alias"`
	IsVerified     bool                                       `json:"isVerified"`
	LinkedAccounts []GetOrganizationUserByUserIDLinkedAccount `json:"linkedAccounts"`
}

func (s *Service) GetOrganizationUserByUserID(ctx context.Context, request *GetOrganizationUserByUserIDQuery) (*GetOrganizationUserByUserIDResponse, error) {
	organizationUser, err := s.organizationUserRepository.GetByOrganizationIDAndUserID(ctx, request.OrganizationID, request.UserID)
	if err != nil {
		return nil, err
	}

	var user *accountContract.User
	var linkedAccounts []*accountContract.LinkedAccount

	user, err = s.userService.GetByID(ctx, organizationUser.UserID)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get user for organization user",
			"component", "Handler",
			"user_id", organizationUser.UserID,
			"error", err)
	}

	linkedAccounts, err = s.linkedAccountService.GetAllByUserID(ctx, organizationUser.UserID)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get linked accounts for organization user",
			"component", "Handler",
			"user_id", organizationUser.UserID,
			"error", err)
	}

	return MapGetOrganizationUserByUserIDQueryToGetOrganizationUserByUserIDResponse(organizationUser, user, linkedAccounts), nil
}

func getOrganizationUserByUserIDMapDomainLinkedAccountsToResponse(linkedAccounts []*accountContract.LinkedAccount) []GetOrganizationUserByUserIDLinkedAccount {
	result := make([]GetOrganizationUserByUserIDLinkedAccount, 0, len(linkedAccounts))
	for _, linkedAccount := range linkedAccounts {
		result = append(result, GetOrganizationUserByUserIDLinkedAccount{
			AuthProvider: linkedAccount.AuthProvider.String(),
			DisplayName:  linkedAccount.DisplayName,
			IsPrimary:    linkedAccount.IsPrimary,
		})
	}
	return result
}

func MapGetOrganizationUserByUserIDQueryToGetOrganizationUserByUserIDResponse(
	organizationUser *organizationusers.OrganizationUser,
	user *accountContract.User,
	linkedAccounts []*accountContract.LinkedAccount,
) *GetOrganizationUserByUserIDResponse {
	permissions := rbac.OrganizationPermissions(organizationUser.Role)

	isVerified := false
	if user != nil {
		isVerified = user.IsVerified
	}

	linkedAccountsResponse := []GetOrganizationUserByUserIDLinkedAccount{}
	if linkedAccounts != nil {
		linkedAccountsResponse = getOrganizationUserByUserIDMapDomainLinkedAccountsToResponse(linkedAccounts)
	}

	return &GetOrganizationUserByUserIDResponse{
		ID:             organizationUser.ID,
		OrganizationID: organizationUser.OrganizationID,
		UserID:         organizationUser.UserID,
		Role:           organizationUser.Role,
		Permissions:    permissions,
		Alias:          organizationUser.Alias,
		IsVerified:     isVerified,
		LinkedAccounts: linkedAccountsResponse,
	}
}
