package grants

import (
	"context"

	mcpOAuthDomainGrants "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/grants"
	mcpOAuthDomainRefreshTokens "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/refreshtokens"
)

type GrantRevoker struct {
	grantRepository        mcpOAuthDomainGrants.GrantRepository
	refreshTokenRepository mcpOAuthDomainRefreshTokens.RefreshTokenRepository
}

func NewGrantRevoker(
	grantRepository mcpOAuthDomainGrants.GrantRepository,
	refreshTokenRepository mcpOAuthDomainRefreshTokens.RefreshTokenRepository,
) *GrantRevoker {
	return &GrantRevoker{
		grantRepository:        grantRepository,
		refreshTokenRepository: refreshTokenRepository,
	}
}

func (r *GrantRevoker) Revoke(ctx context.Context, grant *mcpOAuthDomainGrants.Grant) error {
	revoked, err := grant.Revoke()
	if err != nil {
		return err
	}

	if err := r.grantRepository.Update(ctx, revoked); err != nil {
		return ErrFailedToRevokeGrant.WithCause(err)
	}

	if err := r.refreshTokenRepository.RevokeAllByGrantID(ctx, revoked.ID, *revoked.RevokedAt); err != nil {
		return ErrFailedToRevokeGrant.WithCause(err)
	}

	return nil
}
