package adapter

import (
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonDomainCredential "github.com/blocknextai/platform-api/internal/common/domain/credential"
	commonAuth "github.com/blocknextai/platform-api/internal/common/presentation/auth"
	"github.com/google/uuid"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func extractOwner(req *mcpsdk.CallToolRequest) (commonDomain.OwnerType, uuid.UUID, error) {
	if req.Extra == nil {
		return "", uuid.Nil, ErrMissingOwner
	}

	ownerType := commonDomain.OwnerType(req.Extra.Header.Get(commonAuth.OwnerTypeHeader))
	if ownerType == "" {
		return "", uuid.Nil, ErrMissingOwner
	}

	ownerID, err := uuid.Parse(req.Extra.Header.Get(commonAuth.OwnerIDHeader))
	if err != nil {
		return "", uuid.Nil, ErrMissingOwner.WithCause(err)
	}

	return ownerType, ownerID, nil
}

func extractAPIKeyID(req *mcpsdk.CallToolRequest) *uuid.UUID {
	if req.Extra == nil {
		return nil
	}

	apiKeyID, err := uuid.Parse(req.Extra.Header.Get(commonAuth.APIKeyIDHeader))
	if err != nil {
		return nil
	}

	return &apiKeyID
}

func splitInput(input map[string]any, credentialKeys []string) (map[string]any, map[string]any) {
	if len(credentialKeys) == 0 {
		return map[string]any{}, input
	}

	references := map[string]any{}
	if raw, ok := input[credentialsKey].(map[string]any); ok {
		references = raw
	}

	data := make(map[string]any, len(input))
	for k, v := range input {
		if k == credentialsKey {
			continue
		}
		data[k] = v
	}

	return references, data
}

func credentialReferences(references map[string]any) map[string]any {
	if len(references) == 0 {
		return nil
	}

	kept := make(map[string]any, len(references))
	for key, raw := range references {
		reference, ok := raw.(string)
		if !ok {
			continue
		}
		if _, _, err := commonDomainCredential.ParseReference(reference); err != nil {
			continue
		}
		kept[key] = reference
	}

	return kept
}
