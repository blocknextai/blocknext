package http

import (
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/rbac"
	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"
	wsApplicationHub "github.com/blocknextai/platform-api/internal/modules/ws/internal/application/hub"
	wsHTTPWebsockets "github.com/blocknextai/platform-api/internal/modules/ws/internal/http/websockets"
)

func RegisterRoutes(
	router fiber.Router,
	hub wsApplicationHub.HubService,
	authMiddleware *commonAuth.AuthMiddleware,
) {
	organizationsRouterGroup := router.Group("/organizations/:organizationId")

	organizationsRouterGroup.Get(
		"/ws",
		authMiddleware.AuthenticateWebSocket(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadOrganizationPermission),
		func(c fiber.Ctx) error {
			if !websocket.IsWebSocketUpgrade(c) {
				return fiber.ErrUpgradeRequired
			}
			return c.Next()
		},
		websocket.New(wsHTTPWebsockets.NewWebSocketHandler(hub)),
	)
}
