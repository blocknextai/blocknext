package presentation

import (
	"github.com/blocknextai/go-packages/rbac"
	commonPresentationAuth "github.com/blocknextai/platform-api/internal/common/presentation/auth"
	notificationsInfrastructure "github.com/blocknextai/platform-api/internal/notifications/infrastructure"
	"github.com/blocknextai/platform-api/internal/notifications/presentation/notifications"
	"github.com/gofiber/fiber/v3"
)

func RegisterUserNotificationsPresentation(
	router fiber.Router,
	authMiddleware *commonPresentationAuth.AuthMiddleware,
	handlers *notificationsInfrastructure.Handlers,
) {
	notificationsRouterGroup := router.Group("/users/me/notifications")

	notificationsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadNotificationPermission),
		notifications.NewGetAllUserNotificationsHandler(handlers.GetAllNotifications),
	)

	notificationsRouterGroup.Get(
		"/count",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadNotificationPermission),
		notifications.NewGetUserNotificationCountsHandler(handlers.GetNotificationCounts),
	)

	notificationsRouterGroup.Post(
		"/seen",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.UpdateNotificationPermission),
		notifications.NewMarkAllUserNotificationsSeenHandler(handlers.MarkAllNotificationsSeen),
	)

	notificationsRouterGroup.Post(
		"/read-all",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.UpdateNotificationPermission),
		notifications.NewMarkAllUserNotificationsReadHandler(handlers.MarkAllNotificationsRead),
	)

	notificationsRouterGroup.Patch(
		"/:recipientId/read",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.UpdateNotificationPermission),
		notifications.NewMarkNotificationReadHandler(handlers.MarkNotificationRead),
	)

	notificationsRouterGroup.Delete(
		"/:recipientId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.DeleteNotificationPermission),
		notifications.NewDeleteNotificationHandler(handlers.DeleteNotification),
	)
}

func RegisterOrganizationNotificationsPresentation(
	router fiber.Router,
	authMiddleware *commonPresentationAuth.AuthMiddleware,
	handlers *notificationsInfrastructure.Handlers,
) {
	notificationsRouterGroup := router.Group("/organizations/:organizationId/notifications")

	notificationsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadNotificationPermission),
		notifications.NewGetAllOrganizationNotificationsHandler(handlers.GetAllNotifications),
	)

	notificationsRouterGroup.Get(
		"/count",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadNotificationPermission),
		notifications.NewGetOrganizationNotificationCountsHandler(handlers.GetNotificationCounts),
	)

	notificationsRouterGroup.Post(
		"/seen",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateNotificationPermission),
		notifications.NewMarkAllOrganizationNotificationsSeenHandler(handlers.MarkAllNotificationsSeen),
	)

	notificationsRouterGroup.Post(
		"/read-all",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateNotificationPermission),
		notifications.NewMarkAllOrganizationNotificationsReadHandler(handlers.MarkAllNotificationsRead),
	)

	notificationsRouterGroup.Patch(
		"/:recipientId/read",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateNotificationPermission),
		notifications.NewMarkNotificationReadHandler(handlers.MarkNotificationRead),
	)

	notificationsRouterGroup.Delete(
		"/:recipientId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.DeleteNotificationPermission),
		notifications.NewDeleteNotificationHandler(handlers.DeleteNotification),
	)
}

func RegisterPresentation(
	router fiber.Router,
	authMiddleware *commonPresentationAuth.AuthMiddleware,
	handlers *notificationsInfrastructure.Handlers,
) {
	RegisterUserNotificationsPresentation(router, authMiddleware, handlers)
	RegisterOrganizationNotificationsPresentation(router, authMiddleware, handlers)
}
