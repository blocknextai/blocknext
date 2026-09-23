package usecases

import (
	uploadDomain "github.com/blocknextai/file-gateway-api/internal/modules/upload/internal/domain"
)

func (s *Service) GetRule(uploadID string) (*uploadDomain.UploadRule, error) {
	return uploadDomain.GetUploadRule(uploadID)
}
