package updateorganizationuserinfo

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	organizationsApplicationUsers "github.com/blocknextai/platform-api/internal/organizations/application/organizationusers"
	organizationsDomainOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
)

type Handler struct {
	organizationUserRepository organizationsDomainOrganizationUsers.OrganizationUserRepository
	transactionManager         database.TransactionManager
}

func New(
	organizationUserRepository organizationsDomainOrganizationUsers.OrganizationUserRepository,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		organizationUserRepository: organizationUserRepository,
		transactionManager:         transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *UpdateOrganizationUserInfoCommand) (*UpdateOrganizationUserInfoResponse, error) {
	var response *UpdateOrganizationUserInfoResponse

	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		organizationUser, err := h.organizationUserRepository.GetByIDAndOrganizationID(txCtx, command.OrganizationUserID, command.OrganizationID)
		if err != nil {
			return err
		}

		organizationUser, err = organizationUser.Update(organizationUser.Role, *command.Alias)
		if err != nil {
			return err
		}

		err = h.organizationUserRepository.Update(txCtx, organizationUser)
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
