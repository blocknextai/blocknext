package authorizationrequests

import (
	"context"
	"errors"
	"strings"

	mcpOAuthApplicationAuthorizationRequests "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/authorizationrequests"
	mcpOAuthDomainAuthorizationRequests "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/authorizationrequests"
	mcpOAuthDomainClients "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/clients"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

type CreateAuthorizationRequestCommand struct {
	ClientID            string
	RedirectURI         string
	ResponseType        string
	Scope               string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
	Resource            string
}

func (c *CreateAuthorizationRequestCommand) Validate() error {
	if strings.TrimSpace(c.ClientID) == "" {
		return mcpOAuthDomainClients.ErrInvalidClientID
	}

	if strings.TrimSpace(c.RedirectURI) == "" {
		return mcpOAuthDomainClients.ErrInvalidRedirectURI
	}

	return nil
}

type CreateAuthorizationRequestResponse struct {
	RedirectURI string
}

const (
	requestIDPlaceholder = "{requestId}"
)

func (s *Service) CreateAuthorizationRequest(ctx context.Context, command *CreateAuthorizationRequestCommand) (*CreateAuthorizationRequestResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	client, err := s.clientResolver.Resolve(ctx, command.ClientID)
	if err != nil {
		return nil, err
	}

	if !client.HasRedirectURI(command.RedirectURI) {
		return nil, mcpOAuthDomainClients.ErrInvalidRedirectURI
	}

	var response *CreateAuthorizationRequestResponse

	err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		request, err := s.newAuthorizationRequest(client, command)
		if err != nil {
			return err
		}

		if err := s.authorizationRequestRepository.Create(txCtx, request); err != nil {
			return mcpOAuthApplicationAuthorizationRequests.ErrFailedToCreateAuthorizationRequest.WithCause(err)
		}

		response = &CreateAuthorizationRequestResponse{
			RedirectURI: strings.ReplaceAll(s.consentURLTemplate, requestIDPlaceholder, request.ID.String()),
		}
		return nil
	})
	if err == nil {
		return response, nil
	}

	errorCode, ok := errorCodeOf(err)
	if !ok {
		return nil, err
	}

	redirectURI, err := mcpOAuthDomainOAuth2.BuildRedirectURI(command.RedirectURI, map[string]string{
		"error":             errorCode.String(),
		"error_description": err.Error(),
		"state":             command.State,
		"iss":               s.issuerURL,
	})
	if err != nil {
		return nil, err
	}

	return &CreateAuthorizationRequestResponse{
		RedirectURI: redirectURI,
	}, nil
}

func (s *Service) newAuthorizationRequest(
	client *mcpOAuthDomainClients.Client,
	command *CreateAuthorizationRequestCommand,
) (*mcpOAuthDomainAuthorizationRequests.AuthorizationRequest, error) {
	if !mcpOAuthDomainOAuth2.ResponseType(command.ResponseType).IsValid() {
		return nil, mcpOAuthDomainAuthorizationRequests.ErrUnsupportedResponseType
	}

	if !client.SupportsGrantType(mcpOAuthDomainOAuth2.AuthorizationCodeGrant) {
		return nil, mcpOAuthDomainClients.ErrUnauthorizedGrantType
	}

	resource, ok := mcpOAuthDomainOAuth2.ResolveResource(s.resourceBaseURL, command.Resource)
	if !ok {
		return nil, mcpOAuthDomainAuthorizationRequests.ErrInvalidResource
	}

	requestedScopes := mcpOAuthDomainOAuth2.ParseScopes(command.Scope)
	if len(requestedScopes) == 0 {
		requestedScopes = mcpOAuthDomainOAuth2.DefaultScopes
	}
	if !requestedScopes.IsValid() {
		return nil, mcpOAuthDomainAuthorizationRequests.ErrInvalidScopes
	}

	return mcpOAuthDomainAuthorizationRequests.New(
		client.ClientID,
		command.RedirectURI,
		client.AllowedScopes(requestedScopes),
		optionalString(command.State),
		command.CodeChallenge,
		mcpOAuthDomainOAuth2.CodeChallengeMethod(command.CodeChallengeMethod),
		resource,
		s.authorizationRequestTTL,
	)
}

func errorCodeOf(err error) (mcpOAuthDomainOAuth2.ErrorCode, bool) {
	switch {
	case errors.Is(err, mcpOAuthDomainAuthorizationRequests.ErrUnsupportedResponseType):
		return mcpOAuthDomainOAuth2.UnsupportedResponseTypeError, true
	case errors.Is(err, mcpOAuthDomainClients.ErrUnauthorizedGrantType):
		return mcpOAuthDomainOAuth2.UnauthorizedClientError, true
	case errors.Is(err, mcpOAuthDomainAuthorizationRequests.ErrInvalidCodeChallenge),
		errors.Is(err, mcpOAuthDomainAuthorizationRequests.ErrInvalidCodeChallengeMethod):
		return mcpOAuthDomainOAuth2.InvalidRequestError, true
	case errors.Is(err, mcpOAuthDomainAuthorizationRequests.ErrInvalidResource):
		return mcpOAuthDomainOAuth2.InvalidTargetError, true
	case errors.Is(err, mcpOAuthDomainAuthorizationRequests.ErrInvalidScopes):
		return mcpOAuthDomainOAuth2.InvalidScopeError, true
	default:
		return "", false
	}
}

func optionalString(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return new(trimmed)
}
