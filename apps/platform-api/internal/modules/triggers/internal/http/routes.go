package http

import (
	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/rbac"
	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"
	triggers "github.com/blocknextai/platform-api/internal/modules/triggers/internal/http/triggers"
	triggersUseCases "github.com/blocknextai/platform-api/internal/modules/triggers/internal/usecases"
)

func RegisterRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	useCases *triggersUseCases.Services,
) {
	organizationTriggersRouterGroup := router.Group("/organizations/:organizationId/triggers")

	organizationTriggersRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadTriggersPermission),
		triggers.NewGetAllTriggersHandler(useCases.Triggers),
	)

	organizationTriggersRouterGroup.Patch(
		"/:triggerId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateTriggersPermission),
		triggers.NewUpdateTriggerHandler(useCases.Triggers),
	)

	organizationTriggersRouterGroup.Post(
		"/:triggerId/regenerate-webhook-token",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateTriggersPermission),
		triggers.NewRegenerateWebhookTokenHandler(useCases.Triggers),
	)

	organizationTriggersRouterGroup.Delete(
		"/:triggerId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.DeleteTriggersPermission),
		triggers.NewDeleteTriggerHandler(useCases.Triggers),
	)
}
