package clients

import (
	"context"
	"strings"

	mcpOAuthApplicationClients "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/clients"
	mcpOAuthDomainClients "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/clients"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

type RegisterClientCommand struct {
	ClientName              string
	RedirectURIs            []string
	GrantTypes              []string
	ResponseTypes           []string
	Scope                   string
	TokenEndpointAuthMethod string
	LogoURI                 string
	ClientURI               string
}

func (c *RegisterClientCommand) Validate() error {
	if strings.TrimSpace(c.ClientName) == "" {
		return mcpOAuthDomainClients.ErrInvalidClientName
	}

	if len(c.RedirectURIs) == 0 {
		return mcpOAuthDomainClients.ErrInvalidRedirectURI
	}
	for _, redirectURI := range c.RedirectURIs {
		if !mcpOAuthDomainOAuth2.IsValidRedirectURI(redirectURI) {
			return mcpOAuthDomainClients.ErrInvalidRedirectURI
		}
	}

	if !mcpOAuthDomainOAuth2.ParseScopes(c.Scope).IsValid() {
		return mcpOAuthDomainClients.ErrInvalidScopes
	}

	if strings.TrimSpace(c.TokenEndpointAuthMethod) != "" &&
		!mcpOAuthDomainOAuth2.TokenEndpointAuthMethod(c.TokenEndpointAuthMethod).IsValid() {
		return mcpOAuthDomainClients.ErrInvalidTokenEndpointAuthMethod
	}

	return nil
}

type RegisterClientResponse struct {
	ClientID                string   `json:"client_id"`
	ClientSecret            string   `json:"client_secret,omitempty"`
	ClientIDIssuedAt        int64    `json:"client_id_issued_at"`
	ClientSecretExpiresAt   int64    `json:"client_secret_expires_at"`
	ClientName              string   `json:"client_name"`
	RedirectURIs            []string `json:"redirect_uris"`
	GrantTypes              []string `json:"grant_types"`
	ResponseTypes           []string `json:"response_types"`
	Scope                   string   `json:"scope,omitempty"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
	LogoURI                 *string  `json:"logo_uri,omitempty"`
	ClientURI               *string  `json:"client_uri,omitempty"`
}

func (s *Service) RegisterClient(ctx context.Context, command *RegisterClientCommand) (*RegisterClientResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	authMethod := resolveTokenEndpointAuthMethod(command.TokenEndpointAuthMethod)

	clientID, err := mcpOAuthDomainOAuth2.GenerateSecret(mcpOAuthDomainClients.ClientIDPrefix)
	if err != nil {
		return nil, mcpOAuthApplicationClients.ErrFailedToGenerateClientCredentials.WithCause(err)
	}

	var clientSecret *mcpOAuthDomainOAuth2.Secret
	if authMethod.IsConfidential() {
		clientSecret, err = mcpOAuthDomainOAuth2.GenerateSecret(mcpOAuthDomainClients.ClientSecretPrefix)
		if err != nil {
			return nil, mcpOAuthApplicationClients.ErrFailedToGenerateClientCredentials.WithCause(err)
		}
	}

	var response *RegisterClientResponse

	err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		client, err := mcpOAuthDomainClients.New(
			clientID.Plain,
			secretHash(clientSecret),
			command.ClientName,
			command.RedirectURIs,
			resolveGrantTypes(command.GrantTypes),
			resolveResponseTypes(command.ResponseTypes),
			resolveScopes(command.Scope),
			authMethod,
			optionalString(command.LogoURI),
			optionalString(command.ClientURI),
		)
		if err != nil {
			return err
		}

		if err := s.clientRepository.Create(txCtx, client); err != nil {
			return mcpOAuthApplicationClients.ErrFailedToRegisterClient.WithCause(err)
		}

		response = &RegisterClientResponse{
			ClientID:                client.ClientID,
			ClientSecret:            secretPlain(clientSecret),
			ClientIDIssuedAt:        client.CreatedAt.Unix(),
			ClientSecretExpiresAt:   0,
			ClientName:              client.Name,
			RedirectURIs:            client.RedirectURIs,
			GrantTypes:              grantTypeStrings(client.GrantTypes),
			ResponseTypes:           responseTypeStrings(client.ResponseTypes),
			Scope:                   client.Scopes.String(),
			TokenEndpointAuthMethod: client.TokenEndpointAuthMethod.String(),
			LogoURI:                 client.LogoURI,
			ClientURI:               client.ClientURI,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func resolveTokenEndpointAuthMethod(value string) mcpOAuthDomainOAuth2.TokenEndpointAuthMethod {
	if strings.TrimSpace(value) == "" {
		return mcpOAuthDomainOAuth2.NoneAuthentication
	}
	return mcpOAuthDomainOAuth2.TokenEndpointAuthMethod(value)
}

func resolveGrantTypes(values []string) []mcpOAuthDomainOAuth2.GrantType {
	if len(values) == 0 {
		return []mcpOAuthDomainOAuth2.GrantType{
			mcpOAuthDomainOAuth2.AuthorizationCodeGrant,
			mcpOAuthDomainOAuth2.RefreshTokenGrant,
		}
	}

	return mcpOAuthDomainOAuth2.SupportedGrantTypes(values)
}

func resolveResponseTypes(values []string) []mcpOAuthDomainOAuth2.ResponseType {
	if len(values) == 0 {
		return []mcpOAuthDomainOAuth2.ResponseType{
			mcpOAuthDomainOAuth2.CodeResponseType,
		}
	}

	return mcpOAuthDomainOAuth2.SupportedResponseTypes(values)
}

func resolveScopes(value string) mcpOAuthDomainOAuth2.Scopes {
	scopes := mcpOAuthDomainOAuth2.ParseScopes(value)
	if len(scopes) == 0 {
		return mcpOAuthDomainOAuth2.DefaultScopes
	}
	return scopes
}

func grantTypeStrings(grantTypes []mcpOAuthDomainOAuth2.GrantType) []string {
	values := make([]string, 0, len(grantTypes))
	for _, grantType := range grantTypes {
		values = append(values, grantType.String())
	}
	return values
}

func responseTypeStrings(responseTypes []mcpOAuthDomainOAuth2.ResponseType) []string {
	values := make([]string, 0, len(responseTypes))
	for _, responseType := range responseTypes {
		values = append(values, responseType.String())
	}
	return values
}

func secretHash(secret *mcpOAuthDomainOAuth2.Secret) *string {
	if secret == nil {
		return nil
	}
	return new(secret.Hash)
}

func secretPlain(secret *mcpOAuthDomainOAuth2.Secret) string {
	if secret == nil {
		return ""
	}
	return secret.Plain
}

func optionalString(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return new(trimmed)
}
