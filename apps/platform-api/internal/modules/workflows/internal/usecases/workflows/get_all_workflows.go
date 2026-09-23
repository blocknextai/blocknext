package workflows

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	accountContract "github.com/blocknextai/platform-api/internal/modules/account/contract"
	organizationsContract "github.com/blocknextai/platform-api/internal/modules/organizations/contract"
	workflowsDomainWorkflows "github.com/blocknextai/platform-api/internal/modules/workflows/internal/domain/workflows"
)

type GetAllWorkflowsQuery struct {
	OrganizationID uuid.UUID
	Search         resultPkg.SearchRequest
	Pagination     resultPkg.PaginationRequest
}

type LinkedAccount struct {
	AuthProvider string  `json:"authProvider"`
	DisplayName  *string `json:"displayName"`
	IsPrimary    bool    `json:"isPrimary"`
}

type Owner struct {
	ID             uuid.UUID       `json:"id"`
	Alias          string          `json:"alias"`
	IsVerified     bool            `json:"isVerified"`
	LinkedAccounts []LinkedAccount `json:"linkedAccounts"`
}

type WorkflowResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organizationId"`
	Owner          *Owner    `json:"owner"`
	Title          string    `json:"title"`
	Description    *string   `json:"description"`
	IsPinned       bool      `json:"isPinned"`
}

type GetAllWorkflowsResponse struct {
	Items      []*WorkflowResponse
	TotalCount int64
}

func (s *Service) GetAllWorkflows(ctx context.Context, request *GetAllWorkflowsQuery) (*GetAllWorkflowsResponse, error) {
	workflows, totalCount, err := s.workflowRepository.GetAllByOrganizationID(
		ctx,
		request.OrganizationID,
		request.Search.Query,
		request.Pagination.Offset,
		request.Pagination.Limit,
	)
	if err != nil {
		return nil, err
	}

	organizationUserIDs := make([]uuid.UUID, 0, len(workflows))
	for _, workflow := range workflows {
		organizationUserIDs = append(organizationUserIDs, workflow.OwnerID)
	}

	organizationUsers, err := s.organizationUserService.GetAllByIDs(ctx, organizationUserIDs)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get organization users for workflows",
			"component", "getallworkflows",
			"organization_id", request.OrganizationID,
			"error", err)
		organizationUsers = []*organizationsContract.OrganizationUser{}
	}

	organizationUsersByID := make(map[uuid.UUID]*organizationsContract.OrganizationUser, len(organizationUsers))
	userIDs := make([]uuid.UUID, 0, len(organizationUsers))
	for _, organizationUser := range organizationUsers {
		organizationUsersByID[organizationUser.ID] = organizationUser
		userIDs = append(userIDs, organizationUser.UserID)
	}

	users, err := s.userService.GetAllByIDs(ctx, userIDs)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get users for workflows",
			"component", "getallworkflows",
			"error", err)
		users = []*accountContract.User{}
	}

	usersByID := make(map[uuid.UUID]*accountContract.User, len(users))
	for _, user := range users {
		usersByID[user.ID] = user
	}

	linkedAccounts, err := s.linkedAccountService.GetAllByUserIDs(ctx, userIDs)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get linked accounts for workflows",
			"component", "getallworkflows",
			"error", err)
		linkedAccounts = []*accountContract.LinkedAccount{}
	}

	linkedAccountsByUserID := make(map[uuid.UUID][]*accountContract.LinkedAccount)
	for _, linkedAccount := range linkedAccounts {
		linkedAccountsByUserID[linkedAccount.UserID] = append(linkedAccountsByUserID[linkedAccount.UserID], linkedAccount)
	}

	return &GetAllWorkflowsResponse{
		Items: MapGetAllWorkflowsQueryToGetAllWorkflowsResponse(
			workflows,
			organizationUsersByID,
			usersByID,
			linkedAccountsByUserID,
		),
		TotalCount: totalCount,
	}, nil
}

func MapGetAllWorkflowsQueryToGetAllWorkflowsResponse(
	workflows []*workflowsDomainWorkflows.Workflow,
	organizationUsersByID map[uuid.UUID]*organizationsContract.OrganizationUser,
	usersByID map[uuid.UUID]*accountContract.User,
	linkedAccountsByUserID map[uuid.UUID][]*accountContract.LinkedAccount,
) []*WorkflowResponse {
	items := make([]*WorkflowResponse, 0, len(workflows))
	for _, workflow := range workflows {
		item := &WorkflowResponse{
			ID:             workflow.ID,
			OrganizationID: workflow.OrganizationID,
			Owner:          buildOwner(workflow.OwnerID, organizationUsersByID, usersByID, linkedAccountsByUserID),
			Title:          workflow.Title,
			Description:    workflow.Description,
			IsPinned:       workflow.IsPinned,
		}

		items = append(items, item)
	}
	return items
}

func buildOwner(
	ownerID uuid.UUID,
	organizationUsersByID map[uuid.UUID]*organizationsContract.OrganizationUser,
	usersByID map[uuid.UUID]*accountContract.User,
	linkedAccountsByUserID map[uuid.UUID][]*accountContract.LinkedAccount,
) *Owner {
	organizationUser, ok := organizationUsersByID[ownerID]
	if !ok {
		return nil
	}

	owner := &Owner{
		ID:    organizationUser.ID,
		Alias: organizationUser.Alias,
	}

	if user, ok := usersByID[organizationUser.UserID]; ok {
		owner.IsVerified = user.IsVerified
	}

	if linkedAccounts, ok := linkedAccountsByUserID[organizationUser.UserID]; ok {
		owner.LinkedAccounts = mapLinkedAccounts(linkedAccounts)
	}

	return owner
}

func mapLinkedAccounts(linkedAccounts []*accountContract.LinkedAccount) []LinkedAccount {
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
