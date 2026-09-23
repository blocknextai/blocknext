package credentials

import (
	"context"

	"github.com/google/uuid"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	credentialsApplicationCredentials "github.com/blocknextai/platform-api/internal/modules/credentials/internal/application/credentials"
)

type GetCredentialsForNodesQuery struct {
	OwnerType commonDomain.OwnerType
	OwnerID   uuid.UUID
	NodeIDs   []string
}

type GetCredentialsForNodesResponse = []credentialsApplicationCredentials.CredentialResponse

func (s *Service) GetCredentialsForNodes(ctx context.Context, request *GetCredentialsForNodesQuery) (*GetCredentialsForNodesResponse, error) {
	if len(request.NodeIDs) == 0 {
		return nil, nil
	}

	credentialSchemas := s.nodeEngineCredentialService.GetCredentialSchemasByNodeIDs(request.NodeIDs)

	keys := make([]string, 0, len(credentialSchemas))
	for _, schema := range credentialSchemas {
		keys = append(keys, schema.GetID())
	}

	if len(keys) == 0 {
		return nil, nil
	}

	credentials, err := s.credentialRepository.GetAllByOwnerAndKeys(ctx, request.OwnerType, request.OwnerID, keys)
	if err != nil {
		return nil, err
	}

	return new(GetCredentialsForNodesResponse(credentialsApplicationCredentials.MapCredentialsToResponse(credentials))), nil
}
