package organizations

import (
	"context"

	"github.com/google/uuid"

	organizationsApplicationOrganizations "github.com/blocknextai/platform-api/internal/modules/organizations/internal/application/organizations"
	organizationsDomainOrganizations "github.com/blocknextai/platform-api/internal/modules/organizations/internal/domain/organizations"
)

type UpdateOrganizationCommand struct {
	UserID         uuid.UUID
	OrganizationID uuid.UUID
	Title          string
	Description    *string
}

type UpdateOrganizationResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
}

func (s *Service) UpdateOrganization(ctx context.Context, request *UpdateOrganizationCommand) (*UpdateOrganizationResponse, error) {
	var response *UpdateOrganizationResponse

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		if _, err := s.organizationUserRepository.GetByOrganizationIDAndUserID(txCtx, request.OrganizationID, request.UserID); err != nil {
			return organizationsDomainOrganizations.ErrOrganizationNotFound
		}

		organization, err := s.organizationRepository.GetByID(txCtx, request.OrganizationID)
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

		if err := s.organizationRepository.Update(txCtx, organization); err != nil {
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
