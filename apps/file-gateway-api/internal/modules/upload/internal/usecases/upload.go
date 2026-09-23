package usecases

import (
	"context"

	storage "github.com/blocknextai/file-gateway-api/internal/modules/storage"
	uploadDomain "github.com/blocknextai/file-gateway-api/internal/modules/upload/internal/domain"
)

func (s *Service) Upload(
	ctx context.Context,
	rule *uploadDomain.UploadRule,
	file *uploadDomain.File,
) (*storage.UploadResult, error) {
	if err := uploadDomain.RunValidators(file, uploadDomain.CreateValidators(*rule)); err != nil {
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
