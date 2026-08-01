package getplatformcredentialbyid

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
)

type GetPlatformCredentialByIDResponse struct {
	ID             string                     `json:"id"`
	Name           string                     `json:"name"`
	Description    string                     `json:"description"`
	Icon           credentials.CredentialIcon `json:"icon"`
	IsOAuth2       bool                       `json:"isOAuth2"`
	SupportedNodes *[]string                  `json:"supportedNodes"`
}
