package presentation

import (
	"github.com/blocknextai/go-packages/rbac"
	apiKeysDomain "github.com/blocknextai/platform-api/internal/apikeys/domain/apikeys"
	commonPresentationAuth "github.com/blocknextai/platform-api/internal/common/presentation/auth"
	taskRunnerInfrastructure "github.com/blocknextai/platform-api/internal/taskrunner/infrastructure"
	taskRunnerPresentationTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/presentation/taskrunner"
	"github.com/gofiber/fiber/v3"
)

func RegisterPresentation(
	router fiber.Router,
	authMiddleware *commonPresentationAuth.AuthMiddleware,
	apiKeyMiddleware *commonPresentationAuth.APIKeyMiddleware,
	handlers *taskRunnerInfrastructure.Handlers,
) {
	organizationTaskRunnerRouterGroup := router.Group("/organizations/:organizationId/task-runner")

	organizationTaskRunnerRouterGroup.Post(
		"/trigger",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.TriggerTaskPermission),
		taskRunnerPresentationTaskRunner.NewTriggerTaskHandler(handlers.TriggerTask),
	)

	organizationTaskRunnerRouterGroup.Post(
		"/rerun-all",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.RetryTaskPermission),
		taskRunnerPresentationTaskRunner.NewRerunAllHandler(handlers.RerunAll),
	)

	organizationTaskRunnerRouterGroup.Post(
		"/rerun-failed",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.RetryTaskPermission),
		taskRunnerPresentationTaskRunner.NewRerunFailedHandler(handlers.RerunFailed),
	)

	organizationTaskRunnerRouterGroup.Post(
		"/cancel",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.CancelTaskPermission),
		taskRunnerPresentationTaskRunner.NewCancelTaskHandler(handlers.CancelTask),
	)

	apiRouterGroup := router.Group("/task-runner/trigger")

	apiRouterGroup.Post(
		"/:workflowId",
		apiKeyMiddleware.Authenticate(),
		apiKeyMiddleware.RequireScope(apiKeysDomain.ScopeWorkflowsTrigger),
		taskRunnerPresentationTaskRunner.NewTriggerTaskWithAPIWorkflowHandler(handlers.TriggerTask),
	)
}
