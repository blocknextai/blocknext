package getallapikeys

import (
	apiKeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/apikeys/domain/apikeys"
)

func MapAPIKeysToResponse(apiKeys []*apiKeysDomainAPIKeys.APIKey) []*APIKeyResponse {
	responses := make([]*APIKeyResponse, 0, len(apiKeys))
	for _, k := range apiKeys {
		responses = append(responses, &APIKeyResponse{
			ID:         k.ID,
			Name:       k.Name,
			Scopes:     k.Scopes,
			LastUsedAt: k.LastUsedAt,
			CreatedAt:  k.CreatedAt,
			UpdatedAt:  k.UpdatedAt,
		})
	}
	return responses
}
