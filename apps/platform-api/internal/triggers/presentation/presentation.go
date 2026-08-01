package presentation

import (
	"github.com/blocknextai/go-packages/rbac"
	commonPresentationAuth "github.com/blocknextai/platform-api/internal/common/presentation/auth"
	triggersInfrastructure "github.com/blocknextai/platform-api/internal/triggers/infrastructure"
	triggers "github.com/blocknextai/platform-api/internal/triggers/presentation/triggers"
	"github.com/gofiber/fiber/v3"
)

func RegisterPresentation(
	router fiber.Router,
	authMiddleware *commonPresentationAuth.AuthMiddleware,
	handlers *triggersInfrastructure.Handlers,
) {
	organizationTriggersRouterGroup := router.Group("/organizations/:organizationId/triggers")

	organizationTriggersRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadTriggersPermission),
		triggers.NewGetAllTriggersHandler(handlers.GetAllTriggers),
	)

	organizationTriggersRouterGroup.Patch(
		"/:triggerId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateTriggersPermission),
		triggers.NewUpdateTriggerHandler(handlers.UpdateTrigger),
	)

	organizationTriggersRouterGroup.Post(
		"/:triggerId/regenerate-webhook-token",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateTriggersPermission),
		triggers.NewRegenerateWebhookTokenHandler(handlers.RegenerateWebhookToken),
	)

	organizationTriggersRouterGroup.Delete(
		"/:triggerId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.DeleteTriggersPermission),
		triggers.NewDeleteTriggerHandler(handlers.DeleteTrigger),
	)
}
