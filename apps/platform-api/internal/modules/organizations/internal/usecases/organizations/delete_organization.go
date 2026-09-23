package organizations

import (
	"context"
	"time"

	"github.com/google/uuid"

	organizationsApplicationOrganizations "github.com/blocknextai/platform-api/internal/modules/organizations/internal/application/organizations"
	organizationsDomainOrganizations "github.com/blocknextai/platform-api/internal/modules/organizations/internal/domain/organizations"
)

type DeleteOrganizationCommand struct {
	UserID         uuid.UUID
	OrganizationID uuid.UUID
}

type DeleteOrganizationResponse struct{}

func (s *Service) DeleteOrganization(ctx context.Context, request *DeleteOrganizationCommand) (*DeleteOrganizationResponse, error) {
	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		if _, err := s.organizationUserRepository.GetByOrganizationIDAndUserID(txCtx, request.OrganizationID, request.UserID); err != nil {
			return organizationsDomainOrganizations.ErrOrganizationNotFound
		}

		organization, err := s.organizationRepository.GetByID(txCtx, request.OrganizationID)
		if err != nil {
			return organizationsDomainOrganizations.ErrOrganizationNotFound
		}

		organization, err = organization.Delete()
		if err != nil {
			return err
		}

		if err := s.organizationUserRepository.DeleteAllByOrganizationID(txCtx, request.OrganizationID, time.Now().UTC()); err != nil {
			return organizationsApplicationOrganizations.ErrFailedToDeleteOrganization.WithCause(err)
		}

		if err := s.organizationRepository.Delete(txCtx, organization); err != nil {
			return organizationsApplicationOrganizations.ErrFailedToDeleteOrganization.WithCause(err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &DeleteOrganizationResponse{}, nil
}
