package getcredentialsfornodes

import (
	"context"

	credentialsApplicationCredentials "github.com/blocknextai/platform-api/internal/credentials/application/credentials"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/credentials/domain/credentials"
	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/nodeengine/application/credentials"
)

type Handler struct {
	credentialRepository        credentialsDomainCredentials.CredentialRepository
	nodeEngineCredentialService nodeEngineApplicationCredentials.CredentialService
}

func New(
	credentialRepository credentialsDomainCredentials.CredentialRepository,
	nodeEngineCredentialService nodeEngineApplicationCredentials.CredentialService,
) *Handler {
	return &Handler{
		credentialRepository:        credentialRepository,
		nodeEngineCredentialService: nodeEngineCredentialService,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetCredentialsForNodesQuery) (*GetCredentialsForNodesResponse, error) {
	if len(request.NodeIDs) == 0 {
		return nil, nil
	}

	credentialSchemas := h.nodeEngineCredentialService.GetCredentialSchemasByNodeIDs(request.NodeIDs)

	keys := make([]string, 0, len(credentialSchemas))
	for _, schema := range credentialSchemas {
		keys = append(keys, schema.GetID())
	}

	if len(keys) == 0 {
		return nil, nil
	}

	credentials, err := h.credentialRepository.GetAllByOwnerAndKeys(ctx, request.OwnerType, request.OwnerID, keys)
	if err != nil {
		return nil, err
	}

	return new(GetCredentialsForNodesResponse(credentialsApplicationCredentials.MapCredentialsToResponse(credentials))), nil
}
