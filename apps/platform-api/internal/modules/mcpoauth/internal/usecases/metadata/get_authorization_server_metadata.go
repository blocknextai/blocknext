package metadata

import (
	"context"

	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

type GetAuthorizationServerMetadataQuery struct{}

type GetAuthorizationServerMetadataResponse struct {
	Issuer                                     string   `json:"issuer"`
	AuthorizationEndpoint                      string   `json:"authorization_endpoint"`
	TokenEndpoint                              string   `json:"token_endpoint"`
	RegistrationEndpoint                       string   `json:"registration_endpoint"`
	RevocationEndpoint                         string   `json:"revocation_endpoint"`
	ScopesSupported                            []string `json:"scopes_supported"`
	ResponseTypesSupported                     []string `json:"response_types_supported"`
	ResponseModesSupported                     []string `json:"response_modes_supported"`
	GrantTypesSupported                        []string `json:"grant_types_supported"`
	TokenEndpointAuthMethodsSupported          []string `json:"token_endpoint_auth_methods_supported"`
	RevocationEndpointAuthMethodsSupported     []string `json:"revocation_endpoint_auth_methods_supported"`
	CodeChallengeMethodsSupported              []string `json:"code_challenge_methods_supported"`
	AuthorizationResponseIssParameterSupported bool     `json:"authorization_response_iss_parameter_supported"`
	ClientIDMetadataDocumentSupported          bool     `json:"client_id_metadata_document_supported"`
}

func (s *Service) GetAuthorizationServerMetadata(ctx context.Context, _ *GetAuthorizationServerMetadataQuery) (*GetAuthorizationServerMetadataResponse, error) {
	return MapAuthorizationServerMetadataToResponse(s.issuerURL), nil
}

const (
	queryResponseMode = "query"
)

func MapAuthorizationServerMetadataToResponse(issuerURL string) *GetAuthorizationServerMetadataResponse {
	authMethods := []string{
		mcpOAuthDomainOAuth2.NoneAuthentication.String(),
		mcpOAuthDomainOAuth2.ClientSecretBasicAuthentication.String(),
		mcpOAuthDomainOAuth2.ClientSecretPostAuthentication.String(),
	}

	return &GetAuthorizationServerMetadataResponse{
		Issuer:                issuerURL,
		AuthorizationEndpoint: issuerURL + mcpOAuthDomainOAuth2.AuthorizationEndpointPath,
		TokenEndpoint:         issuerURL + mcpOAuthDomainOAuth2.TokenEndpointPath,
		RegistrationEndpoint:  issuerURL + mcpOAuthDomainOAuth2.RegistrationEndpointPath,
		RevocationEndpoint:    issuerURL + mcpOAuthDomainOAuth2.RevocationEndpointPath,
		ScopesSupported: mcpOAuthDomainOAuth2.Scopes{
			mcpOAuthDomainOAuth2.MCPInvokeScope,
			mcpOAuthDomainOAuth2.PlatformReadScope,
			mcpOAuthDomainOAuth2.PlatformWriteScope,
			mcpOAuthDomainOAuth2.OfflineAccessScope,
		}.Strings(),
		ResponseTypesSupported: []string{mcpOAuthDomainOAuth2.CodeResponseType.String()},
		ResponseModesSupported: []string{queryResponseMode},
		GrantTypesSupported: []string{
			mcpOAuthDomainOAuth2.AuthorizationCodeGrant.String(),
			mcpOAuthDomainOAuth2.RefreshTokenGrant.String(),
		},
		TokenEndpointAuthMethodsSupported:          authMethods,
		RevocationEndpointAuthMethodsSupported:     authMethods,
		CodeChallengeMethodsSupported:              []string{mcpOAuthDomainOAuth2.S256CodeChallengeMethod.String()},
		AuthorizationResponseIssParameterSupported: true,
		ClientIDMetadataDocumentSupported:          true,
	}
}
