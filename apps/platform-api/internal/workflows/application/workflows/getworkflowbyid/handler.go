package getworkflowbyid

import (
	"context"
	"log/slog"

	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/account/application/linkedaccounts"
	accountApplicationUsers "github.com/blocknextai/platform-api/internal/account/application/users"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
	organizationsApplicationOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/application/organizationusers"
	organizationsDomainOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
	workflowsDomainWorkflows "github.com/blocknextai/platform-api/internal/workflows/domain/workflows"
	"github.com/google/uuid"
)

type Handler struct {
	workflowRepository      workflowsDomainWorkflows.WorkflowRepository
	organizationUserService organizationsApplicationOrganizationUsers.OrganizationUserService
	userService             accountApplicationUsers.UserService
	linkedAccountService    accountApplicationLinkedAccounts.LinkedAccountService
}

func New(
	workflowRepository workflowsDomainWorkflows.WorkflowRepository,
	organizationUserService organizationsApplicationOrganizationUsers.OrganizationUserService,
	userService accountApplicationUsers.UserService,
	linkedAccountService accountApplicationLinkedAccounts.LinkedAccountService,
) *Handler {
	return &Handler{
		workflowRepository:      workflowRepository,
		organizationUserService: organizationUserService,
		userService:             userService,
		linkedAccountService:    linkedAccountService,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetWorkflowByIDQuery) (*GetWorkflowByIDResponse, error) {
	workflow, err := h.workflowRepository.GetByOrganizationIDAndID(ctx, request.OrganizationID, request.WorkflowID)
	if err != nil {
		return nil, err
	}

	organizationUsers, err := h.organizationUserService.GetAllByIDs(ctx, []uuid.UUID{workflow.OwnerID})
	if err != nil {
		slog.WarnContext(ctx, "Failed to get organization users for workflow",
			"component", "getworkflowbyid",
			"organization_id", request.OrganizationID,
			"error", err)
		organizationUsers = []*organizationsDomainOrganizationUsers.OrganizationUser{}
	}

	organizationUsersByID := make(map[uuid.UUID]*organizationsDomainOrganizationUsers.OrganizationUser, len(organizationUsers))
	userIDs := make([]uuid.UUID, 0, len(organizationUsers))
	for _, organizationUser := range organizationUsers {
		organizationUsersByID[organizationUser.ID] = organizationUser
		userIDs = append(userIDs, organizationUser.UserID)
	}

	users, err := h.userService.GetAllByIDs(ctx, userIDs)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get users for workflow",
			"component", "getworkflowbyid",
			"error", err)
		users = []*accountDomainUsers.User{}
	}

	usersByID := make(map[uuid.UUID]*accountDomainUsers.User, len(users))
	for _, user := range users {
		usersByID[user.ID] = user
	}

	linkedAccounts, err := h.linkedAccountService.GetAllByUserIDs(ctx, userIDs)
	if err != nil {
		slog.WarnContext(ctx, "Failed to get linked accounts for workflow",
			"component", "getworkflowbyid",
			"error", err)
		linkedAccounts = []*accountDomainLinkedAccounts.LinkedAccount{}
	}

	linkedAccountsByUserID := make(map[uuid.UUID][]*accountDomainLinkedAccounts.LinkedAccount)
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
