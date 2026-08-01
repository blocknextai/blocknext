package presentation

import (
	"github.com/blocknextai/go-packages/rbac"
	commonPresentationAuth "github.com/blocknextai/platform-api/internal/common/presentation/auth"
	wsApplicationHub "github.com/blocknextai/platform-api/internal/ws/application/hub"
	wsPresentationWebsockets "github.com/blocknextai/platform-api/internal/ws/presentation/websockets"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

func RegisterPresentation(
	router fiber.Router,
	hub wsApplicationHub.HubService,
	authMiddleware *commonPresentationAuth.AuthMiddleware,
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
		websocket.New(wsPresentationWebsockets.NewWebSocketHandler(hub)),
	)
}
