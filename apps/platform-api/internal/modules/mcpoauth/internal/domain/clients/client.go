package clients

import (
	"slices"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

const (
	ClientIDPrefix               = "mcpc_"
	ClientSecretPrefix           = "mcps_"
	metadataDocumentClientPrefix = "https://"
)

type Client struct {
	database.BaseEntity

	ClientID                string
	ClientSecretHash        *string
	Name                    string
	RedirectURIs            []string
	GrantTypes              []mcpOAuthDomainOAuth2.GrantType
	ResponseTypes           []mcpOAuthDomainOAuth2.ResponseType
	Scopes                  mcpOAuthDomainOAuth2.Scopes
	TokenEndpointAuthMethod mcpOAuthDomainOAuth2.TokenEndpointAuthMethod
	LogoURI                 *string
	ClientURI               *string
}

func New(
	clientID string,
	clientSecretHash *string,
	name string,
	redirectURIs []string,
	grantTypes []mcpOAuthDomainOAuth2.GrantType,
	responseTypes []mcpOAuthDomainOAuth2.ResponseType,
	scopes mcpOAuthDomainOAuth2.Scopes,
	tokenEndpointAuthMethod mcpOAuthDomainOAuth2.TokenEndpointAuthMethod,
	logoURI *string,
	clientURI *string,
) (*Client, error) {
	now := time.Now().UTC()

	client := &Client{
		ID:                      bnuuid.NewV7(),
		CreatedAt:               now,
		UpdatedAt:               now,
		ClientID:                clientID,
		ClientSecretHash:        clientSecretHash,
		Name:                    name,
		RedirectURIs:            redirectURIs,
		GrantTypes:              grantTypes,
		ResponseTypes:           responseTypes,
		Scopes:                  scopes,
		TokenEndpointAuthMethod: tokenEndpointAuthMethod,
		LogoURI:                 logoURI,
		ClientURI:               clientURI,
	}

	return client.validateThenReturn()
}

func IsMetadataDocumentClientID(clientID string) bool {
	return strings.HasPrefix(clientID, metadataDocumentClientPrefix)
}

func (c *Client) HasRedirectURI(redirectURI string) bool {
	return slices.Contains(c.RedirectURIs, redirectURI)
}

func (c *Client) SupportsGrantType(grantType mcpOAuthDomainOAuth2.GrantType) bool {
	return slices.Contains(c.GrantTypes, grantType)
}

func (c *Client) VerifySecret(plainSecret string) bool {
	if c.ClientSecretHash == nil {
		return strings.TrimSpace(plainSecret) == ""
	}
	return mcpOAuthDomainOAuth2.HashSecret(plainSecret) == *c.ClientSecretHash
}

func (c *Client) AllowedScopes(requested mcpOAuthDomainOAuth2.Scopes) mcpOAuthDomainOAuth2.Scopes {
	if len(c.Scopes) == 0 {
		return requested
	}

	allowed := make(mcpOAuthDomainOAuth2.Scopes, 0, len(requested))
	for _, scope := range requested {
		if c.Scopes.Has(scope) {
			allowed = append(allowed, scope)
		}
	}
	return allowed
}

func (c *Client) validateThenReturn() (*Client, error) {
	if strings.TrimSpace(c.Name) == "" {
		return nil, ErrInvalidClientName
	}

	if strings.TrimSpace(c.ClientID) == "" {
		return nil, ErrInvalidClientID
	}

	if len(c.RedirectURIs) == 0 {
		return nil, ErrInvalidRedirectURI
	}
	for _, redirectURI := range c.RedirectURIs {
		if !mcpOAuthDomainOAuth2.IsValidRedirectURI(redirectURI) {
			return nil, ErrInvalidRedirectURI
		}
	}

	if len(c.GrantTypes) == 0 {
		return nil, ErrInvalidGrantTypes
	}
	for _, grantType := range c.GrantTypes {
		if !grantType.IsValid() {
			return nil, ErrInvalidGrantTypes
		}
	}

	if len(c.ResponseTypes) == 0 {
		return nil, ErrInvalidResponseTypes
	}
	for _, responseType := range c.ResponseTypes {
		if !responseType.IsValid() {
			return nil, ErrInvalidResponseTypes
		}
	}

	if !c.Scopes.IsValid() {
		return nil, ErrInvalidScopes
	}

	if !c.TokenEndpointAuthMethod.IsValid() {
		return nil, ErrInvalidTokenEndpointAuthMethod
	}

	return c, nil
}
