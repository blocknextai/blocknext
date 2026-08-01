package executions

import (
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/account/application/linkedaccounts"
	accountApplicationUsers "github.com/blocknextai/platform-api/internal/account/application/users"
	"github.com/blocknextai/platform-api/internal/common/presentation/auth"
	nodeexecutionsApplication "github.com/blocknextai/platform-api/internal/executions/application/nodeexecutions"
	taskclaimsApplication "github.com/blocknextai/platform-api/internal/executions/application/taskclaims"
	taskexecutionsApplication "github.com/blocknextai/platform-api/internal/executions/application/taskexecutions"
	executionsInfrastructure "github.com/blocknextai/platform-api/internal/executions/infrastructure"
	nodeexecutionsInfrastructure "github.com/blocknextai/platform-api/internal/executions/infrastructure/nodeexecutions"
	taskclaimsInfrastructure "github.com/blocknextai/platform-api/internal/executions/infrastructure/taskclaims"
	taskexecutionsInfrastructure "github.com/blocknextai/platform-api/internal/executions/infrastructure/taskexecutions"
	executionsPresentation "github.com/blocknextai/platform-api/internal/executions/presentation"
	organizationsApplicationOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/application/organizationusers"
	workflowsApplicationWorkflows "github.com/blocknextai/platform-api/internal/workflows/application/workflows"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	DB                 *sql.DB
	TransactionManager database.TransactionManager

	OrganizationUserService organizationsApplicationOrganizationUsers.OrganizationUserService
	WorkflowService         workflowsApplicationWorkflows.WorkflowService
	UserService             accountApplicationUsers.UserService
	LinkedAccountService    accountApplicationLinkedAccounts.LinkedAccountService
}

type Module struct {
	TaskExecutionService taskexecutionsApplication.TaskExecutionService
	NodeExecutionService nodeexecutionsApplication.NodeExecutionService
	TaskClaimService     taskclaimsApplication.TaskClaimService

	handlers *executionsInfrastructure.Handlers
}

func NewModule(deps Dependencies) *Module {
	taskExecutionRepository := taskexecutionsInfrastructure.NewTaskExecutionRepository(deps.DB)
	taskClaimRepository := taskclaimsInfrastructure.NewTaskClaimRepository(deps.DB)
	nodeExecutionRepository := nodeexecutionsInfrastructure.NewNodeExecutionRepository(deps.DB)

	taskClaimService := taskclaimsApplication.NewTaskClaimService(taskClaimRepository, deps.TransactionManager)
	taskExecutionService := taskexecutionsApplication.NewTaskExecutionService(
		taskExecutionRepository,
		taskClaimService,
		deps.OrganizationUserService,
		deps.TransactionManager,
	)
	nodeExecutionService := nodeexecutionsApplication.NewNodeExecutionService(
		nodeExecutionRepository,
		deps.TransactionManager,
	)

	handlers := executionsInfrastructure.RegisterInfrastructure(executionsInfrastructure.RegisterInfrastructureDeps{
		TransactionManager: deps.TransactionManager,

		TaskExecutionRepository: taskExecutionRepository,
		NodeExecutionService:    nodeExecutionService,
		WorkflowService:         deps.WorkflowService,
		OrganizationUserService: deps.OrganizationUserService,
		UserService:             deps.UserService,
		LinkedAccountService:    deps.LinkedAccountService,
	})

	return &Module{
		TaskExecutionService: taskExecutionService,
		NodeExecutionService: nodeExecutionService,
		TaskClaimService:     taskClaimService,
		handlers:             handlers,
	}
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware) {
	executionsPresentation.RegisterPresentation(router, authMiddleware, m.handlers)
}
