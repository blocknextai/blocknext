package authorizationrequests

import (
	"context"

	"github.com/google/uuid"

	mcpOAuthApplicationAuthorizationRequests "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/authorizationrequests"
	mcpOAuthDomainAuthorizationRequests "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/authorizationrequests"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

type DenyAuthorizationRequestCommand struct {
	AuthorizationRequestID uuid.UUID
	UserID                 uuid.UUID
}

func (c *DenyAuthorizationRequestCommand) Validate() error {
	if c.AuthorizationRequestID == uuid.Nil {
		return mcpOAuthDomainAuthorizationRequests.ErrInvalidAuthorizationRequestID
	}

	if c.UserID == uuid.Nil {
		return mcpOAuthDomainAuthorizationRequests.ErrInvalidUserID
	}

	return nil
}

type DenyAuthorizationRequestResponse struct {
	RedirectURI string `json:"redirectUri"`
}

func (s *Service) DenyAuthorizationRequest(ctx context.Context, command *DenyAuthorizationRequestCommand) (*DenyAuthorizationRequestResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	var response *DenyAuthorizationRequestResponse

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		request, err := s.authorizationRequestRepository.GetByID(txCtx, command.AuthorizationRequestID)
		if err != nil {
			return err
		}

		denied, err := request.Deny(command.UserID)
		if err != nil {
			return err
		}

		if err := s.authorizationRequestRepository.Update(txCtx, denied); err != nil {
			return mcpOAuthApplicationAuthorizationRequests.ErrFailedToDenyAuthorizationRequest.WithCause(err)
		}

		state := ""
		if denied.State != nil {
			state = *denied.State
		}

		redirectURI, err := mcpOAuthDomainOAuth2.BuildRedirectURI(denied.RedirectURI, map[string]string{
			"error":             mcpOAuthDomainOAuth2.AccessDeniedError.String(),
			"error_description": mcpOAuthDomainAuthorizationRequests.ErrAuthorizationRequestAccessDenied.Error(),
			"state":             state,
			"iss":               s.issuerURL,
		})
		if err != nil {
			return err
		}

		response = &DenyAuthorizationRequestResponse{
			RedirectURI: redirectURI,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
