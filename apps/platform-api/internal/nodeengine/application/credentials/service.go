package credentials

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
)

type CredentialService interface {
	GetAllCredentials() []credentials.CredentialManager
	GetCredentialByID(id string) (credentials.CredentialManager, bool)
	GetCredentialSchemasByNodeIDs(nodeIDs []string) []credentials.CredentialManager
	GetHiddenObjectsByCredentialID(id string) map[string]any
	IsOAuth1Credential(id string) bool
	IsOAuth2Credential(id string) bool
	GetRefreshableCredential(id string) (credentials.RefreshableCredential, bool)
}

type credentialService struct{}

func NewCredentialService() CredentialService {
	return &credentialService{}
}

func (s *credentialService) GetAllCredentials() []credentials.CredentialManager {
	return credentials.GetAllCredentials()
}

func (s *credentialService) GetCredentialByID(id string) (credentials.CredentialManager, bool) {
	return credentials.GetCredential(id)
}

func (s *credentialService) GetCredentialSchemasByNodeIDs(nodeIDs []string) []credentials.CredentialManager {
	credentialSchemasByKey := make(map[string]credentials.CredentialManager)
	credentialSchemas := make([]credentials.CredentialManager, 0)

	for _, nodeID := range nodeIDs {
		node, exists := nodes.GetNode(nodeID)
		if !exists {
			continue
		}

		for _, supportedCredential := range node.GetSupportedCredentials() {
			if _, exists := credentialSchemasByKey[supportedCredential]; exists {
				continue
			}

			credential, exists := credentials.GetCredential(supportedCredential)
			if !exists {
				continue
			}

			credentialSchemasByKey[supportedCredential] = credential
			credentialSchemas = append(credentialSchemas, credential)
		}
	}

	return credentialSchemas
}

func (s *credentialService) GetHiddenObjectsByCredentialID(id string) map[string]any {
	data := make(map[string]any)

	credential, exists := credentials.GetCredential(id)
	if !exists {
		return nil
	}

	schema := credential.GetSchema()
	if schema == nil || schema.Properties == nil {
		return nil
	}

	for key, propSchema := range schema.Properties {
		if propSchema == nil {
			continue
		}
		hidden, _ := propSchema.Extra["hidden"].(bool)
		if !hidden {
			continue
		}

		var defaultValue any
		if len(propSchema.Default) > 0 {
			_ = json.Unmarshal(propSchema.Default, &defaultValue)
		}
		data[key] = defaultValue
	}

	return data
}

func (s *credentialService) IsOAuth1Credential(id string) bool {
	credential, exists := credentials.GetCredential(id)
	if !exists {
		return false
	}

	return credential.GetIsOAuth1()
}

func (s *credentialService) IsOAuth2Credential(id string) bool {
	credential, exists := credentials.GetCredential(id)
	if !exists {
		return false
	}

	return credential.GetIsOAuth2()
}

func (s *credentialService) GetRefreshableCredential(id string) (credentials.RefreshableCredential, bool) {
	credential, exists := credentials.GetCredential(id)
	if !exists {
		return nil, false
	}

	refreshable, ok := credential.(credentials.RefreshableCredential)
	return refreshable, ok
}
