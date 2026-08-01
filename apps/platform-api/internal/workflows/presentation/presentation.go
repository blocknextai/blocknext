package presentation

import (
	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/rbac"
	commonPresentationAuth "github.com/blocknextai/platform-api/internal/common/presentation/auth"
	"github.com/blocknextai/platform-api/internal/workflows/application/generation/chat"
	workflowsInfrastructure "github.com/blocknextai/platform-api/internal/workflows/infrastructure"
	"github.com/blocknextai/platform-api/internal/workflows/presentation/generation"
	workflowsPresentationWorkflows "github.com/blocknextai/platform-api/internal/workflows/presentation/workflows"
)

func RegisterPresentation(
	router fiber.Router,
	authMiddleware *commonPresentationAuth.AuthMiddleware,
	chatService chat.ChatService,
	handlers *workflowsInfrastructure.Handlers,
) {
	workflowsRouterGroup := router.Group("/organizations/:organizationId/workflows")

	workflowsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadWorkflowPermission),
		workflowsPresentationWorkflows.NewGetAllWorkflowsHandler(handlers.GetAllWorkflows),
	)

	workflowsRouterGroup.Post(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.CreateWorkflowPermission),
		workflowsPresentationWorkflows.NewCreateWorkflowHandler(handlers.CreateWorkflow),
	)

	workflowsRouterGroup.Get(
		"/:workflowId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadWorkflowPermission),
		workflowsPresentationWorkflows.NewGetWorkflowByIDHandler(handlers.GetWorkflowByID),
	)

	workflowsRouterGroup.Patch(
		"/:workflowId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateWorkflowPermission),
		workflowsPresentationWorkflows.NewUpdateWorkflowHandler(handlers.UpdateWorkflow),
	)

	workflowsRouterGroup.Delete(
		"/:workflowId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.DeleteWorkflowPermission),
		workflowsPresentationWorkflows.NewDeleteWorkflowHandler(handlers.DeleteWorkflow),
	)

	workflowsRouterGroup.Post(
		"/:workflowId/duplicate",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.CreateWorkflowPermission),
		workflowsPresentationWorkflows.NewDuplicateWorkflowHandler(handlers.DuplicateWorkflow),
	)

	workflowsRouterGroup.Get(
		"/:workflowId/run",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadWorkflowPermission),
		workflowsPresentationWorkflows.NewGetWorkflowForRunHandler(handlers.GetWorkflowForRun),
	)

	generationRouterGroup := workflowsRouterGroup.Group("/generation/sessions")

	generationRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadWorkflowGenerationSessionPermission),
		generation.NewGetAllSessionsHandler(handlers.GetAllSessions),
	)
	generationRouterGroup.Post(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.CreateWorkflowGenerationSessionPermission),
		generation.NewCreateSessionHandler(handlers.CreateSession),
	)
	generationRouterGroup.Get(
		"/:sessionId/messages",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadWorkflowGenerationMessagePermission),
		generation.NewGetAllSessionMessagesHandler(handlers.GetAllSessionMessages),
	)
	generationRouterGroup.Post(
		"/:sessionId/messages",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.CreateWorkflowGenerationMessagePermission),
		generation.NewSendMessageHandler(chatService),
	)
	generationRouterGroup.Patch(
		"/:sessionId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateWorkflowGenerationSessionPermission),
		generation.NewUpdateSessionHandler(handlers.UpdateSession),
	)
	generationRouterGroup.Delete(
		"/:sessionId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.DeleteWorkflowGenerationSessionPermission),
		generation.NewDeleteSessionHandler(handlers.DeleteSession),
	)
}
