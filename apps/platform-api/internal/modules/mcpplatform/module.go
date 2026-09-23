package mcpplatform

import (
	accountContract "github.com/blocknextai/platform-api/internal/modules/account/contract"
	credentialsContract "github.com/blocknextai/platform-api/internal/modules/credentials/contract"
	executionsContract "github.com/blocknextai/platform-api/internal/modules/executions/contract"
	mcpContract "github.com/blocknextai/platform-api/internal/modules/mcp/contract"
	mcpPlatformApplicationServer "github.com/blocknextai/platform-api/internal/modules/mcpplatform/internal/application/server"
	organizationsContract "github.com/blocknextai/platform-api/internal/modules/organizations/contract"
	triggersContract "github.com/blocknextai/platform-api/internal/modules/triggers/contract"
	workflowsContract "github.com/blocknextai/platform-api/internal/modules/workflows/contract"
)

const (
	serverVersion = "0.0.1"
)

type Dependencies struct {
	UserService             accountContract.UserService
	OrganizationService     organizationsContract.OrganizationService
	OrganizationUserService organizationsContract.OrganizationUserService
	WorkflowService         workflowsContract.WorkflowService
	TriggerService          triggersContract.TriggerService
	TaskExecutionService    executionsContract.TaskExecutionService
	ToolInvocationService   executionsContract.ToolInvocationService
	CredentialService       credentialsContract.CredentialService
}

type Module struct {
	ServerProvider mcpContract.ServerProvider
}

func NewModule(deps Dependencies) *Module {
	return &Module{
		ServerProvider: mcpPlatformApplicationServer.NewServerProvider(mcpPlatformApplicationServer.Dependencies{
			Version: serverVersion,

			UserService:             deps.UserService,
			OrganizationService:     deps.OrganizationService,
			OrganizationUserService: deps.OrganizationUserService,
			WorkflowService:         deps.WorkflowService,
			TriggerService:          deps.TriggerService,
			TaskExecutionService:    deps.TaskExecutionService,
			ToolInvocationService:   deps.ToolInvocationService,
			CredentialService:       deps.CredentialService,
		}),
	}
}
