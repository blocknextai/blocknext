package tokens

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	mcpOAuthDomainClients "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/clients"
	mcpOAuthDomainGrants "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/grants"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

type RevokeTokenCommand struct {
	Token         string
	TokenTypeHint string
	ClientID      string
	ClientSecret  string
}

func (c *RevokeTokenCommand) Validate() error {
	if strings.TrimSpace(c.ClientID) == "" {
		return mcpOAuthDomainClients.ErrInvalidClientID
	}

	return nil
}

type RevokeTokenResponse struct{}

func (s *Service) RevokeToken(ctx context.Context, command *RevokeTokenCommand) (*RevokeTokenResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	client, err := s.clientResolver.Authenticate(ctx, command.ClientID, command.ClientSecret)
	if err != nil {
		return nil, err
	}

	err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		grantID, ok := s.resolveGrantID(txCtx, client, command.Token)
		if !ok {
			return nil
		}

		grant, err := s.grantRepository.GetByID(txCtx, grantID)
		if errors.Is(err, mcpOAuthDomainGrants.ErrGrantNotFound) {
			return nil
		}
		if err != nil {
			return err
		}

		if !grant.IsActive() {
			return nil
		}

		return s.grantRevoker.Revoke(txCtx, grant)
	})
	if err != nil {
		return nil, err
	}

	return &RevokeTokenResponse{}, nil
}

func (s *Service) resolveGrantID(ctx context.Context, client *mcpOAuthDomainClients.Client, token string) (uuid.UUID, bool) {
	if strings.TrimSpace(token) == "" {
		return uuid.Nil, false
	}

	refreshToken, err := s.refreshTokenRepository.GetByTokenHash(ctx, mcpOAuthDomainOAuth2.HashSecret(token))
	if err == nil {
		return refreshToken.GrantID, refreshToken.ClientID == client.ClientID
	}

	accessToken, err := s.accessTokenValidator.Validate(ctx, token, nil)
	if err != nil {
		return uuid.Nil, false
	}

	return accessToken.GrantID, accessToken.ClientID == client.ClientID
}
