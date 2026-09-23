package workflows

import (
	"github.com/blocknextai/go-packages/database"
	accountContract "github.com/blocknextai/platform-api/internal/modules/account/contract"
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
	organizationsContract "github.com/blocknextai/platform-api/internal/modules/organizations/contract"
	workflowsDomainWorkflows "github.com/blocknextai/platform-api/internal/modules/workflows/internal/domain/workflows"
)

type Service struct {
	workflowRepository      workflowsDomainWorkflows.WorkflowRepository
	transactionManager      database.TransactionManager
	organizationUserService organizationsContract.OrganizationUserService
	userService             accountContract.UserService
	linkedAccountService    accountContract.LinkedAccountService
	credentialService       nodeengineContract.CredentialService
	nodeService             nodeengineContract.NodeService
}

func NewService(
	workflowRepository workflowsDomainWorkflows.WorkflowRepository,
	transactionManager database.TransactionManager,
	organizationUserService organizationsContract.OrganizationUserService,
	userService accountContract.UserService,
	linkedAccountService accountContract.LinkedAccountService,
	credentialService nodeengineContract.CredentialService,
	nodeService nodeengineContract.NodeService,
) *Service {
	return &Service{
		workflowRepository:      workflowRepository,
		transactionManager:      transactionManager,
		organizationUserService: organizationUserService,
		userService:             userService,
		linkedAccountService:    linkedAccountService,
		credentialService:       credentialService,
		nodeService:             nodeService,
	}
}
