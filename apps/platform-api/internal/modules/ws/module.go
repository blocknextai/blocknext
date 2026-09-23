package ws

import (
	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/platform-api/internal/common/auth"
	wsApplicationHub "github.com/blocknextai/platform-api/internal/modules/ws/internal/application/hub"
	wsHTTP "github.com/blocknextai/platform-api/internal/modules/ws/internal/http"
	"github.com/blocknextai/platform-api/internal/realtime"
)

type Dependencies struct {
	Broadcaster realtime.Broadcaster

	MaxConnectionsPerRoom int
}

type Module struct {
	hub wsApplicationHub.HubService
}

func NewModule(deps Dependencies) *Module {
	hub := wsApplicationHub.NewHubService(deps.Broadcaster, deps.MaxConnectionsPerRoom)
	return &Module{hub: hub}
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware) {
	wsHTTP.RegisterRoutes(router, m.hub, authMiddleware)
}

func (m *Module) Shutdown() {
	m.hub.Shutdown()
}
