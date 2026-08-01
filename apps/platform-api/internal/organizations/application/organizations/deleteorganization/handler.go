package deleteorganization

import (
	"context"
	"time"

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

func (h *Handler) Handle(ctx context.Context, request *DeleteOrganizationCommand) (*DeleteOrganizationResponse, error) {
	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		if _, err := h.organizationUserRepository.GetByOrganizationIDAndUserID(txCtx, request.OrganizationID, request.UserID); err != nil {
			return organizationsDomainOrganizations.ErrOrganizationNotFound
		}

		organization, err := h.organizationRepository.GetByID(txCtx, request.OrganizationID)
		if err != nil {
			return organizationsDomainOrganizations.ErrOrganizationNotFound
		}

		organization, err = organization.Delete()
		if err != nil {
			return err
		}

		if err := h.organizationUserRepository.DeleteAllByOrganizationID(txCtx, request.OrganizationID, time.Now().UTC()); err != nil {
			return organizationsApplicationOrganizations.ErrFailedToDeleteOrganization.WithCause(err)
		}

		if err := h.organizationRepository.Delete(txCtx, organization); err != nil {
			return organizationsApplicationOrganizations.ErrFailedToDeleteOrganization.WithCause(err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &DeleteOrganizationResponse{}, nil
}
