package credentials

import (
	"context"

	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/credentials"
)

type GetCredentialByIDQuery struct {
	ID string
}

type GetCredentialByIDResponse = nodeEngineDomainCredentials.CredentialManager

func (s *Service) GetCredentialByID(ctx context.Context, request *GetCredentialByIDQuery) (*GetCredentialByIDResponse, error) {
	credential, exists := s.credentialService.GetCredentialByID(request.ID)
	if !exists {
		return nil, nodeEngineDomainCredentials.ErrCredentialNotFound
	}

	return &credential, nil
}
