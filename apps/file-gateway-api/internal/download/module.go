package download

import (
	"github.com/blocknextai/file-gateway-api/internal/config"
	downloadApplicationDownload "github.com/blocknextai/file-gateway-api/internal/download/application/download"
	downloadDomainDownload "github.com/blocknextai/file-gateway-api/internal/download/domain/download"
	downloadPresentation "github.com/blocknextai/file-gateway-api/internal/download/presentation"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	Options config.DownloadOptions
}

type Module struct {
	DownloadService downloadDomainDownload.Service
}

func NewModule(deps Dependencies) *Module {
	service := downloadApplicationDownload.NewService(deps.Options.MaxSize, deps.Options.Timeout)
	return &Module{
		DownloadService: service,
	}
}

func (m *Module) Register(router fiber.Router) {
	downloadPresentation.RegisterRoutes(router, m.DownloadService)
}
