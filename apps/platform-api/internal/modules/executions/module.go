package executions

import (
	"database/sql"

	"github.com/blocknextai/platform-api/internal/realtime"

	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/common/auth"
	accountContract "github.com/blocknextai/platform-api/internal/modules/account/contract"
	apikeysContract "github.com/blocknextai/platform-api/internal/modules/apikeys/contract"
	nodeexecutionsApplication "github.com/blocknextai/platform-api/internal/modules/executions/internal/application/nodeexecutions"
	taskclaimsApplication "github.com/blocknextai/platform-api/internal/modules/executions/internal/application/taskclaims"
	taskexecutionsApplication "github.com/blocknextai/platform-api/internal/modules/executions/internal/application/taskexecutions"
	toolinvocationsApplication "github.com/blocknextai/platform-api/internal/modules/executions/internal/application/toolinvocations"
	"github.com/blocknextai/platform-api/internal/modules/executions/internal/domain/nodeexecutions"
	"github.com/blocknextai/platform-api/internal/modules/executions/internal/domain/taskclaims"
	"github.com/blocknextai/platform-api/internal/modules/executions/internal/domain/taskexecutions"
	"github.com/blocknextai/platform-api/internal/modules/executions/internal/domain/toolinvocations"
	executionsHTTP "github.com/blocknextai/platform-api/internal/modules/executions/internal/http"
	nodeexecutionsPostgres "github.com/blocknextai/platform-api/internal/modules/executions/internal/postgres/nodeexecutions"
	taskclaimsPostgres "github.com/blocknextai/platform-api/internal/modules/executions/internal/postgres/taskclaims"
	taskexecutionsPostgres "github.com/blocknextai/platform-api/internal/modules/executions/internal/postgres/taskexecutions"
	toolinvocationsPostgres "github.com/blocknextai/platform-api/internal/modules/executions/internal/postgres/toolinvocations"
	executionsUseCases "github.com/blocknextai/platform-api/internal/modules/executions/internal/usecases"
	organizationsContract "github.com/blocknextai/platform-api/internal/modules/organizations/contract"
	workflowsContract "github.com/blocknextai/platform-api/internal/modules/workflows/contract"
)

type ServicesDependencies struct {
	DB                 *sql.DB
	TransactionManager database.TransactionManager
	Broadcaster        realtime.Broadcaster

	OrganizationUserService organizationsContract.OrganizationUserService
}

type Dependencies struct {
	DB                 *sql.DB
	TransactionManager database.TransactionManager
	Broadcaster        realtime.Broadcaster

	OrganizationUserService organizationsContract.OrganizationUserService
	WorkflowService         workflowsContract.WorkflowService
	APIKeyService           apikeysContract.APIKeyService
	UserService             accountContract.UserService
	LinkedAccountService    accountContract.LinkedAccountService
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

	useCases *executionsUseCases.Services
}

func NewServices(deps ServicesDependencies) *Services {
	return newServices(
		taskexecutionsPostgres.NewTaskExecutionRepository(deps.DB),
		taskclaimsPostgres.NewTaskClaimRepository(deps.DB),
		nodeexecutionsPostgres.NewNodeExecutionRepository(deps.DB),
		toolinvocationsPostgres.NewToolInvocationRepository(deps.DB),
		deps.OrganizationUserService,
		deps.TransactionManager,
		deps.Broadcaster,
	)
}

func newServices(
	taskExecutionRepository taskexecutions.TaskExecutionRepository,
	taskClaimRepository taskclaims.TaskClaimRepository,
	nodeExecutionRepository nodeexecutions.NodeExecutionRepository,
	toolInvocationRepository toolinvocations.ToolInvocationRepository,
	organizationUserService organizationsContract.OrganizationUserService,
	transactionManager database.TransactionManager,
	broadcaster realtime.Broadcaster,
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
		broadcaster,
	)

	return &Services{
		TaskExecutionService:  taskExecutionService,
		NodeExecutionService:  nodeExecutionService,
		TaskClaimService:      taskClaimService,
		ToolInvocationService: toolInvocationService,
	}
}

func NewModule(deps Dependencies) *Module {
	taskExecutionRepository := taskexecutionsPostgres.NewTaskExecutionRepository(deps.DB)
	taskClaimRepository := taskclaimsPostgres.NewTaskClaimRepository(deps.DB)
	nodeExecutionRepository := nodeexecutionsPostgres.NewNodeExecutionRepository(deps.DB)
	toolInvocationRepository := toolinvocationsPostgres.NewToolInvocationRepository(deps.DB)

	services := newServices(
		taskExecutionRepository,
		taskClaimRepository,
		nodeExecutionRepository,
		toolInvocationRepository,
		deps.OrganizationUserService,
		deps.TransactionManager,
		deps.Broadcaster,
	)

	useCases := executionsUseCases.NewServices(executionsUseCases.ServiceDependencies{
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
		useCases:              useCases,
	}
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware) {
	executionsHTTP.RegisterRoutes(router, authMiddleware, m.useCases)
}
