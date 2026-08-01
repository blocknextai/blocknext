package platformconfig

import (
	"log/slog"
	"strings"
	"sync"

	"github.com/blocknextai/go-packages/base64"
	"github.com/blocknextai/go-packages/json"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	platformDomain "github.com/blocknextai/platform-api/internal/platform/domain/platformcredentials"
)

type PlatformConfigLoader struct {
	platforms map[string]*platformDomain.PlatformCredential
	mu        sync.RWMutex
}

func NewPlatformConfigLoader(
	credentialConfigs map[string]string,
	credentials []nodeEngineDomainCredentials.CredentialManager,
) *PlatformConfigLoader {
	loader := &PlatformConfigLoader{
		platforms: make(map[string]*platformDomain.PlatformCredential),
	}

	for _, credential := range credentials {
		platformID := credential.GetPlatformID()
		if platformID == "" {
			continue
		}
		data := strings.TrimSpace(credentialConfigs[platformID])
		if data == "" {
			continue
		}
		if !loader.addCredential(credential.GetID(), data) {
			continue
		}
		credential.SetIsSupportPlatform(true)
	}

	return loader
}

func (loader *PlatformConfigLoader) addCredential(credentialID string, credentialData string) bool {
	decoded, err := base64.Decode(credentialData)
	if err != nil {
		slog.Error("Failed to decode base64 for credential",
			"component", "platform_config",
			"credential_id", credentialID,
			"error", err)
		return false
	}

	var data map[string]any
	if err := json.Unmarshal(decoded, &data); err != nil {
		slog.Error("Failed to parse credential JSON",
			"component", "platform_config",
			"credential_id", credentialID,
			"error", err)
		return false
	}

	if data == nil {
		slog.Error("Credential data is nil",
			"component", "platform_config",
			"credential_id", credentialID)
		return false
	}

	loader.platforms[credentialID] = platformDomain.NewPlatformCredential(
		credentialID,
		data,
	)
	return true
}

func (loader *PlatformConfigLoader) GetPlatformCredential(id string) (*platformDomain.PlatformCredential, bool) {
	loader.mu.RLock()
	defer loader.mu.RUnlock()

	platform, ok := loader.platforms[id]
	return platform, ok
}

func (loader *PlatformConfigLoader) GetAllPlatformCredentials() []*platformDomain.PlatformCredential {
	loader.mu.RLock()
	defer loader.mu.RUnlock()

	result := make([]*platformDomain.PlatformCredential, 0, len(loader.platforms))
	for _, platform := range loader.platforms {
		result = append(result, platform)
	}
	return result
}
