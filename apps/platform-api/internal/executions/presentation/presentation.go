package presentation

import (
	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/rbac"
	commonPresentationAuth "github.com/blocknextai/platform-api/internal/common/presentation/auth"
	executionsInfrastructure "github.com/blocknextai/platform-api/internal/executions/infrastructure"
	executionsPresentationTaskExecutions "github.com/blocknextai/platform-api/internal/executions/presentation/taskexecutions"
)

func RegisterPresentation(
	router fiber.Router,
	authMiddleware *commonPresentationAuth.AuthMiddleware,
	handlers *executionsInfrastructure.Handlers,
) {
	organizationRouterGroup := router.Group("/organizations/:organizationId")

	executionsRouterGroup := organizationRouterGroup.Group("/executions")
	executionsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadTaskExecutionPermission),
		executionsPresentationTaskExecutions.NewGetAllExecutionsHandler(handlers.GetAllTaskExecutions),
	)

	executionsRouterGroup.Delete(
		"/bulk",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.DeleteTaskExecutionPermission),
		executionsPresentationTaskExecutions.NewBulkDeleteTaskExecutionsHandler(handlers.BulkDeleteTaskExecutions),
	)

	executionsRouterGroup.Get(
		"/:executionId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadTaskExecutionPermission),
		executionsPresentationTaskExecutions.NewGetExecutionByIDHandler(handlers.GetTaskExecutionByID),
	)

	executionsRouterGroup.Delete(
		"/:executionId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.DeleteTaskExecutionPermission),
		executionsPresentationTaskExecutions.NewDeleteExecutionHandler(handlers.DeleteTaskExecution),
	)
}
