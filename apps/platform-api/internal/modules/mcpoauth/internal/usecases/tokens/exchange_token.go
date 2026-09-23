package tokens

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"

	bnuuid "github.com/blocknextai/go-packages/uuid"
	mcpOAuthApplicationTokens "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/tokens"
	mcpOAuthDomainAuthorizationCodes "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/authorizationcodes"
	mcpOAuthDomainClients "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/clients"
	mcpOAuthDomainGrants "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/grants"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
	mcpOAuthDomainRefreshTokens "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/refreshtokens"
)

type ExchangeTokenCommand struct {
	GrantType    string
	ClientID     string
	ClientSecret string
	Code         string
	RedirectURI  string
	CodeVerifier string
	RefreshToken string
	Scope        string
	Resource     string
}

func (c *ExchangeTokenCommand) Validate() error {
	if strings.TrimSpace(c.ClientID) == "" {
		return mcpOAuthDomainClients.ErrInvalidClientID
	}

	switch mcpOAuthDomainOAuth2.GrantType(c.GrantType) {
	case mcpOAuthDomainOAuth2.AuthorizationCodeGrant:
		if strings.TrimSpace(c.Code) == "" {
			return mcpOAuthDomainAuthorizationCodes.ErrMissingCode
		}
		if strings.TrimSpace(c.CodeVerifier) == "" {
			return mcpOAuthDomainAuthorizationCodes.ErrMissingCodeVerifier
		}
	case mcpOAuthDomainOAuth2.RefreshTokenGrant:
		if strings.TrimSpace(c.RefreshToken) == "" {
			return mcpOAuthDomainRefreshTokens.ErrMissingRefreshToken
		}
	default:
		return mcpOAuthDomainOAuth2.ErrUnsupportedGrantType
	}

	return nil
}

type ExchangeTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

const (
	bearerTokenType = "Bearer"
)

func (s *Service) ExchangeToken(ctx context.Context, command *ExchangeTokenCommand) (*ExchangeTokenResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	client, err := s.clientResolver.Authenticate(ctx, command.ClientID, command.ClientSecret)
	if err != nil {
		return nil, err
	}

	grantType := mcpOAuthDomainOAuth2.GrantType(command.GrantType)
	if !client.SupportsGrantType(grantType) {
		return nil, mcpOAuthDomainClients.ErrUnauthorizedGrantType
	}

	switch grantType {
	case mcpOAuthDomainOAuth2.AuthorizationCodeGrant:
		return s.exchangeAuthorizationCode(ctx, client, command)
	case mcpOAuthDomainOAuth2.RefreshTokenGrant:
		return s.exchangeRefreshToken(ctx, client, command)
	default:
		return nil, mcpOAuthDomainOAuth2.ErrUnsupportedGrantType
	}
}

func (s *Service) exchangeAuthorizationCode(
	ctx context.Context,
	client *mcpOAuthDomainClients.Client,
	command *ExchangeTokenCommand,
) (*ExchangeTokenResponse, error) {
	code, err := s.authorizationCodeRepository.GetByCodeHash(ctx, mcpOAuthDomainOAuth2.HashSecret(command.Code))
	if err != nil {
		return nil, err
	}

	if code.ClientID != client.ClientID {
		return nil, mcpOAuthDomainAuthorizationCodes.ErrAuthorizationCodeClientMismatch
	}

	if code.IsUsed() {
		s.revokeGrant(ctx, code.GrantID, "authorization code replay")
		return nil, mcpOAuthDomainAuthorizationCodes.ErrAuthorizationCodeUsed
	}

	if code.IsExpired() {
		return nil, mcpOAuthDomainAuthorizationCodes.ErrAuthorizationCodeExpired
	}

	if strings.TrimSpace(command.RedirectURI) != code.RedirectURI {
		return nil, mcpOAuthDomainAuthorizationCodes.ErrAuthorizationCodeRedirectMismatch
	}

	if !code.VerifyCodeVerifier(command.CodeVerifier) {
		return nil, mcpOAuthDomainAuthorizationCodes.ErrInvalidCodeVerifier
	}

	if err := ensureResourceMatches(command.Resource, code.Resource); err != nil {
		return nil, err
	}

	var response *ExchangeTokenResponse

	err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		grant, err := s.activeGrant(txCtx, code.GrantID)
		if err != nil {
			return err
		}

		used, err := code.Use()
		if err != nil {
			return err
		}

		if err := s.authorizationCodeRepository.Update(txCtx, used); err != nil {
			if errors.Is(err, mcpOAuthDomainAuthorizationCodes.ErrAuthorizationCodeUsed) {
				return err
			}
			return mcpOAuthApplicationTokens.ErrFailedToExchangeToken.WithCause(err)
		}

		refreshToken, err := s.issueRefreshToken(txCtx, client, grant, code.Scopes, code.Resource)
		if err != nil {
			return err
		}

		response, err = s.newResponse(client, grant, code.Scopes, code.Resource, refreshToken)
		return err
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) exchangeRefreshToken(
	ctx context.Context,
	client *mcpOAuthDomainClients.Client,
	command *ExchangeTokenCommand,
) (*ExchangeTokenResponse, error) {
	token, err := s.refreshTokenRepository.GetByTokenHash(ctx, mcpOAuthDomainOAuth2.HashSecret(command.RefreshToken))
	if err != nil {
		return nil, err
	}

	if token.ClientID != client.ClientID {
		return nil, mcpOAuthDomainRefreshTokens.ErrRefreshTokenClientMismatch
	}

	if token.IsUsed() {
		s.revokeGrant(ctx, token.GrantID, "refresh token replay")
		return nil, mcpOAuthDomainRefreshTokens.ErrRefreshTokenUsed
	}

	if token.IsRevoked() {
		return nil, mcpOAuthDomainRefreshTokens.ErrRefreshTokenRevoked
	}

	if token.IsExpired() {
		return nil, mcpOAuthDomainRefreshTokens.ErrRefreshTokenExpired
	}

	if !s.isRefreshTokenResourceServed(command.Resource, token.Resource) {
		return nil, mcpOAuthDomainRefreshTokens.ErrRefreshTokenResourceMismatch
	}

	scopes := token.Scopes
	if requested := mcpOAuthDomainOAuth2.ParseScopes(command.Scope); len(requested) > 0 {
		if !token.Scopes.HasAll(requested) {
			return nil, mcpOAuthDomainGrants.ErrScopeExceedsGrant
		}
		scopes = requested
	}

	var response *ExchangeTokenResponse

	err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		grant, err := s.activeGrant(txCtx, token.GrantID)
		if err != nil {
			return err
		}

		used, err := token.Use()
		if err != nil {
			return err
		}

		if err := s.refreshTokenRepository.Update(txCtx, used); err != nil {
			if errors.Is(err, mcpOAuthDomainRefreshTokens.ErrRefreshTokenUsed) {
				return err
			}
			return mcpOAuthApplicationTokens.ErrFailedToExchangeToken.WithCause(err)
		}

		refreshToken, err := s.issueRefreshToken(txCtx, client, grant, scopes, token.Resource)
		if err != nil {
			return err
		}

		response, err = s.newResponse(client, grant, scopes, token.Resource, refreshToken)
		return err
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) activeGrant(ctx context.Context, grantID uuid.UUID) (*mcpOAuthDomainGrants.Grant, error) {
	grant, err := s.grantRepository.GetByID(ctx, grantID)
	if err != nil {
		return nil, err
	}

	if !grant.IsActive() {
		return nil, mcpOAuthDomainGrants.ErrGrantRevoked
	}

	return grant, nil
}

func (s *Service) issueRefreshToken(
	ctx context.Context,
	client *mcpOAuthDomainClients.Client,
	grant *mcpOAuthDomainGrants.Grant,
	scopes mcpOAuthDomainOAuth2.Scopes,
	resource string,
) (*mcpOAuthDomainOAuth2.Secret, error) {
	if !client.SupportsGrantType(mcpOAuthDomainOAuth2.RefreshTokenGrant) {
		return nil, nil
	}

	secret, err := mcpOAuthDomainOAuth2.GenerateSecret(mcpOAuthDomainRefreshTokens.TokenPrefix)
	if err != nil {
		return nil, mcpOAuthApplicationTokens.ErrFailedToGenerateToken.WithCause(err)
	}

	token, err := mcpOAuthDomainRefreshTokens.New(grant.ID, client.ClientID, secret.Hash, scopes, resource, s.refreshTokenTTL)
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepository.Create(ctx, token); err != nil {
		return nil, mcpOAuthApplicationTokens.ErrFailedToExchangeToken.WithCause(err)
	}

	return secret, nil
}

func (s *Service) newResponse(
	client *mcpOAuthDomainClients.Client,
	grant *mcpOAuthDomainGrants.Grant,
	scopes mcpOAuthDomainOAuth2.Scopes,
	resource string,
	refreshToken *mcpOAuthDomainOAuth2.Secret,
) (*ExchangeTokenResponse, error) {
	now := time.Now().UTC()

	accessToken, err := s.accessTokenSigner.Sign(&mcpOAuthDomainOAuth2.AccessToken{
		ID:             bnuuid.NewV7(),
		Subject:        grant.UserID,
		OrganizationID: grant.OrganizationID,
		ClientID:       client.ClientID,
		GrantID:        grant.ID,
		Scopes:         scopes,
		Resource:       resource,
		IssuedAt:       now,
		ExpiresAt:      now.Add(s.accessTokenTTL),
	})
	if err != nil {
		return nil, mcpOAuthApplicationTokens.ErrFailedToSignAccessToken.WithCause(err)
	}

	response := &ExchangeTokenResponse{
		AccessToken: accessToken,
		TokenType:   bearerTokenType,
		ExpiresIn:   int64(s.accessTokenTTL.Seconds()),
		Scope:       scopes.String(),
	}

	if refreshToken != nil {
		response.RefreshToken = refreshToken.Plain
	}

	return response, nil
}

func (s *Service) revokeGrant(ctx context.Context, grantID uuid.UUID, reason string) {
	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		grant, err := s.grantRepository.GetByID(txCtx, grantID)
		if err != nil {
			return err
		}

		return s.grantRevoker.Revoke(txCtx, grant)
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to revoke mcp oauth grant",
			"component", "mcpoauth",
			"grant_id", grantID.String(),
			"reason", reason,
			"error", err.Error(),
		)
		return
	}

	slog.WarnContext(ctx, "revoked mcp oauth grant",
		"component", "mcpoauth",
		"grant_id", grantID.String(),
		"reason", reason,
	)
}

func (s *Service) isRefreshTokenResourceServed(requested string, granted string) bool {
	if _, ok := mcpOAuthDomainOAuth2.ResolveResource(s.resourceBaseURL, granted); !ok {
		return false
	}

	return strings.TrimSpace(requested) == "" || mcpOAuthDomainOAuth2.CanonicalizeResource(requested) == granted
}

func ensureResourceMatches(requested string, granted string) error {
	if strings.TrimSpace(requested) == "" {
		return nil
	}

	if mcpOAuthDomainOAuth2.CanonicalizeResource(requested) != granted {
		return mcpOAuthDomainGrants.ErrResourceMismatch
	}

	return nil
}
