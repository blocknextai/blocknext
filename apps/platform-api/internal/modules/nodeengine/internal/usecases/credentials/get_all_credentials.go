package credentials

import (
	"context"

	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/credentials"
)

type GetAllCredentialsQuery struct{}

type GetAllCredentialsResponse = []credentials.CredentialManager

func (s *Service) GetAllCredentials(ctx context.Context, _ *GetAllCredentialsQuery) (*GetAllCredentialsResponse, error) {
	return new(s.credentialService.GetAllCredentials()), nil
}
