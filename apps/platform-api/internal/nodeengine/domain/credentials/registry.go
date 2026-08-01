package credentials

import (
	"strings"
)

var (
	instance *CredentialRegistry
)

type CredentialRegistry struct {
	credentialsMap map[string]CredentialManager
	credentials    []CredentialManager
}

func NewCredentialRegistry() *CredentialRegistry {
	if instance == nil {
		instance = &CredentialRegistry{
			credentialsMap: make(map[string]CredentialManager),
			credentials:    make([]CredentialManager, 0),
		}
	}
	return instance
}

func (r *CredentialRegistry) RegisterCredential(credential CredentialManager) {
	r.credentialsMap[credential.GetID()] = credential
	r.credentials = append(r.credentials, credential)
}

func (r *CredentialRegistry) GetCredential(id string) (CredentialManager, bool) {
	credential, ok := r.credentialsMap[id]
	return credential, ok
}

func (r *CredentialRegistry) GetAllCredentials() []CredentialManager {
	return r.credentials
}

func RegisterCredential(credential CredentialManager) {
	if credential.GetDisabled() {
		return
	}

	id := strings.ToLower(credential.GetID())
	if strings.Contains(id, "oauth2") {
		credential.SetIsOAuth2(true)
	}
	if strings.Contains(id, "oauth1") {
		credential.SetIsOAuth1(true)
	}

	NewCredentialRegistry().RegisterCredential(credential)
}

func GetCredential(id string) (CredentialManager, bool) {
	return NewCredentialRegistry().GetCredential(id)
}

func GetAllCredentials() []CredentialManager {
	return NewCredentialRegistry().GetAllCredentials()
}
