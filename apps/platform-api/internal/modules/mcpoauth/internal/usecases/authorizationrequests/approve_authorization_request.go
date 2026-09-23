package authorizationrequests

import (
	"context"
	"errors"

	"github.com/google/uuid"

	mcpOAuthApplicationAuthorizationRequests "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/authorizationrequests"
	mcpOAuthApplicationGrants "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/grants"
	mcpOAuthDomainAuthorizationCodes "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/authorizationcodes"
	mcpOAuthDomainAuthorizationRequests "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/authorizationrequests"
	mcpOAuthDomainGrants "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/grants"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

type ApproveAuthorizationRequestCommand struct {
	AuthorizationRequestID uuid.UUID
	UserID                 uuid.UUID
	OrganizationID         uuid.UUID
	Scopes                 []string
}

func (c *ApproveAuthorizationRequestCommand) Validate() error {
	if c.AuthorizationRequestID == uuid.Nil {
		return mcpOAuthDomainAuthorizationRequests.ErrInvalidAuthorizationRequestID
	}

	if c.UserID == uuid.Nil {
		return mcpOAuthDomainAuthorizationRequests.ErrInvalidUserID
	}

	if c.OrganizationID == uuid.Nil {
		return mcpOAuthDomainAuthorizationRequests.ErrInvalidOrganizationID
	}

	if !mcpOAuthDomainOAuth2.ScopesFromStrings(c.Scopes).IsValid() {
		return mcpOAuthDomainAuthorizationRequests.ErrInvalidScopes
	}

	return nil
}

type ApproveAuthorizationRequestResponse struct {
	RedirectURI string `json:"redirectUri"`
}

func (s *Service) ApproveAuthorizationRequest(ctx context.Context, command *ApproveAuthorizationRequestCommand) (*ApproveAuthorizationRequestResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	var response *ApproveAuthorizationRequestResponse

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		request, err := s.authorizationRequestRepository.GetByID(txCtx, command.AuthorizationRequestID)
		if err != nil {
			return err
		}

		if _, err := s.organizationUserService.GetByOrganizationIDAndUserID(txCtx, command.OrganizationID, command.UserID); err != nil {
			return mcpOAuthDomainAuthorizationRequests.ErrOrganizationMembershipRequired.WithCause(err)
		}

		scopes := request.Scopes
		if len(command.Scopes) > 0 {
			scopes = mcpOAuthDomainOAuth2.ScopesFromStrings(command.Scopes)
		}

		approved, err := request.Approve(command.UserID, command.OrganizationID, scopes)
		if err != nil {
			return err
		}

		grant, err := s.saveGrant(txCtx, approved)
		if err != nil {
			return err
		}

		codeSecret, err := mcpOAuthDomainOAuth2.GenerateSecret(mcpOAuthDomainAuthorizationCodes.CodePrefix)
		if err != nil {
			return mcpOAuthApplicationAuthorizationRequests.ErrFailedToApproveAuthorizationRequest.WithCause(err)
		}

		code, err := mcpOAuthDomainAuthorizationCodes.New(
			grant.ID,
			approved.ClientID,
			codeSecret.Hash,
			approved.RedirectURI,
			approved.CodeChallenge,
			approved.CodeChallengeMethod,
			approved.Resource,
			approved.Scopes,
			s.authorizationCodeTTL,
		)
		if err != nil {
			return err
		}

		if err := s.authorizationCodeRepository.Create(txCtx, code); err != nil {
			return mcpOAuthApplicationAuthorizationRequests.ErrFailedToApproveAuthorizationRequest.WithCause(err)
		}

		if err := s.authorizationRequestRepository.Update(txCtx, approved); err != nil {
			return mcpOAuthApplicationAuthorizationRequests.ErrFailedToApproveAuthorizationRequest.WithCause(err)
		}

		redirectURI, err := mcpOAuthDomainOAuth2.BuildRedirectURI(approved.RedirectURI, map[string]string{
			"code":  codeSecret.Plain,
			"state": stringValue(approved.State),
			"iss":   s.issuerURL,
		})
		if err != nil {
			return err
		}

		response = &ApproveAuthorizationRequestResponse{
			RedirectURI: redirectURI,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) saveGrant(
	ctx context.Context,
	request *mcpOAuthDomainAuthorizationRequests.AuthorizationRequest,
) (*mcpOAuthDomainGrants.Grant, error) {
	existing, err := s.grantRepository.GetActiveByClientIDAndUserIDAndOrganizationIDAndResource(
		ctx,
		request.ClientID,
		*request.UserID,
		*request.OrganizationID,
		request.Resource,
	)
	if err != nil && !errors.Is(err, mcpOAuthDomainGrants.ErrGrantNotFound) {
		return nil, err
	}

	if existing != nil {
		updated, err := existing.Update(request.Scopes)
		if err != nil {
			return nil, err
		}

		if err := s.grantRepository.Update(ctx, updated); err != nil {
			return nil, mcpOAuthApplicationGrants.ErrFailedToUpdateGrant.WithCause(err)
		}

		return updated, nil
	}

	grant, err := mcpOAuthDomainGrants.New(
		request.ClientID,
		*request.UserID,
		*request.OrganizationID,
		request.Scopes,
		request.Resource,
	)
	if err != nil {
		return nil, err
	}

	if err := s.grantRepository.Create(ctx, grant); err != nil {
		return nil, mcpOAuthApplicationGrants.ErrFailedToCreateGrant.WithCause(err)
	}

	return grant, nil
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
