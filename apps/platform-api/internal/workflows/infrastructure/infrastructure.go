package infrastructure

import (
	"github.com/blocknextai/go-packages/database"
	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/account/application/linkedaccounts"
	accountApplicationUsers "github.com/blocknextai/platform-api/internal/account/application/users"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/nodeengine/application/credentials"
	nodeEngineApplicationNodes "github.com/blocknextai/platform-api/internal/nodeengine/application/nodes"
	organizationsApplicationOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/application/organizationusers"
	"github.com/blocknextai/platform-api/internal/workflows/application/generation/sessions/createsession"
	"github.com/blocknextai/platform-api/internal/workflows/application/generation/sessions/deletesession"
	"github.com/blocknextai/platform-api/internal/workflows/application/generation/sessions/getallsessionmessages"
	"github.com/blocknextai/platform-api/internal/workflows/application/generation/sessions/getallsessions"
	"github.com/blocknextai/platform-api/internal/workflows/application/generation/sessions/updatesession"
	"github.com/blocknextai/platform-api/internal/workflows/application/workflows/createworkflow"
	"github.com/blocknextai/platform-api/internal/workflows/application/workflows/deleteworkflow"
	"github.com/blocknextai/platform-api/internal/workflows/application/workflows/duplicateworkflow"
	"github.com/blocknextai/platform-api/internal/workflows/application/workflows/getallworkflows"
	"github.com/blocknextai/platform-api/internal/workflows/application/workflows/getworkflowbyid"
	"github.com/blocknextai/platform-api/internal/workflows/application/workflows/getworkflowforrun"
	"github.com/blocknextai/platform-api/internal/workflows/application/workflows/updateworkflow"
	generationDomainMessages "github.com/blocknextai/platform-api/internal/workflows/domain/generation/messages"
	generationDomainSessions "github.com/blocknextai/platform-api/internal/workflows/domain/generation/sessions"
	"github.com/blocknextai/platform-api/internal/workflows/domain/workflows"
)

type Handlers struct {
	CreateWorkflow        cqrs.Handler[*createworkflow.CreateWorkflowCommand, *createworkflow.CreateWorkflowResponse]
	GetAllWorkflows       cqrs.Handler[*getallworkflows.GetAllWorkflowsQuery, *getallworkflows.GetAllWorkflowsResponse]
	GetWorkflowByID       cqrs.Handler[*getworkflowbyid.GetWorkflowByIDQuery, *getworkflowbyid.GetWorkflowByIDResponse]
	GetWorkflowForRun     cqrs.Handler[*getworkflowforrun.GetWorkflowForRunQuery, *getworkflowforrun.GetWorkflowForRunResponse]
	UpdateWorkflow        cqrs.Handler[*updateworkflow.UpdateWorkflowCommand, *updateworkflow.UpdateWorkflowResponse]
	DeleteWorkflow        cqrs.Handler[*deleteworkflow.DeleteWorkflowCommand, *deleteworkflow.DeleteWorkflowResponse]
	DuplicateWorkflow     cqrs.Handler[*duplicateworkflow.DuplicateWorkflowCommand, *duplicateworkflow.DuplicateWorkflowResponse]
	CreateSession         cqrs.Handler[*createsession.CreateSessionCommand, *createsession.CreateSessionResponse]
	GetAllSessions        cqrs.Handler[*getallsessions.GetAllSessionsQuery, *getallsessions.GetAllSessionsResponse]
	GetAllSessionMessages cqrs.Handler[*getallsessionmessages.GetAllSessionMessagesQuery, *getallsessionmessages.GetAllSessionMessagesResponse]
	UpdateSession         cqrs.Handler[*updatesession.UpdateSessionCommand, *updatesession.UpdateSessionResponse]
	DeleteSession         cqrs.Handler[*deletesession.DeleteSessionCommand, *deletesession.DeleteSessionResponse]

	// TODO: Violation!!!
	deps RegisterInfrastructureDeps
}

type RegisterInfrastructureDeps struct {
	TransactionManager database.TransactionManager

	WorkflowRepository      workflows.WorkflowRepository
	SessionRepository       generationDomainSessions.SessionRepository
	MessageRepository       generationDomainMessages.MessageRepository
	OrganizationUserService organizationsApplicationOrganizationUsers.OrganizationUserService
	UserService             accountApplicationUsers.UserService
	LinkedAccountService    accountApplicationLinkedAccounts.LinkedAccountService
	CredentialService       nodeEngineApplicationCredentials.CredentialService
	NodeService             nodeEngineApplicationNodes.NodeService
}

func RegisterInfrastructure(deps RegisterInfrastructureDeps) *Handlers {
	return &Handlers{
		CreateWorkflow:        cqrs.ValidationBehavior(createworkflow.New(deps.WorkflowRepository, deps.TransactionManager, deps.OrganizationUserService)),
		GetAllWorkflows:       cqrs.ValidationBehavior(getallworkflows.New(deps.WorkflowRepository, deps.OrganizationUserService, deps.UserService, deps.LinkedAccountService)),
		GetWorkflowByID:       cqrs.ValidationBehavior(getworkflowbyid.New(deps.WorkflowRepository, deps.OrganizationUserService, deps.UserService, deps.LinkedAccountService)),
		GetWorkflowForRun:     cqrs.ValidationBehavior(getworkflowforrun.New(deps.WorkflowRepository, deps.CredentialService, deps.NodeService)),
		UpdateWorkflow:        cqrs.ValidationBehavior(updateworkflow.New(deps.WorkflowRepository, deps.TransactionManager)),
		DeleteWorkflow:        cqrs.ValidationBehavior(deleteworkflow.New(deps.WorkflowRepository, deps.TransactionManager)),
		DuplicateWorkflow:     cqrs.ValidationBehavior(duplicateworkflow.New(deps.WorkflowRepository, deps.TransactionManager, deps.OrganizationUserService)),
		CreateSession:         cqrs.ValidationBehavior(createsession.New(deps.SessionRepository)),
		GetAllSessions:        cqrs.ValidationBehavior(getallsessions.New(deps.SessionRepository)),
		GetAllSessionMessages: cqrs.ValidationBehavior(getallsessionmessages.New(deps.SessionRepository, deps.MessageRepository)),
		UpdateSession:         cqrs.ValidationBehavior(updatesession.New(deps.SessionRepository)),
		DeleteSession:         cqrs.ValidationBehavior(deletesession.New(deps.SessionRepository, deps.MessageRepository, deps.TransactionManager)),

		// TODO: Violation!!!
		deps: deps,
	}
}
