package usecases

import (
	"errors"
	"log/slog"

	storage "github.com/blocknextai/file-gateway-api/internal/modules/storage"
	uploadDomain "github.com/blocknextai/file-gateway-api/internal/modules/upload/internal/domain"
	bnfile "github.com/blocknextai/go-packages/file"
	"github.com/blocknextai/go-packages/uuid"
)

type Service struct {
	storageProvider storage.Provider
}

func NewService(storageProvider storage.Provider) *Service {
	return &Service{
		storageProvider: storageProvider,
	}
}

func resolveFilename(file *uploadDomain.File, rule *uploadDomain.UploadRule) string {
	if rule.IsOverrideFilename {
		return uuid.NewV7().String() + bnfile.ExtensionFromMIMEType(file.ContentType)
	}
	return file.Filename
}

func handleStorageError(err error) error {
	if errors.Is(err, uploadDomain.ErrMaxSizeExceeded) {
		return uploadDomain.ErrMaxSizeExceeded
	}
	slog.Error("storage upload failed", "component", "UploadService", "error", err)
	return uploadDomain.ErrStorageError
}
