package createorganization

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/rbac"
	accountApplicationUsers "github.com/blocknextai/platform-api/internal/account/application/users"
	"github.com/blocknextai/platform-api/internal/eventbus/application/publishing"
	organizationsDomainOrganizations "github.com/blocknextai/platform-api/internal/organizations/domain/organizations"
	organizationsDomainOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
)

type Handler struct {
	organizationRepository     organizationsDomainOrganizations.OrganizationRepository
	organizationUserRepository organizationsDomainOrganizationUsers.OrganizationUserRepository
	eventBusPublisherService   publishing.PublisherService
	transactionManager         database.TransactionManager
	userService                accountApplicationUsers.UserService
}

func New(
	organizationRepository organizationsDomainOrganizations.OrganizationRepository,
	organizationUserRepository organizationsDomainOrganizationUsers.OrganizationUserRepository,
	eventBusPublisherService publishing.PublisherService,
	transactionManager database.TransactionManager,
	userService accountApplicationUsers.UserService,
) *Handler {
	return &Handler{
		organizationRepository:     organizationRepository,
		organizationUserRepository: organizationUserRepository,
		eventBusPublisherService:   eventBusPublisherService,
		transactionManager:         transactionManager,
		userService:                userService,
	}
}

func (h *Handler) Handle(ctx context.Context, request *CreateOrganizationCommand) (*CreateOrganizationResponse, error) {
	var response *CreateOrganizationResponse

	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		organization, err := organizationsDomainOrganizations.New(
			request.Title,
			request.Description,
			false,
		)
		if err != nil {
			return err
		}

		err = h.organizationRepository.Create(txCtx, organization)
		if err != nil {
			return err
		}

		if _, err = h.userService.GetByID(txCtx, request.UserID); err != nil {
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

		if err = h.organizationUserRepository.Create(txCtx, owner); err != nil {
			return err
		}

		response = &CreateOrganizationResponse{
			ID:          organization.ID,
			Title:       organization.Title,
			Description: organization.Description,
		}

		if err = h.eventBusPublisherService.Enqueue(txCtx, organizationsDomainOrganizationUsers.OrganizationUserCreatedDomainEvent{
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
