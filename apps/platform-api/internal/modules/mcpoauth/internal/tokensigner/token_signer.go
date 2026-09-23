package tokensigner

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	mcpOAuthApplicationTokens "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/tokens"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

const (
	accessTokenType         = "at+jwt"
	minimumSigningKeyLength = 32
	clockSkewLeeway         = 30 * time.Second
)

var (
	signingMethod = jwt.SigningMethodHS256
)

type accessTokenClaims struct {
	jwt.RegisteredClaims

	OrganizationID string `json:"organization_id"`
	ClientID       string `json:"client_id"`
	GrantID        string `json:"grant_id"`
	Scope          string `json:"scope"`
}

type TokenSigner struct {
	issuer    string
	secretKey []byte
}

func New(issuer string, signingKey string) (*TokenSigner, error) {
	if strings.TrimSpace(issuer) == "" {
		return nil, mcpOAuthApplicationTokens.ErrInvalidIssuer
	}

	trimmedKey := strings.TrimSpace(signingKey)
	if len(trimmedKey) < minimumSigningKeyLength {
		return nil, mcpOAuthApplicationTokens.ErrInvalidSigningKey
	}

	return &TokenSigner{
		issuer:    issuer,
		secretKey: []byte(trimmedKey),
	}, nil
}

func (s *TokenSigner) Sign(accessToken *mcpOAuthDomainOAuth2.AccessToken) (string, error) {
	token := jwt.NewWithClaims(signingMethod, accessTokenClaims{
		Issuer:         s.issuer,
		Subject:        accessToken.Subject.String(),
		Audience:       []string{accessToken.Resource},
		ExpiresAt:      jwt.NewNumericDate(accessToken.ExpiresAt),
		IssuedAt:       jwt.NewNumericDate(accessToken.IssuedAt),
		NotBefore:      jwt.NewNumericDate(accessToken.IssuedAt),
		ID:             accessToken.ID.String(),
		OrganizationID: accessToken.OrganizationID.String(),
		ClientID:       accessToken.ClientID,
		GrantID:        accessToken.GrantID.String(),
		Scope:          accessToken.Scopes.String(),
	})
	token.Header["typ"] = accessTokenType

	return token.SignedString(s.secretKey)
}

func (s *TokenSigner) Validate(ctx context.Context, rawToken string, allowedResources []string) (*mcpOAuthDomainOAuth2.AccessToken, error) {
	parsed, err := jwt.ParseWithClaims(rawToken, &accessTokenClaims{}, s.keyFunc,
		jwt.WithValidMethods([]string{signingMethod.Name}),
		jwt.WithIssuer(s.issuer),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
		jwt.WithLeeway(clockSkewLeeway),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, mcpOAuthApplicationTokens.ErrExpiredAccessToken
		}

		slog.DebugContext(ctx, "mcp access token validation failed",
			"component", "mcpoauth",
			"error", err.Error(),
		)
		return nil, mcpOAuthApplicationTokens.ErrInvalidAccessToken
	}

	claims, ok := parsed.Claims.(*accessTokenClaims)
	if !ok || !parsed.Valid {
		return nil, mcpOAuthApplicationTokens.ErrInvalidAccessToken
	}

	resource, err := resolveAudience(claims.Audience, allowedResources)
	if err != nil {
		return nil, err
	}

	return toAccessToken(claims, resource)
}

func (s *TokenSigner) keyFunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, mcpOAuthApplicationTokens.ErrInvalidAccessToken
	}
	return s.secretKey, nil
}

func resolveAudience(audience jwt.ClaimStrings, allowedResources []string) (string, error) {
	if len(audience) == 0 {
		return "", mcpOAuthApplicationTokens.ErrInvalidAccessTokenAudience
	}

	if len(allowedResources) == 0 {
		return audience[0], nil
	}

	for _, value := range audience {
		if slices.Contains(allowedResources, value) {
			return value, nil
		}
	}

	return "", mcpOAuthApplicationTokens.ErrInvalidAccessTokenAudience
}

func toAccessToken(claims *accessTokenClaims, resource string) (*mcpOAuthDomainOAuth2.AccessToken, error) {
	id, err := uuid.Parse(claims.ID)
	if err != nil {
		return nil, mcpOAuthApplicationTokens.ErrInvalidAccessToken
	}

	subject, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, mcpOAuthApplicationTokens.ErrInvalidAccessToken
	}

	organizationID, err := uuid.Parse(claims.OrganizationID)
	if err != nil {
		return nil, mcpOAuthApplicationTokens.ErrInvalidAccessToken
	}

	grantID, err := uuid.Parse(claims.GrantID)
	if err != nil {
		return nil, mcpOAuthApplicationTokens.ErrInvalidAccessToken
	}

	return &mcpOAuthDomainOAuth2.AccessToken{
		ID:             id,
		Subject:        subject,
		OrganizationID: organizationID,
		ClientID:       claims.ClientID,
		GrantID:        grantID,
		Scopes:         mcpOAuthDomainOAuth2.ParseScopes(claims.Scope),
		Resource:       resource,
		IssuedAt:       claims.IssuedAt.Time,
		ExpiresAt:      claims.ExpiresAt.Time,
	}, nil
}
