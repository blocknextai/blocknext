package grants

import (
	"context"

	"github.com/google/uuid"

	mcpOAuthDomainGrants "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/grants"
)

type RevokeGrantCommand struct {
	GrantID uuid.UUID
	UserID  uuid.UUID
}

func (c *RevokeGrantCommand) Validate() error {
	if c.GrantID == uuid.Nil {
		return mcpOAuthDomainGrants.ErrInvalidGrantID
	}

	if c.UserID == uuid.Nil {
		return mcpOAuthDomainGrants.ErrInvalidUserID
	}

	return nil
}

type RevokeGrantResponse struct {
	ID uuid.UUID `json:"id"`
}

func (s *Service) RevokeGrant(ctx context.Context, command *RevokeGrantCommand) (*RevokeGrantResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	var response *RevokeGrantResponse

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		grant, err := s.grantRepository.GetByID(txCtx, command.GrantID)
		if err != nil {
			return err
		}

		if grant.UserID != command.UserID || !grant.IsActive() {
			return mcpOAuthDomainGrants.ErrGrantNotFound
		}

		if err := s.grantRevoker.Revoke(txCtx, grant); err != nil {
			return err
		}

		response = &RevokeGrantResponse{
			ID: grant.ID,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
