package organizationusers

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/rbac"
	organizationsApplicationUsers "github.com/blocknextai/platform-api/internal/modules/organizations/internal/application/organizationusers"
	"github.com/blocknextai/platform-api/internal/modules/organizations/internal/domain/organizationusers"
)

type CreateOrganizationUserCommand struct {
	OrganizationID uuid.UUID
	Identifier     string
	Alias          string
	Role           string
}

func (c *CreateOrganizationUserCommand) Validate() error {
	if strings.TrimSpace(c.Identifier) == "" {
		return organizationsApplicationUsers.ErrIdentifierIsRequired
	}

	return nil
}

type CreateOrganizationUserResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organizationId"`
	UserID         uuid.UUID `json:"userId"`
	Role           string    `json:"role"`
}

func (s *Service) CreateOrganizationUser(ctx context.Context, command *CreateOrganizationUserCommand) (*CreateOrganizationUserResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	var response *CreateOrganizationUserResponse

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		user, err := s.linkedAccountService.GetByIdentifier(txCtx, command.Identifier)
		if err != nil {
			return err
		}

		organizationUser, err := s.organizationUserRepository.GetByOrganizationIDAndUserID(txCtx, command.OrganizationID, user.UserID)
		if err != nil && !errors.Is(err, organizationusers.ErrOrganizationUserNotFound) {
			return err
		}

		if organizationUser != nil {
			return organizationusers.ErrOrganizationUserAlreadyExists
		}

		role := rbac.OrganizationViewerRole.Name
		if command.Role != "" {
			role = command.Role
		}

		organizationUser, err = organizationusers.New(command.OrganizationID, user.UserID, role, command.Alias)
		if err != nil {
			return err
		}

		err = s.organizationUserRepository.Create(txCtx, organizationUser)
		if err != nil {
			return err
		}

		response = &CreateOrganizationUserResponse{
			ID:             organizationUser.ID,
			OrganizationID: organizationUser.OrganizationID,
			UserID:         organizationUser.UserID,
			Role:           organizationUser.Role,
		}

		return s.eventBusPublisherService.Enqueue(txCtx, organizationusers.OrganizationUserCreatedDomainEvent{
			ID:             organizationUser.ID,
			OrganizationID: organizationUser.OrganizationID,
			UserID:         organizationUser.UserID,
			Role:           organizationUser.Role,
		})
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
