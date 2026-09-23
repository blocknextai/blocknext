package organizations

import (
	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/eventbus/publishing"
	accountContract "github.com/blocknextai/platform-api/internal/modules/account/contract"
	organizationsDomainOrganizations "github.com/blocknextai/platform-api/internal/modules/organizations/internal/domain/organizations"
	organizationsDomainOrganizationUsers "github.com/blocknextai/platform-api/internal/modules/organizations/internal/domain/organizationusers"
)

type Service struct {
	organizationRepository     organizationsDomainOrganizations.OrganizationRepository
	organizationUserRepository organizationsDomainOrganizationUsers.OrganizationUserRepository
	eventBusPublisherService   publishing.PublisherService
	transactionManager         database.TransactionManager
	userService                accountContract.UserService
}

func NewService(
	organizationRepository organizationsDomainOrganizations.OrganizationRepository,
	organizationUserRepository organizationsDomainOrganizationUsers.OrganizationUserRepository,
	eventBusPublisherService publishing.PublisherService,
	transactionManager database.TransactionManager,
	userService accountContract.UserService,
) *Service {
	return &Service{
		organizationRepository:     organizationRepository,
		organizationUserRepository: organizationUserRepository,
		eventBusPublisherService:   eventBusPublisherService,
		transactionManager:         transactionManager,
		userService:                userService,
	}
}
