package credentials

import (
	"github.com/blocknextai/go-packages/cast"
)

type CredentialHelper struct {
	credentials map[string]any
}

func NewCredentialHelper(credentials map[string]any) *CredentialHelper {
	return &CredentialHelper{
		credentials: credentials,
	}
}

func GetCredentials(credentials map[string]any, credentialType string) *CredentialData {
	return NewCredentialHelper(credentials).GetCredentials(credentialType)
}

func (c *CredentialHelper) GetCredentials(credentialType string) *CredentialData {
	if creds, ok := c.credentials[credentialType].(map[string]any); ok {
		return &CredentialData{data: creds}
	}
	return &CredentialData{data: make(map[string]any)}
}

type CredentialData struct {
	data map[string]any
}

func (c *CredentialData) String(key string) string {
	if val, exists := c.data[key]; exists {
		return cast.ToString(val)
	}
	return ""
}

func (c *CredentialData) Object(key string) *CredentialData {
	if val, exists := c.data[key]; exists {
		if obj, ok := val.(map[string]any); ok {
			return &CredentialData{data: obj}
		}
	}
	return &CredentialData{data: make(map[string]any)}
}
