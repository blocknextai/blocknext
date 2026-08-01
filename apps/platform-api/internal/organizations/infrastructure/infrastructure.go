package infrastructure

import (
	"github.com/blocknextai/go-packages/database"
	linkedaccountsApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/account/application/linkedaccounts"
	usersApplicationUsers "github.com/blocknextai/platform-api/internal/account/application/users"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	"github.com/blocknextai/platform-api/internal/eventbus/application/publishing"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizations/createorganization"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizations/deleteorganization"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizations/getallorganizations"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizations/getorganizationbyid"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizations/updateorganization"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizationusers/createorganizationuser"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizationusers/deleteorganizationuser"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizationusers/getallorganizationusers"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizationusers/getorganizationme"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizationusers/getorganizationuserbyuserid"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizationusers/getroles"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizationusers/updateorganizationuserinfo"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizationusers/updateorganizationuserrole"
	"github.com/blocknextai/platform-api/internal/organizations/domain/organizations"
	"github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
)

type Handlers struct {
	CreateOrganization          cqrs.Handler[*createorganization.CreateOrganizationCommand, *createorganization.CreateOrganizationResponse]
	UpdateOrganization          cqrs.Handler[*updateorganization.UpdateOrganizationCommand, *updateorganization.UpdateOrganizationResponse]
	DeleteOrganization          cqrs.Handler[*deleteorganization.DeleteOrganizationCommand, *deleteorganization.DeleteOrganizationResponse]
	GetOrganizationByID         cqrs.Handler[*getorganizationbyid.GetOrganizationByIDQuery, *getorganizationbyid.GetOrganizationByIDResponse]
	GetAllOrganizations         cqrs.Handler[*getallorganizations.GetAllOrganizationsQuery, *getallorganizations.GetAllOrganizationsResponse]
	GetOrganizationMe           cqrs.Handler[*getorganizationme.GetOrganizationMeQuery, *getorganizationme.GetOrganizationMeResponse]
	GetRoles                    cqrs.Handler[*getroles.GetRolesQuery, *getroles.GetRolesResponse]
	CreateOrganizationUser      cqrs.Handler[*createorganizationuser.CreateOrganizationUserCommand, *createorganizationuser.CreateOrganizationUserResponse]
	GetAllOrganizationUsers     cqrs.Handler[*getallorganizationusers.GetAllOrganizationUsersQuery, *getallorganizationusers.GetAllOrganizationUsersResponse]
	GetOrganizationUserByUserID cqrs.Handler[*getorganizationuserbyuserid.GetOrganizationUserByUserIDQuery, *getorganizationuserbyuserid.GetOrganizationUserByUserIDResponse]
	UpdateOrganizationUserInfo  cqrs.Handler[*updateorganizationuserinfo.UpdateOrganizationUserInfoCommand, *updateorganizationuserinfo.UpdateOrganizationUserInfoResponse]
	UpdateOrganizationUserRole  cqrs.Handler[*updateorganizationuserrole.UpdateOrganizationUserRoleCommand, *updateorganizationuserrole.UpdateOrganizationUserRoleResponse]
	DeleteOrganizationUser      cqrs.Handler[*deleteorganizationuser.DeleteOrganizationUserCommand, *deleteorganizationuser.DeleteOrganizationUserResponse]
}

type RegisterInfrastructureDeps struct {
	TransactionManager       database.TransactionManager
	EventBusPublisherService publishing.PublisherService

	OrganizationRepository     organizations.OrganizationRepository
	OrganizationUserRepository organizationusers.OrganizationUserRepository
	UserService                usersApplicationUsers.UserService
	LinkedAccountService       linkedaccountsApplicationLinkedAccounts.LinkedAccountService
}

func RegisterInfrastructure(deps RegisterInfrastructureDeps) *Handlers {
	createOrganizationHandler := cqrs.ValidationBehavior(
		createorganization.New(deps.OrganizationRepository, deps.OrganizationUserRepository, deps.EventBusPublisherService, deps.TransactionManager, deps.UserService),
	)

	return &Handlers{
		CreateOrganization:          createOrganizationHandler,
		UpdateOrganization:          cqrs.ValidationBehavior(updateorganization.New(deps.OrganizationRepository, deps.OrganizationUserRepository, deps.TransactionManager)),
		DeleteOrganization:          cqrs.ValidationBehavior(deleteorganization.New(deps.OrganizationRepository, deps.OrganizationUserRepository, deps.TransactionManager)),
		GetOrganizationByID:         cqrs.ValidationBehavior(getorganizationbyid.New(deps.OrganizationRepository, deps.OrganizationUserRepository)),
		GetAllOrganizations:         cqrs.ValidationBehavior(getallorganizations.New(deps.OrganizationRepository)),
		GetOrganizationMe:           cqrs.ValidationBehavior(getorganizationme.New(deps.OrganizationUserRepository)),
		GetRoles:                    cqrs.ValidationBehavior(getroles.New()),
		CreateOrganizationUser:      cqrs.ValidationBehavior(createorganizationuser.New(deps.OrganizationUserRepository, deps.EventBusPublisherService, deps.TransactionManager, deps.UserService, deps.LinkedAccountService)),
		GetAllOrganizationUsers:     cqrs.ValidationBehavior(getallorganizationusers.New(deps.OrganizationUserRepository, deps.UserService, deps.LinkedAccountService)),
		GetOrganizationUserByUserID: cqrs.ValidationBehavior(getorganizationuserbyuserid.New(deps.OrganizationUserRepository, deps.UserService, deps.LinkedAccountService)),
		UpdateOrganizationUserInfo:  cqrs.ValidationBehavior(updateorganizationuserinfo.New(deps.OrganizationUserRepository, deps.TransactionManager)),
		UpdateOrganizationUserRole:  cqrs.ValidationBehavior(updateorganizationuserrole.New(deps.OrganizationUserRepository, deps.EventBusPublisherService, deps.TransactionManager)),
		DeleteOrganizationUser:      cqrs.ValidationBehavior(deleteorganizationuser.New(deps.OrganizationUserRepository, deps.TransactionManager)),
	}
}
