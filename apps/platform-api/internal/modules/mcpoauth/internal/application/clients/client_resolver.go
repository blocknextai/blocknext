package clients

import (
	"context"

	mcpOAuthDomainClients "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/clients"
)

type ClientResolver struct {
	clientRepository        mcpOAuthDomainClients.ClientRepository
	metadataDocumentFetcher mcpOAuthDomainClients.MetadataDocumentFetcher
}

func NewClientResolver(
	clientRepository mcpOAuthDomainClients.ClientRepository,
	metadataDocumentFetcher mcpOAuthDomainClients.MetadataDocumentFetcher,
) *ClientResolver {
	return &ClientResolver{
		clientRepository:        clientRepository,
		metadataDocumentFetcher: metadataDocumentFetcher,
	}
}

func (r *ClientResolver) Resolve(ctx context.Context, clientID string) (*mcpOAuthDomainClients.Client, error) {
	if mcpOAuthDomainClients.IsMetadataDocumentClientID(clientID) {
		return r.metadataDocumentFetcher.Fetch(ctx, clientID)
	}

	return r.clientRepository.GetByClientID(ctx, clientID)
}

func (r *ClientResolver) Authenticate(ctx context.Context, clientID string, clientSecret string) (*mcpOAuthDomainClients.Client, error) {
	client, err := r.Resolve(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if client.TokenEndpointAuthMethod.IsConfidential() && !client.VerifySecret(clientSecret) {
		return nil, mcpOAuthDomainClients.ErrInvalidClientCredentials
	}

	return client, nil
}
