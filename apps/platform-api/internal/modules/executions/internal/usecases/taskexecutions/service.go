package taskexecutions

import (
	"github.com/blocknextai/go-packages/database"
	accountContract "github.com/blocknextai/platform-api/internal/modules/account/contract"
	executionsApplicationNodeExecutions "github.com/blocknextai/platform-api/internal/modules/executions/internal/application/nodeexecutions"
	"github.com/blocknextai/platform-api/internal/modules/executions/internal/application/taskexecutions/workflowresolver"
	"github.com/blocknextai/platform-api/internal/modules/executions/internal/domain/taskexecutions"
	organizationsContract "github.com/blocknextai/platform-api/internal/modules/organizations/contract"
)

type Service struct {
	taskExecutionRepository taskexecutions.TaskExecutionRepository
	transactionManager      database.TransactionManager
	workflowResolver        workflowresolver.Resolver
	nodeExecutionService    executionsApplicationNodeExecutions.NodeExecutionService
	organizationUserService organizationsContract.OrganizationUserService
	userService             accountContract.UserService
	linkedAccountService    accountContract.LinkedAccountService
}

func NewService(
	taskExecutionRepository taskexecutions.TaskExecutionRepository,
	transactionManager database.TransactionManager,
	workflowResolver workflowresolver.Resolver,
	nodeExecutionService executionsApplicationNodeExecutions.NodeExecutionService,
	organizationUserService organizationsContract.OrganizationUserService,
	userService accountContract.UserService,
	linkedAccountService accountContract.LinkedAccountService,
) *Service {
	return &Service{
		taskExecutionRepository: taskExecutionRepository,
		transactionManager:      transactionManager,
		workflowResolver:        workflowResolver,
		nodeExecutionService:    nodeExecutionService,
		organizationUserService: organizationUserService,
		userService:             userService,
		linkedAccountService:    linkedAccountService,
	}
}
