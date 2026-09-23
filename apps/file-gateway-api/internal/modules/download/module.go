package download

import (
	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/file-gateway-api/internal/config"
	downloadHTTP "github.com/blocknextai/file-gateway-api/internal/modules/download/internal/http"
	downloadUseCases "github.com/blocknextai/file-gateway-api/internal/modules/download/internal/usecases"
)

type Dependencies struct {
	Options config.DownloadOptions
}

type Module struct {
	DownloadService *downloadUseCases.Service
}

func NewModule(deps Dependencies) *Module {
	service := downloadUseCases.NewService(deps.Options.MaxSize, deps.Options.Timeout)
	return &Module{
		DownloadService: service,
	}
}

func (m *Module) Register(router fiber.Router) {
	downloadHTTP.RegisterRoutes(router, m.DownloadService)
}
