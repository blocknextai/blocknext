package upload

import (
	"context"
	"errors"
	"log/slog"

	storageDomain "github.com/blocknextai/file-gateway-api/internal/storage/domain/storage"
	uploadDomain "github.com/blocknextai/file-gateway-api/internal/upload/domain"
	uploadDomainUpload "github.com/blocknextai/file-gateway-api/internal/upload/domain/upload"
	bnfile "github.com/blocknextai/go-packages/file"
	"github.com/blocknextai/go-packages/uuid"
)

type Service interface {
	GetRule(uploadID string) (*uploadDomainUpload.UploadRule, error)
	Upload(
		ctx context.Context,
		rule *uploadDomainUpload.UploadRule,
		file *uploadDomainUpload.File,
	) (*storageDomain.UploadResult, error)
}

type service struct {
	storageProvider storageDomain.Provider
}

func NewService(storageProvider storageDomain.Provider) Service {
	return &service{
		storageProvider: storageProvider,
	}
}

func (s *service) GetRule(uploadID string) (*uploadDomainUpload.UploadRule, error) {
	return uploadDomainUpload.GetUploadRule(uploadID)
}

func (s *service) Upload(
	ctx context.Context,
	rule *uploadDomainUpload.UploadRule,
	file *uploadDomainUpload.File,
) (*storageDomain.UploadResult, error) {
	if err := uploadDomainUpload.RunValidators(file, uploadDomainUpload.CreateValidators(*rule)); err != nil {
		return nil, err
	}

	filename := resolveFilename(file, rule)
	key := rule.DefaultFolder + filename

	uploadResult, err := s.storageProvider.Upload(ctx, rule.IsPublic, key, file.ContentReader, file.ContentType, file.Size)
	if err != nil {
		return nil, handleStorageError(err)
	}

	return uploadResult, nil
}

func resolveFilename(file *uploadDomainUpload.File, rule *uploadDomainUpload.UploadRule) string {
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
