package usecases

import (
	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/eventbus/publishing"
	accountContract "github.com/blocknextai/platform-api/internal/modules/account/contract"
	"github.com/blocknextai/platform-api/internal/modules/organizations/internal/domain/organizations"
	"github.com/blocknextai/platform-api/internal/modules/organizations/internal/domain/organizationusers"
	organizationsUseCases "github.com/blocknextai/platform-api/internal/modules/organizations/internal/usecases/organizations"
	organizationusersUseCases "github.com/blocknextai/platform-api/internal/modules/organizations/internal/usecases/organizationusers"
)

type Services struct {
	Organizations     *organizationsUseCases.Service
	OrganizationUsers *organizationusersUseCases.Service
}

type ServiceDependencies struct {
	TransactionManager       database.TransactionManager
	EventBusPublisherService publishing.PublisherService

	OrganizationRepository     organizations.OrganizationRepository
	OrganizationUserRepository organizationusers.OrganizationUserRepository
	UserService                accountContract.UserService
	LinkedAccountService       accountContract.LinkedAccountService
}

func NewServices(deps ServiceDependencies) *Services {

	return &Services{
		Organizations:     organizationsUseCases.NewService(deps.OrganizationRepository, deps.OrganizationUserRepository, deps.EventBusPublisherService, deps.TransactionManager, deps.UserService),
		OrganizationUsers: organizationusersUseCases.NewService(deps.OrganizationUserRepository, deps.EventBusPublisherService, deps.TransactionManager, deps.UserService, deps.LinkedAccountService),
	}
}
