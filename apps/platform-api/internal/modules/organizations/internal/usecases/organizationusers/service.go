package organizationusers

import (
	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/eventbus/publishing"
	accountContract "github.com/blocknextai/platform-api/internal/modules/account/contract"
	"github.com/blocknextai/platform-api/internal/modules/organizations/internal/domain/organizationusers"
)

type Service struct {
	organizationUserRepository organizationusers.OrganizationUserRepository
	eventBusPublisherService   publishing.PublisherService
	transactionManager         database.TransactionManager
	userService                accountContract.UserService
	linkedAccountService       accountContract.LinkedAccountService
}

func NewService(
	organizationUserRepository organizationusers.OrganizationUserRepository,
	eventBusPublisherService publishing.PublisherService,
	transactionManager database.TransactionManager,
	userService accountContract.UserService,
	linkedAccountService accountContract.LinkedAccountService,
) *Service {
	return &Service{
		organizationUserRepository: organizationUserRepository,
		eventBusPublisherService:   eventBusPublisherService,
		transactionManager:         transactionManager,
		userService:                userService,
		linkedAccountService:       linkedAccountService,
	}
}
