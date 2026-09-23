package http

import (
	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/rbac"
	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"
	apikeysContract "github.com/blocknextai/platform-api/internal/modules/apikeys/contract"
	taskRunnerHTTPTaskRunner "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/http/taskrunner"
	taskRunnerUseCases "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/usecases"
)

func RegisterRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	apiKeyMiddleware *commonAuth.APIKeyMiddleware[apikeysContract.Scope],
	useCases *taskRunnerUseCases.Services,
) {
	organizationTaskRunnerRouterGroup := router.Group("/organizations/:organizationId/task-runner")

	organizationTaskRunnerRouterGroup.Post(
		"/trigger",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.TriggerTaskPermission),
		taskRunnerHTTPTaskRunner.NewTriggerTaskHandler(useCases.Tasks),
	)

	organizationTaskRunnerRouterGroup.Post(
		"/rerun-all",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.RetryTaskPermission),
		taskRunnerHTTPTaskRunner.NewRerunAllHandler(useCases.Tasks),
	)

	organizationTaskRunnerRouterGroup.Post(
		"/rerun-failed",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.RetryTaskPermission),
		taskRunnerHTTPTaskRunner.NewRerunFailedHandler(useCases.Tasks),
	)

	organizationTaskRunnerRouterGroup.Post(
		"/cancel",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.CancelTaskPermission),
		taskRunnerHTTPTaskRunner.NewCancelTaskHandler(useCases.Tasks),
	)

	apiRouterGroup := router.Group("/task-runner/trigger")

	apiRouterGroup.Post(
		"/:workflowId",
		apiKeyMiddleware.Authenticate(),
		apiKeyMiddleware.RequireScope(apikeysContract.ScopeWorkflowsTrigger),
		taskRunnerHTTPTaskRunner.NewTriggerTaskWithAPIWorkflowHandler(useCases.Tasks),
	)
}
