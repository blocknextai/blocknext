package credentials

import (
	"maps"

	commonDomainCredential "github.com/blocknextai/platform-api/internal/common/domain/credential"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type CredentialProcessor interface {
	ProcessCredentialData(schema *gjs.Schema, key string, data map[string]any) map[string]any
	ShouldSkipField(schema *gjs.Schema, key string, value any) bool
	ProcessUpdateData(schema *gjs.Schema, data map[string]any) map[string]any
}

type credentialProcessor struct {
	credentialService CredentialService
}

func NewCredentialProcessor(
	credentialService CredentialService,
) CredentialProcessor {
	return &credentialProcessor{
		credentialService: credentialService,
	}
}

func (cp *credentialProcessor) ProcessCredentialData(
	schema *gjs.Schema,
	key string,
	data map[string]any,
) map[string]any {
	processedData := make(map[string]any)
	maps.Copy(processedData, data)

	if schema != nil && schema.Properties != nil {
		for propKey, propSchema := range schema.Properties {
			if propSchema == nil || !propSchema.WriteOnly {
				continue
			}
			if _, exists := processedData[propKey]; exists {
				processedData[propKey] = commonDomainCredential.CredentialMaskValue
			}
		}
	}

	// TODO: update static strings with a more generic solution
	if cp.credentialService != nil && (cp.credentialService.IsOAuth1Credential(key) || cp.credentialService.IsOAuth2Credential(key)) {
		_, exists := processedData["oauthToken"]
		processedData["oauthToken"] = exists
	}

	return processedData
}

func (cp *credentialProcessor) ShouldSkipField(
	schema *gjs.Schema,
	key string,
	value any,
) bool {
	if schema == nil || schema.Properties == nil {
		return false
	}

	propSchema, ok := schema.Properties[key]
	if !ok || propSchema == nil || !propSchema.WriteOnly {
		return false
	}

	if strValue, ok := value.(string); ok && strValue == commonDomainCredential.CredentialMaskValue {
		return true
	}

	return false
}

func (cp *credentialProcessor) ProcessUpdateData(
	schema *gjs.Schema,
	data map[string]any,
) map[string]any {
	if schema == nil || schema.Properties == nil {
		return data
	}

	processedData := make(map[string]any)
	for key, value := range data {
		if !cp.ShouldSkipField(schema, key, value) {
			processedData[key] = value
		}
	}

	return processedData
}
