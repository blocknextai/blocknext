package updateorganization

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	organizationsApplicationOrganizations "github.com/blocknextai/platform-api/internal/organizations/application/organizations"
	organizationsDomainOrganizations "github.com/blocknextai/platform-api/internal/organizations/domain/organizations"
	organizationsDomainOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
)

type Handler struct {
	organizationRepository     organizationsDomainOrganizations.OrganizationRepository
	organizationUserRepository organizationsDomainOrganizationUsers.OrganizationUserRepository
	transactionManager         database.TransactionManager
}

func New(
	organizationRepository organizationsDomainOrganizations.OrganizationRepository,
	organizationUserRepository organizationsDomainOrganizationUsers.OrganizationUserRepository,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		organizationRepository:     organizationRepository,
		organizationUserRepository: organizationUserRepository,
		transactionManager:         transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, request *UpdateOrganizationCommand) (*UpdateOrganizationResponse, error) {
	var response *UpdateOrganizationResponse

	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		if _, err := h.organizationUserRepository.GetByOrganizationIDAndUserID(txCtx, request.OrganizationID, request.UserID); err != nil {
			return organizationsDomainOrganizations.ErrOrganizationNotFound
		}

		organization, err := h.organizationRepository.GetByID(txCtx, request.OrganizationID)
		if err != nil {
			return err
		}

		organization, err = organization.Update(
			request.Title,
			request.Description,
		)
		if err != nil {
			return err
		}

		if err := h.organizationRepository.Update(txCtx, organization); err != nil {
			return organizationsApplicationOrganizations.ErrFailedToUpdateOrganization.WithCause(err)
		}

		response = &UpdateOrganizationResponse{
			ID:          organization.ID,
			Title:       organization.Title,
			Description: organization.Description,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
