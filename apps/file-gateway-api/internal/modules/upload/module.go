package upload

import (
	"github.com/gofiber/fiber/v3"

	storage "github.com/blocknextai/file-gateway-api/internal/modules/storage"
	uploadHTTP "github.com/blocknextai/file-gateway-api/internal/modules/upload/internal/http"
	uploadUseCases "github.com/blocknextai/file-gateway-api/internal/modules/upload/internal/usecases"
)

type Dependencies struct {
	StorageProvider storage.Provider
}

type Module struct {
	UploadService *uploadUseCases.Service
}

func NewModule(deps Dependencies) *Module {
	service := uploadUseCases.NewService(deps.StorageProvider)
	return &Module{
		UploadService: service,
	}
}

func (m *Module) Register(router fiber.Router) {
	uploadHTTP.RegisterRoutes(router, m.UploadService)
}
