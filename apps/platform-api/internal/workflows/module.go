package workflows

import (
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/account/application/linkedaccounts"
	accountApplicationUsers "github.com/blocknextai/platform-api/internal/account/application/users"
	"github.com/blocknextai/platform-api/internal/common/presentation/auth"
	"github.com/blocknextai/platform-api/internal/config"
	"github.com/blocknextai/platform-api/internal/llm/streamingchat"
	nodeEngineApplicationAdapters "github.com/blocknextai/platform-api/internal/nodeengine/application/adapters"
	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/nodeengine/application/credentials"
	nodeEngineApplicationNodes "github.com/blocknextai/platform-api/internal/nodeengine/application/nodes"
	organizationsApplicationOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/application/organizationusers"
	generationChat "github.com/blocknextai/platform-api/internal/workflows/application/generation/chat"
	generationCredentialSchema "github.com/blocknextai/platform-api/internal/workflows/application/generation/credentialschema"
	generationNodeSchema "github.com/blocknextai/platform-api/internal/workflows/application/generation/nodeschema"
	generationTriggerVariables "github.com/blocknextai/platform-api/internal/workflows/application/generation/triggervariables"
	workflowsApplicationWorkflows "github.com/blocknextai/platform-api/internal/workflows/application/workflows"
	workflowsInfrastructure "github.com/blocknextai/platform-api/internal/workflows/infrastructure"
	generationInfraMessages "github.com/blocknextai/platform-api/internal/workflows/infrastructure/generation/messages"
	generationInfraSessions "github.com/blocknextai/platform-api/internal/workflows/infrastructure/generation/sessions"
	workflowsInfrastructureWorkflows "github.com/blocknextai/platform-api/internal/workflows/infrastructure/workflows"
	workflowsPresentation "github.com/blocknextai/platform-api/internal/workflows/presentation"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	DB                 *sql.DB
	TransactionManager database.TransactionManager

	WorkflowsOptions config.WorkflowsOptions

	GenerationProvider          streamingchat.Provider
	NodeEngineNodeService       nodeEngineApplicationNodes.NodeService
	NodeEngineCredentialService nodeEngineApplicationCredentials.CredentialService
	NodeEngineAdapterService    nodeEngineApplicationAdapters.AdapterService
	OrganizationUserService     organizationsApplicationOrganizationUsers.OrganizationUserService
	UserService                 accountApplicationUsers.UserService
	LinkedAccountService        accountApplicationLinkedAccounts.LinkedAccountService
}

type Module struct {
	WorkflowService workflowsApplicationWorkflows.WorkflowService

	chatService generationChat.ChatService
	handlers    *workflowsInfrastructure.Handlers
}

func NewModule(deps Dependencies) (*Module, error) {
	workflowRepository := workflowsInfrastructureWorkflows.NewWorkflowRepository(deps.DB)
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

	handlers := workflowsInfrastructure.RegisterInfrastructure(workflowsInfrastructure.RegisterInfrastructureDeps{
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
		handlers:        handlers,
	}, nil
}

// TODO: Violation!!!

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware) {
	workflowsPresentation.RegisterPresentation(router, authMiddleware, m.chatService, m.handlers)
}
