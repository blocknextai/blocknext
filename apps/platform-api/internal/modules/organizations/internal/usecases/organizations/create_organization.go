package organizations

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/rbac"
	organizationsApplicationOrganizations "github.com/blocknextai/platform-api/internal/modules/organizations/internal/application/organizations"
	organizationsDomainOrganizations "github.com/blocknextai/platform-api/internal/modules/organizations/internal/domain/organizations"
	organizationsDomainOrganizationUsers "github.com/blocknextai/platform-api/internal/modules/organizations/internal/domain/organizationusers"
)

type CreateOrganizationCommand struct {
	UserID      uuid.UUID
	Title       string
	Description *string
}

const (
	MaxTitleLength = 255
)

func (c *CreateOrganizationCommand) Validate() error {
	if strings.TrimSpace(c.Title) == "" {
		return organizationsApplicationOrganizations.ErrInvalidTitle
	}

	if len(strings.TrimSpace(c.Title)) > MaxTitleLength {
		return organizationsApplicationOrganizations.ErrTitleTooLong
	}

	return nil
}

type CreateOrganizationResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
}

func (s *Service) CreateOrganization(ctx context.Context, request *CreateOrganizationCommand) (*CreateOrganizationResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	var response *CreateOrganizationResponse

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		organization, err := organizationsDomainOrganizations.New(
			request.Title,
			request.Description,
			false,
		)
		if err != nil {
			return err
		}

		err = s.organizationRepository.Create(txCtx, organization)
		if err != nil {
			return err
		}

		if _, err = s.userService.GetByID(txCtx, request.UserID); err != nil {
			return err
		}

		owner, err := organizationsDomainOrganizationUsers.New(
			organization.ID,
			request.UserID,
			rbac.OrganizationOwnerRole.Name,
			"",
		)
		if err != nil {
			return err
		}

		if err = s.organizationUserRepository.Create(txCtx, owner); err != nil {
			return err
		}

		response = &CreateOrganizationResponse{
			ID:          organization.ID,
			Title:       organization.Title,
			Description: organization.Description,
		}

		if err = s.eventBusPublisherService.Enqueue(txCtx, organizationsDomainOrganizationUsers.OrganizationUserCreatedDomainEvent{
			ID:             owner.ID,
			OrganizationID: owner.OrganizationID,
			UserID:         owner.UserID,
			Role:           owner.Role,
		}); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
