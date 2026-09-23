package http

import (
	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/rbac"
	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"
	"github.com/blocknextai/platform-api/internal/modules/workflows/internal/application/generation/chat"
	"github.com/blocknextai/platform-api/internal/modules/workflows/internal/http/generation"
	workflowsHTTPWorkflows "github.com/blocknextai/platform-api/internal/modules/workflows/internal/http/workflows"
	workflowsUseCases "github.com/blocknextai/platform-api/internal/modules/workflows/internal/usecases"
)

func RegisterRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	chatService chat.ChatService,
	useCases *workflowsUseCases.Services,
) {
	workflowsRouterGroup := router.Group("/organizations/:organizationId/workflows")

	workflowsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadWorkflowPermission),
		workflowsHTTPWorkflows.NewGetAllWorkflowsHandler(useCases.Workflows),
	)

	workflowsRouterGroup.Post(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.CreateWorkflowPermission),
		workflowsHTTPWorkflows.NewCreateWorkflowHandler(useCases.Workflows),
	)

	workflowsRouterGroup.Get(
		"/:workflowId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadWorkflowPermission),
		workflowsHTTPWorkflows.NewGetWorkflowByIDHandler(useCases.Workflows),
	)

	workflowsRouterGroup.Patch(
		"/:workflowId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateWorkflowPermission),
		workflowsHTTPWorkflows.NewUpdateWorkflowHandler(useCases.Workflows),
	)

	workflowsRouterGroup.Delete(
		"/:workflowId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.DeleteWorkflowPermission),
		workflowsHTTPWorkflows.NewDeleteWorkflowHandler(useCases.Workflows),
	)

	workflowsRouterGroup.Post(
		"/:workflowId/duplicate",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.CreateWorkflowPermission),
		workflowsHTTPWorkflows.NewDuplicateWorkflowHandler(useCases.Workflows),
	)

	workflowsRouterGroup.Get(
		"/:workflowId/run",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadWorkflowPermission),
		workflowsHTTPWorkflows.NewGetWorkflowForRunHandler(useCases.Workflows),
	)

	generationRouterGroup := workflowsRouterGroup.Group("/generation/sessions")

	generationRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadWorkflowGenerationSessionPermission),
		generation.NewGetAllSessionsHandler(useCases.Generation),
	)
	generationRouterGroup.Post(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.CreateWorkflowGenerationSessionPermission),
		generation.NewCreateSessionHandler(useCases.Generation),
	)
	generationRouterGroup.Get(
		"/:sessionId/messages",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadWorkflowGenerationMessagePermission),
		generation.NewGetAllSessionMessagesHandler(useCases.Generation),
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
		generation.NewUpdateSessionHandler(useCases.Generation),
	)
	generationRouterGroup.Delete(
		"/:sessionId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.DeleteWorkflowGenerationSessionPermission),
		generation.NewDeleteSessionHandler(useCases.Generation),
	)
}
