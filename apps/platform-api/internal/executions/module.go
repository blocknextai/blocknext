package executions

import (
	"database/sql"
	apiKeysApplicationAPIKeys "github.com/blocknextai/platform-api/internal/apikeys/application/apikeys"

	"github.com/blocknextai/go-packages/database"
	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/account/application/linkedaccounts"
	accountApplicationUsers "github.com/blocknextai/platform-api/internal/account/application/users"
	"github.com/blocknextai/platform-api/internal/common/presentation/auth"
	nodeexecutionsApplication "github.com/blocknextai/platform-api/internal/executions/application/nodeexecutions"
	taskclaimsApplication "github.com/blocknextai/platform-api/internal/executions/application/taskclaims"
	taskexecutionsApplication "github.com/blocknextai/platform-api/internal/executions/application/taskexecutions"
	toolinvocationsApplication "github.com/blocknextai/platform-api/internal/executions/application/toolinvocations"
	"github.com/blocknextai/platform-api/internal/executions/domain/nodeexecutions"
	"github.com/blocknextai/platform-api/internal/executions/domain/taskclaims"
	"github.com/blocknextai/platform-api/internal/executions/domain/taskexecutions"
	"github.com/blocknextai/platform-api/internal/executions/domain/toolinvocations"
	executionsInfrastructure "github.com/blocknextai/platform-api/internal/executions/infrastructure"
	nodeexecutionsInfrastructure "github.com/blocknextai/platform-api/internal/executions/infrastructure/nodeexecutions"
	taskclaimsInfrastructure "github.com/blocknextai/platform-api/internal/executions/infrastructure/taskclaims"
	taskexecutionsInfrastructure "github.com/blocknextai/platform-api/internal/executions/infrastructure/taskexecutions"
	toolinvocationsInfrastructure "github.com/blocknextai/platform-api/internal/executions/infrastructure/toolinvocations"
	executionsPresentation "github.com/blocknextai/platform-api/internal/executions/presentation"
	organizationsApplicationOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/application/organizationusers"
	workflowsApplicationWorkflows "github.com/blocknextai/platform-api/internal/workflows/application/workflows"
	"github.com/gofiber/fiber/v3"
)

type ServicesDependencies struct {
	DB                 *sql.DB
	TransactionManager database.TransactionManager

	OrganizationUserService organizationsApplicationOrganizationUsers.OrganizationUserService
}

type Dependencies struct {
	DB                 *sql.DB
	TransactionManager database.TransactionManager

	OrganizationUserService organizationsApplicationOrganizationUsers.OrganizationUserService
	WorkflowService         workflowsApplicationWorkflows.WorkflowService
	APIKeyService           apiKeysApplicationAPIKeys.APIKeyService
	UserService             accountApplicationUsers.UserService
	LinkedAccountService    accountApplicationLinkedAccounts.LinkedAccountService
}

type Services struct {
	TaskExecutionService  taskexecutionsApplication.TaskExecutionService
	NodeExecutionService  nodeexecutionsApplication.NodeExecutionService
	TaskClaimService      taskclaimsApplication.TaskClaimService
	ToolInvocationService toolinvocationsApplication.ToolInvocationService
}

type Module struct {
	TaskExecutionService  taskexecutionsApplication.TaskExecutionService
	NodeExecutionService  nodeexecutionsApplication.NodeExecutionService
	TaskClaimService      taskclaimsApplication.TaskClaimService
	ToolInvocationService toolinvocationsApplication.ToolInvocationService

	handlers *executionsInfrastructure.Handlers
}

func NewServices(deps ServicesDependencies) *Services {
	return newServices(
		taskexecutionsInfrastructure.NewTaskExecutionRepository(deps.DB),
		taskclaimsInfrastructure.NewTaskClaimRepository(deps.DB),
		nodeexecutionsInfrastructure.NewNodeExecutionRepository(deps.DB),
		toolinvocationsInfrastructure.NewToolInvocationRepository(deps.DB),
		deps.OrganizationUserService,
		deps.TransactionManager,
	)
}

func newServices(
	taskExecutionRepository taskexecutions.TaskExecutionRepository,
	taskClaimRepository taskclaims.TaskClaimRepository,
	nodeExecutionRepository nodeexecutions.NodeExecutionRepository,
	toolInvocationRepository toolinvocations.ToolInvocationRepository,
	organizationUserService organizationsApplicationOrganizationUsers.OrganizationUserService,
	transactionManager database.TransactionManager,
) *Services {
	taskClaimService := taskclaimsApplication.NewTaskClaimService(taskClaimRepository, transactionManager)
	taskExecutionService := taskexecutionsApplication.NewTaskExecutionService(
		taskExecutionRepository,
		taskClaimService,
		organizationUserService,
		transactionManager,
	)
	nodeExecutionService := nodeexecutionsApplication.NewNodeExecutionService(
		nodeExecutionRepository,
		transactionManager,
	)

	toolInvocationService := toolinvocationsApplication.NewToolInvocationService(
		toolInvocationRepository,
		transactionManager,
	)

	return &Services{
		TaskExecutionService:  taskExecutionService,
		NodeExecutionService:  nodeExecutionService,
		TaskClaimService:      taskClaimService,
		ToolInvocationService: toolInvocationService,
	}
}

func NewModule(deps Dependencies) *Module {
	taskExecutionRepository := taskexecutionsInfrastructure.NewTaskExecutionRepository(deps.DB)
	taskClaimRepository := taskclaimsInfrastructure.NewTaskClaimRepository(deps.DB)
	nodeExecutionRepository := nodeexecutionsInfrastructure.NewNodeExecutionRepository(deps.DB)
	toolInvocationRepository := toolinvocationsInfrastructure.NewToolInvocationRepository(deps.DB)

	services := newServices(
		taskExecutionRepository,
		taskClaimRepository,
		nodeExecutionRepository,
		toolInvocationRepository,
		deps.OrganizationUserService,
		deps.TransactionManager,
	)

	handlers := executionsInfrastructure.RegisterInfrastructure(executionsInfrastructure.RegisterInfrastructureDeps{
		TransactionManager: deps.TransactionManager,

		TaskExecutionRepository:  taskExecutionRepository,
		ToolInvocationRepository: toolInvocationRepository,
		APIKeyService:            deps.APIKeyService,
		NodeExecutionService:     services.NodeExecutionService,
		WorkflowService:          deps.WorkflowService,
		OrganizationUserService:  deps.OrganizationUserService,
		UserService:              deps.UserService,
		LinkedAccountService:     deps.LinkedAccountService,
	})

	return &Module{
		TaskExecutionService:  services.TaskExecutionService,
		NodeExecutionService:  services.NodeExecutionService,
		TaskClaimService:      services.TaskClaimService,
		ToolInvocationService: services.ToolInvocationService,
		handlers:              handlers,
	}
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware) {
	executionsPresentation.RegisterPresentation(router, authMiddleware, m.handlers)
}
