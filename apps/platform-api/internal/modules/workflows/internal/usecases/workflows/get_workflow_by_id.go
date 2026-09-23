package workflows

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/dag"
	accountContract "github.com/blocknextai/platform-api/internal/modules/account/contract"
	organizationsContract "github.com/blocknextai/platform-api/internal/modules/organizations/contract"
	workflowsDomainWorkflows "github.com/blocknextai/platform-api/internal/modules/workflows/internal/domain/workflows"
)

type GetWorkflowByIDQuery struct {
	OrganizationID uuid.UUID
	WorkflowID     uuid.UUID
}

type GetWorkflowByIDLinkedAccount struct {
	AuthProvider string  `json:"authProvider"`
	DisplayName  *string `json:"displayName"`
	IsPrimary    bool    `json:"isPrimary"`
}

type GetWorkflowByIDOwner struct {
	ID             uuid.UUID                      `json:"id"`
	Alias          string                         `json:"alias"`
	IsVerified     bool                           `json:"isVerified"`
	LinkedAccounts []GetWorkflowByIDLinkedAccount `json:"linkedAccounts"`
}

type GetWorkflowByIDResponse struct {
	ID             uuid.UUID             `json:"id"`
	OrganizationID uuid.UUID             `json:"organizationId"`
	Owner          *GetWorkflowByIDOwner `json:"owner"`
	Title          string                `json:"title"`
	Description    *string               `json:"description"`
	IsPinned       bool                  `json:"isPinned"`
	Nodes          []dag.Node            `json:"nodes"`
	Edges          []dag.Edge            `json:"edges"`
}

func (s *Service) GetWorkflowByID(ctx context.Context, request *GetWorkflowByIDQuery) (*GetWorkflowByIDResponse, error) {
	workflow, err := s.workflowRepository.GetByOrganizationIDAndID(ctx, request.OrganizationID, request.WorkflowID)
	if err != nil {
		return nil, err
	}

	organizationUsers, err := s.organizationUserService.GetAllByIDs(ctx, []uuid.UUID{workflow.OwnerID})
	if err != nil {
		slog.WarnContext(ctx, "Failed to get organization users for workflow",
			"component", "getworkflowbyid",
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
		slog.WarnContext(ctx, "Failed to get users for workflow",
			"component", "getworkflowbyid",
			"error", err)
		users = []*accountContract.User{}
	}

	usersByID := make(map[uuid.UUID]*accountContract.User, len(users))
	for _, user := range users {
		usersByID[user.ID] = user
	}

	linkedAccounts, err := s.linkedAccountService.GetAllByUserIDs(ctx, userIDs)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get linked accounts for workflow",
			"component", "getworkflowbyid",
			"error", err)
		linkedAccounts = []*accountContract.LinkedAccount{}
	}

	linkedAccountsByUserID := make(map[uuid.UUID][]*accountContract.LinkedAccount)
	for _, linkedAccount := range linkedAccounts {
		linkedAccountsByUserID[linkedAccount.UserID] = append(linkedAccountsByUserID[linkedAccount.UserID], linkedAccount)
	}

	return MapWorkflowToResponse(
		workflow,
		workflow.Nodes,
		workflow.Edges,
		organizationUsersByID,
		usersByID,
		linkedAccountsByUserID,
	), nil
}

func MapWorkflowToResponse(
	workflow *workflowsDomainWorkflows.Workflow,
	nodes []dag.Node,
	edges []dag.Edge,
	organizationUsersByID map[uuid.UUID]*organizationsContract.OrganizationUser,
	usersByID map[uuid.UUID]*accountContract.User,
	linkedAccountsByUserID map[uuid.UUID][]*accountContract.LinkedAccount,
) *GetWorkflowByIDResponse {
	response := &GetWorkflowByIDResponse{
		ID:             workflow.ID,
		OrganizationID: workflow.OrganizationID,
		Owner:          getWorkflowByIDBuildOwner(workflow.OwnerID, organizationUsersByID, usersByID, linkedAccountsByUserID),
		Title:          workflow.Title,
		Description:    workflow.Description,
		IsPinned:       workflow.IsPinned,
		Nodes:          nodes,
		Edges:          edges,
	}

	return response
}

func getWorkflowByIDBuildOwner(
	ownerID uuid.UUID,
	organizationUsersByID map[uuid.UUID]*organizationsContract.OrganizationUser,
	usersByID map[uuid.UUID]*accountContract.User,
	linkedAccountsByUserID map[uuid.UUID][]*accountContract.LinkedAccount,
) *GetWorkflowByIDOwner {
	organizationUser, ok := organizationUsersByID[ownerID]
	if !ok {
		return nil
	}

	owner := &GetWorkflowByIDOwner{
		ID:    organizationUser.ID,
		Alias: organizationUser.Alias,
	}

	if user, ok := usersByID[organizationUser.UserID]; ok {
		owner.IsVerified = user.IsVerified
	}

	if linkedAccounts, ok := linkedAccountsByUserID[organizationUser.UserID]; ok {
		owner.LinkedAccounts = getWorkflowByIDMapLinkedAccounts(linkedAccounts)
	}

	return owner
}

func getWorkflowByIDMapLinkedAccounts(linkedAccounts []*accountContract.LinkedAccount) []GetWorkflowByIDLinkedAccount {
	result := make([]GetWorkflowByIDLinkedAccount, 0, len(linkedAccounts))
	for _, linkedAccount := range linkedAccounts {
		result = append(result, GetWorkflowByIDLinkedAccount{
			AuthProvider: linkedAccount.AuthProvider.String(),
			DisplayName:  linkedAccount.DisplayName,
			IsPrimary:    linkedAccount.IsPrimary,
		})
	}
	return result
}
