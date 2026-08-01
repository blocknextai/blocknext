package createorganizationuser

import (
	"context"
	"errors"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/rbac"
	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/account/application/linkedaccounts"
	accountApplicationUsers "github.com/blocknextai/platform-api/internal/account/application/users"
	"github.com/blocknextai/platform-api/internal/eventbus/application/publishing"
	"github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
)

type Handler struct {
	organizationUserRepository organizationusers.OrganizationUserRepository
	eventBusPublisherService   publishing.PublisherService
	transactionManager         database.TransactionManager
	userService                accountApplicationUsers.UserService
	linkedAccountService       accountApplicationLinkedAccounts.LinkedAccountService
}

func New(
	organizationUserRepository organizationusers.OrganizationUserRepository,
	eventBusPublisherService publishing.PublisherService,
	transactionManager database.TransactionManager,
	userService accountApplicationUsers.UserService,
	linkedAccountService accountApplicationLinkedAccounts.LinkedAccountService,
) *Handler {
	return &Handler{
		organizationUserRepository: organizationUserRepository,
		eventBusPublisherService:   eventBusPublisherService,
		transactionManager:         transactionManager,
		userService:                userService,
		linkedAccountService:       linkedAccountService,
	}
}

func (h *Handler) Handle(ctx context.Context, command *CreateOrganizationUserCommand) (*CreateOrganizationUserResponse, error) {
	var response *CreateOrganizationUserResponse

	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		user, err := h.linkedAccountService.GetByIdentifier(txCtx, command.Identifier)
		if err != nil {
			return err
		}

		organizationUser, err := h.organizationUserRepository.GetByOrganizationIDAndUserID(txCtx, command.OrganizationID, user.UserID)
		if err != nil && !errors.Is(err, organizationusers.ErrOrganizationUserNotFound) {
			return err
		}

		if organizationUser != nil {
			return organizationusers.ErrOrganizationUserAlreadyExists
		}

		role := rbac.OrganizationViewerRole.Name
		if command.Role != "" {
			role = command.Role
		}

		organizationUser, err = organizationusers.New(command.OrganizationID, user.UserID, role, command.Alias)
		if err != nil {
			return err
		}

		err = h.organizationUserRepository.Create(txCtx, organizationUser)
		if err != nil {
			return err
		}

		response = &CreateOrganizationUserResponse{
			ID:             organizationUser.ID,
			OrganizationID: organizationUser.OrganizationID,
			UserID:         organizationUser.UserID,
			Role:           organizationUser.Role,
		}

		return h.eventBusPublisherService.Enqueue(txCtx, organizationusers.OrganizationUserCreatedDomainEvent{
			ID:             organizationUser.ID,
			OrganizationID: organizationUser.OrganizationID,
			UserID:         organizationUser.UserID,
			Role:           organizationUser.Role,
		})
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
