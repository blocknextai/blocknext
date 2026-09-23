package workflows

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/common/auth"
	"github.com/blocknextai/platform-api/internal/config"
	accountContract "github.com/blocknextai/platform-api/internal/modules/account/contract"
	llmContract "github.com/blocknextai/platform-api/internal/modules/llm/contract"
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
	organizationsContract "github.com/blocknextai/platform-api/internal/modules/organizations/contract"
	generationChat "github.com/blocknextai/platform-api/internal/modules/workflows/internal/application/generation/chat"
	generationCredentialSchema "github.com/blocknextai/platform-api/internal/modules/workflows/internal/application/generation/credentialschema"
	generationNodeSchema "github.com/blocknextai/platform-api/internal/modules/workflows/internal/application/generation/nodeschema"
	generationTriggerVariables "github.com/blocknextai/platform-api/internal/modules/workflows/internal/application/generation/triggervariables"
	workflowsApplicationWorkflows "github.com/blocknextai/platform-api/internal/modules/workflows/internal/application/workflows"
	workflowsHTTP "github.com/blocknextai/platform-api/internal/modules/workflows/internal/http"
	generationInfraMessages "github.com/blocknextai/platform-api/internal/modules/workflows/internal/postgres/generation/messages"
	generationInfraSessions "github.com/blocknextai/platform-api/internal/modules/workflows/internal/postgres/generation/sessions"
	workflowsPostgresWorkflows "github.com/blocknextai/platform-api/internal/modules/workflows/internal/postgres/workflows"
	workflowsUseCases "github.com/blocknextai/platform-api/internal/modules/workflows/internal/usecases"
)

type Dependencies struct {
	DB                 *sql.DB
	TransactionManager database.TransactionManager

	WorkflowsOptions config.WorkflowsOptions

	GenerationProvider          llmContract.Provider
	NodeEngineNodeService       nodeengineContract.NodeService
	NodeEngineCredentialService nodeengineContract.CredentialService
	NodeEngineAdapterService    nodeengineContract.AdapterService
	OrganizationUserService     organizationsContract.OrganizationUserService
	UserService                 accountContract.UserService
	LinkedAccountService        accountContract.LinkedAccountService
}

type Module struct {
	WorkflowService workflowsApplicationWorkflows.WorkflowService

	chatService generationChat.ChatService
	useCases    *workflowsUseCases.Services
}

func NewModule(deps Dependencies) (*Module, error) {
	workflowRepository := workflowsPostgresWorkflows.NewWorkflowRepository(deps.DB)
	sessionRepository := generationInfraSessions.NewSessionRepository(deps.DB)
	messageRepository := generationInfraMessages.NewMessageRepository(deps.DB)

	workflowService := workflowsApplicationWorkflows.NewWorkflowService(workflowRepository)

	var chatService generationChat.ChatService
	if deps.GenerationProvider != nil {
		nodeSchemaContextBuilder := generationNodeSchema.NewNodeSchemaContextBuilder(deps.NodeEngineNodeService)
		credentialSchemaContextBuilder := generationCredentialSchema.NewCredentialSchemaContextBuilder(deps.NodeEngineCredentialService)
		triggerVariablesContextBuilder := generationTriggerVariables.NewTriggerVariablesContextBuilder(deps.NodeEngineAdapterService)

		chatService = generationChat.NewChatService(
			deps.WorkflowsOptions.Generation.SystemInstruction,
			sessionRepository,
			messageRepository,
			deps.GenerationProvider,
			nodeSchemaContextBuilder,
			credentialSchemaContextBuilder,
			triggerVariablesContextBuilder,
		)
	}

	useCases := workflowsUseCases.NewServices(workflowsUseCases.ServiceDependencies{
		TransactionManager: deps.TransactionManager,

		WorkflowRepository:      workflowRepository,
		SessionRepository:       sessionRepository,
		MessageRepository:       messageRepository,
		OrganizationUserService: deps.OrganizationUserService,
		UserService:             deps.UserService,
		LinkedAccountService:    deps.LinkedAccountService,
		CredentialService:       deps.NodeEngineCredentialService,
		NodeService:             deps.NodeEngineNodeService,
	})

	return &Module{
		WorkflowService: workflowService,
		chatService:     chatService,
		useCases:        useCases,
	}, nil
}

type Services struct {
	WorkflowService workflowsApplicationWorkflows.WorkflowService
}

type ServicesDependencies struct {
	DB *sql.DB
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		WorkflowService: workflowsApplicationWorkflows.NewWorkflowService(
			workflowsPostgresWorkflows.NewWorkflowRepository(deps.DB),
		),
	}
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware) {
	workflowsHTTP.RegisterRoutes(router, authMiddleware, m.chatService, m.useCases)
}
