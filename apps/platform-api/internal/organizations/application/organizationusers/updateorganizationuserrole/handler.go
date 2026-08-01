package updateorganizationuserrole

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/rbac"
	"github.com/blocknextai/platform-api/internal/eventbus/application/publishing"
	organizationsApplicationUsers "github.com/blocknextai/platform-api/internal/organizations/application/organizationusers"
	"github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
)

type Handler struct {
	organizationUserRepository organizationusers.OrganizationUserRepository
	eventBusPublisherService   publishing.PublisherService
	transactionManager         database.TransactionManager
}

func New(
	organizationUserRepository organizationusers.OrganizationUserRepository,
	eventBusPublisherService publishing.PublisherService,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		organizationUserRepository: organizationUserRepository,
		eventBusPublisherService:   eventBusPublisherService,
		transactionManager:         transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *UpdateOrganizationUserRoleCommand) (*UpdateOrganizationUserRoleResponse, error) {
	var response *UpdateOrganizationUserRoleResponse

	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		organizationUser, err := h.organizationUserRepository.GetByIDAndOrganizationID(txCtx, command.OrganizationUserID, command.OrganizationID)
		if err != nil {
			return err
		}

		if command.UserID == organizationUser.UserID {
			return organizationsApplicationUsers.ErrUserCannotChangeOwnRole
		}

		if command.Role == rbac.OrganizationOwnerRole.Name {
			hasOwner, err := h.organizationUserRepository.HasOwner(txCtx, command.OrganizationID)
			if err != nil {
				return err
			}

			if hasOwner {
				return organizationsApplicationUsers.ErrCouldBeOneOwner
			}
		}

		oldRole := organizationUser.Role

		organizationUser, err = organizationUser.Update(command.Role, organizationUser.Alias)
		if err != nil {
			return err
		}

		err = h.organizationUserRepository.Update(txCtx, organizationUser)
		if err != nil {
			return organizationsApplicationUsers.ErrFailedToUpdateOrganizationUserRole.WithCause(err)
		}

		response = &UpdateOrganizationUserRoleResponse{
			ID:             organizationUser.ID,
			OrganizationID: organizationUser.OrganizationID,
			UserID:         organizationUser.UserID,
			Role:           organizationUser.Role,
		}

		return h.eventBusPublisherService.Enqueue(txCtx, organizationusers.OrganizationUserRoleChangedDomainEvent{
			OrganizationID: organizationUser.OrganizationID,
			UserID:         organizationUser.UserID,
			OldRole:        oldRole,
			NewRole:        organizationUser.Role,
		})
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
