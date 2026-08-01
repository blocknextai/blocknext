package deleteorganizationuser

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/rbac"
	organizationsApplicationUsers "github.com/blocknextai/platform-api/internal/organizations/application/organizationusers"
	"github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
)

type Handler struct {
	organizationUserRepository organizationusers.OrganizationUserRepository
	transactionManager         database.TransactionManager
}

func New(
	organizationUserRepository organizationusers.OrganizationUserRepository,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		organizationUserRepository: organizationUserRepository,
		transactionManager:         transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *DeleteOrganizationUserCommand) (*DeleteOrganizationUserResponse, error) {
	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		foundUser, err := h.organizationUserRepository.GetByIDAndOrganizationID(txCtx, command.UserID, command.OrganizationID)
		if err != nil {
			return organizationusers.ErrOrganizationUserNotFound
		}

		if !command.ForceDelete {
			if foundUser.Role == rbac.OrganizationOwnerRole.Name {
				return organizationsApplicationUsers.ErrCantDeleteTheOwnerOrganizationUser
			}
		}

		foundUser, err = foundUser.Delete()
		if err != nil {
			return err
		}

		err = h.organizationUserRepository.Delete(txCtx, foundUser)
		if err != nil {
			return organizationsApplicationUsers.ErrFailedToDeleteOrganizationUser.WithCause(err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &DeleteOrganizationUserResponse{}, nil
}
