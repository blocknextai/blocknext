package upload

import (
	storageDomain "github.com/blocknextai/file-gateway-api/internal/storage/domain/storage"
	uploadApplicationUpload "github.com/blocknextai/file-gateway-api/internal/upload/application/upload"
	uploadPresentation "github.com/blocknextai/file-gateway-api/internal/upload/presentation"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	StorageProvider storageDomain.Provider
}

type Module struct {
	UploadService uploadApplicationUpload.Service
}

func NewModule(deps Dependencies) *Module {
	service := uploadApplicationUpload.NewService(deps.StorageProvider)
	return &Module{
		UploadService: service,
	}
}

func (m *Module) Register(router fiber.Router) {
	uploadPresentation.RegisterRoutes(router, m.UploadService)
}
