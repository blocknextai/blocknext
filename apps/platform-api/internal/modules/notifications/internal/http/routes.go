package http

import (
	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/rbac"
	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"
	"github.com/blocknextai/platform-api/internal/modules/notifications/internal/http/notifications"
	notificationsUseCases "github.com/blocknextai/platform-api/internal/modules/notifications/internal/usecases"
)

func RegisterUserNotificationsRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	useCases *notificationsUseCases.Services,
) {
	notificationsRouterGroup := router.Group("/users/me/notifications")

	notificationsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadNotificationPermission),
		notifications.NewGetAllUserNotificationsHandler(useCases.Notificationrecipients),
	)

	notificationsRouterGroup.Get(
		"/count",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadNotificationPermission),
		notifications.NewGetUserNotificationCountsHandler(useCases.Notificationrecipients),
	)

	notificationsRouterGroup.Post(
		"/seen",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.UpdateNotificationPermission),
		notifications.NewMarkAllUserNotificationsSeenHandler(useCases.Notificationrecipients),
	)

	notificationsRouterGroup.Post(
		"/read-all",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.UpdateNotificationPermission),
		notifications.NewMarkAllUserNotificationsReadHandler(useCases.Notificationrecipients),
	)

	notificationsRouterGroup.Patch(
		"/:recipientId/read",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.UpdateNotificationPermission),
		notifications.NewMarkNotificationReadHandler(useCases.Notificationrecipients),
	)

	notificationsRouterGroup.Delete(
		"/:recipientId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.DeleteNotificationPermission),
		notifications.NewDeleteNotificationHandler(useCases.Notificationrecipients),
	)
}

func RegisterOrganizationNotificationsRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	useCases *notificationsUseCases.Services,
) {
	notificationsRouterGroup := router.Group("/organizations/:organizationId/notifications")

	notificationsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadNotificationPermission),
		notifications.NewGetAllOrganizationNotificationsHandler(useCases.Notificationrecipients),
	)

	notificationsRouterGroup.Get(
		"/count",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadNotificationPermission),
		notifications.NewGetOrganizationNotificationCountsHandler(useCases.Notificationrecipients),
	)

	notificationsRouterGroup.Post(
		"/seen",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateNotificationPermission),
		notifications.NewMarkAllOrganizationNotificationsSeenHandler(useCases.Notificationrecipients),
	)

	notificationsRouterGroup.Post(
		"/read-all",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateNotificationPermission),
		notifications.NewMarkAllOrganizationNotificationsReadHandler(useCases.Notificationrecipients),
	)

	notificationsRouterGroup.Patch(
		"/:recipientId/read",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateNotificationPermission),
		notifications.NewMarkNotificationReadHandler(useCases.Notificationrecipients),
	)

	notificationsRouterGroup.Delete(
		"/:recipientId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.DeleteNotificationPermission),
		notifications.NewDeleteNotificationHandler(useCases.Notificationrecipients),
	)
}

func RegisterRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	useCases *notificationsUseCases.Services,
) {
	RegisterUserNotificationsRoutes(router, authMiddleware, useCases)
	RegisterOrganizationNotificationsRoutes(router, authMiddleware, useCases)
}
