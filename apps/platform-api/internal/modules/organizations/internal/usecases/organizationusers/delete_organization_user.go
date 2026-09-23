package organizationusers

import (
	"context"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/rbac"
	organizationsApplicationUsers "github.com/blocknextai/platform-api/internal/modules/organizations/internal/application/organizationusers"
	"github.com/blocknextai/platform-api/internal/modules/organizations/internal/domain/organizationusers"
)

type DeleteOrganizationUserCommand struct {
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	ForceDelete    bool
}

type DeleteOrganizationUserResponse struct{}

func (s *Service) DeleteOrganizationUser(ctx context.Context, command *DeleteOrganizationUserCommand) (*DeleteOrganizationUserResponse, error) {
	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		foundUser, err := s.organizationUserRepository.GetByIDAndOrganizationID(txCtx, command.UserID, command.OrganizationID)
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

		err = s.organizationUserRepository.Delete(txCtx, foundUser)
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
