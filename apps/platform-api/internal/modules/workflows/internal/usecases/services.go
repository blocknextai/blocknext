package usecases

import (
	"github.com/blocknextai/go-packages/database"
	accountContract "github.com/blocknextai/platform-api/internal/modules/account/contract"
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
	organizationsContract "github.com/blocknextai/platform-api/internal/modules/organizations/contract"
	generationDomainMessages "github.com/blocknextai/platform-api/internal/modules/workflows/internal/domain/generation/messages"
	generationDomainSessions "github.com/blocknextai/platform-api/internal/modules/workflows/internal/domain/generation/sessions"
	"github.com/blocknextai/platform-api/internal/modules/workflows/internal/domain/workflows"
	generationUseCases "github.com/blocknextai/platform-api/internal/modules/workflows/internal/usecases/generation"
	workflowsUseCases "github.com/blocknextai/platform-api/internal/modules/workflows/internal/usecases/workflows"
)

type Services struct {
	Generation *generationUseCases.Service
	Workflows  *workflowsUseCases.Service

	// TODO: Violation!!!
	deps ServiceDependencies
}

type ServiceDependencies struct {
	TransactionManager database.TransactionManager

	WorkflowRepository      workflows.WorkflowRepository
	SessionRepository       generationDomainSessions.SessionRepository
	MessageRepository       generationDomainMessages.MessageRepository
	OrganizationUserService organizationsContract.OrganizationUserService
	UserService             accountContract.UserService
	LinkedAccountService    accountContract.LinkedAccountService
	CredentialService       nodeengineContract.CredentialService
	NodeService             nodeengineContract.NodeService
}

func NewServices(deps ServiceDependencies) *Services {
	return &Services{
		Generation: generationUseCases.NewService(deps.SessionRepository, deps.MessageRepository, deps.TransactionManager),
		Workflows:  workflowsUseCases.NewService(deps.WorkflowRepository, deps.TransactionManager, deps.OrganizationUserService, deps.UserService, deps.LinkedAccountService, deps.CredentialService, deps.NodeService),

		// TODO: Violation!!!
		deps: deps,
	}
}
