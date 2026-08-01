package ws

import (
	"github.com/blocknextai/platform-api/internal/common/presentation/auth"
	"github.com/blocknextai/platform-api/internal/realtime"
	wsApplicationHub "github.com/blocknextai/platform-api/internal/ws/application/hub"
	wsPresentation "github.com/blocknextai/platform-api/internal/ws/presentation"
	"github.com/gofiber/fiber/v3"
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
	wsPresentation.RegisterPresentation(router, m.hub, authMiddleware)
}

func (m *Module) Shutdown() {
	m.hub.Shutdown()
}
