package http

import (
	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/rbac"
	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"
	executionsHTTPTaskExecutions "github.com/blocknextai/platform-api/internal/modules/executions/internal/http/taskexecutions"
	executionsHTTPToolInvocations "github.com/blocknextai/platform-api/internal/modules/executions/internal/http/toolinvocations"
	executionsUseCases "github.com/blocknextai/platform-api/internal/modules/executions/internal/usecases"
)

func RegisterRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	useCases *executionsUseCases.Services,
) {
	organizationRouterGroup := router.Group("/organizations/:organizationId")

	executionsRouterGroup := organizationRouterGroup.Group("/executions")
	executionsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadTaskExecutionPermission),
		executionsHTTPTaskExecutions.NewGetAllExecutionsHandler(useCases.TaskExecutions),
	)

	executionsRouterGroup.Delete(
		"/bulk",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.DeleteTaskExecutionPermission),
		executionsHTTPTaskExecutions.NewBulkDeleteTaskExecutionsHandler(useCases.TaskExecutions),
	)

	executionsRouterGroup.Get(
		"/:executionId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadTaskExecutionPermission),
		executionsHTTPTaskExecutions.NewGetExecutionByIDHandler(useCases.TaskExecutions),
	)

	executionsRouterGroup.Delete(
		"/:executionId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.DeleteTaskExecutionPermission),
		executionsHTTPTaskExecutions.NewDeleteExecutionHandler(useCases.TaskExecutions),
	)

	toolInvocationsRouterGroup := organizationRouterGroup.Group("/tool-invocations")
	toolInvocationsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadTaskExecutionPermission),
		executionsHTTPToolInvocations.NewGetAllToolInvocationsHandler(useCases.Toolinvocations),
	)

	toolInvocationsRouterGroup.Get(
		"/:toolInvocationId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadTaskExecutionPermission),
		executionsHTTPToolInvocations.NewGetToolInvocationByIDHandler(useCases.Toolinvocations),
	)
}
