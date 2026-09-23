package organizationusers

import (
	"context"

	"github.com/google/uuid"

	organizationsApplicationUsers "github.com/blocknextai/platform-api/internal/modules/organizations/internal/application/organizationusers"
)

type UpdateOrganizationUserInfoCommand struct {
	UserID             uuid.UUID
	OrganizationID     uuid.UUID
	OrganizationUserID uuid.UUID
	Alias              *string
}

type UpdateOrganizationUserInfoResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organizationId"`
	UserID         uuid.UUID `json:"userId"`
	Role           string    `json:"role"`
	Alias          *string   `json:"alias"`
}

func (s *Service) UpdateOrganizationUserInfo(ctx context.Context, command *UpdateOrganizationUserInfoCommand) (*UpdateOrganizationUserInfoResponse, error) {
	var response *UpdateOrganizationUserInfoResponse

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		organizationUser, err := s.organizationUserRepository.GetByIDAndOrganizationID(txCtx, command.OrganizationUserID, command.OrganizationID)
		if err != nil {
			return err
		}

		organizationUser, err = organizationUser.Update(organizationUser.Role, *command.Alias)
		if err != nil {
			return err
		}

		err = s.organizationUserRepository.Update(txCtx, organizationUser)
		if err != nil {
			return organizationsApplicationUsers.ErrFailedToUpdateOrganizationUserInfo.WithCause(err)
		}

		response = &UpdateOrganizationUserInfoResponse{
			ID:             organizationUser.ID,
			OrganizationID: organizationUser.OrganizationID,
			UserID:         organizationUser.UserID,
			Role:           organizationUser.Role,
			Alias:          new(organizationUser.Alias),
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
