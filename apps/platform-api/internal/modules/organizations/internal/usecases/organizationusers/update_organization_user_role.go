package organizationusers

import (
	"context"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/rbac"
	organizationsApplicationUsers "github.com/blocknextai/platform-api/internal/modules/organizations/internal/application/organizationusers"
	"github.com/blocknextai/platform-api/internal/modules/organizations/internal/domain/organizationusers"
)

type UpdateOrganizationUserRoleCommand struct {
	UserID             uuid.UUID
	OrganizationID     uuid.UUID
	OrganizationUserID uuid.UUID
	Role               string
}

type UpdateOrganizationUserRoleResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organizationId"`
	UserID         uuid.UUID `json:"userId"`
	Role           string    `json:"role"`
}

func (s *Service) UpdateOrganizationUserRole(ctx context.Context, command *UpdateOrganizationUserRoleCommand) (*UpdateOrganizationUserRoleResponse, error) {
	var response *UpdateOrganizationUserRoleResponse

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		organizationUser, err := s.organizationUserRepository.GetByIDAndOrganizationID(txCtx, command.OrganizationUserID, command.OrganizationID)
		if err != nil {
			return err
		}

		if command.UserID == organizationUser.UserID {
			return organizationsApplicationUsers.ErrUserCannotChangeOwnRole
		}

		if command.Role == rbac.OrganizationOwnerRole.Name {
			hasOwner, err := s.organizationUserRepository.HasOwner(txCtx, command.OrganizationID)
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

		err = s.organizationUserRepository.Update(txCtx, organizationUser)
		if err != nil {
			return organizationsApplicationUsers.ErrFailedToUpdateOrganizationUserRole.WithCause(err)
		}

		response = &UpdateOrganizationUserRoleResponse{
			ID:             organizationUser.ID,
			OrganizationID: organizationUser.OrganizationID,
			UserID:         organizationUser.UserID,
			Role:           organizationUser.Role,
		}

		return s.eventBusPublisherService.Enqueue(txCtx, organizationusers.OrganizationUserRoleChangedDomainEvent{
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
