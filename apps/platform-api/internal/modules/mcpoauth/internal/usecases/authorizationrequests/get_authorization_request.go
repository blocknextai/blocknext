package authorizationrequests

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	mcpOAuthDomainAuthorizationRequests "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/authorizationrequests"
	mcpOAuthDomainClients "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/clients"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

type GetAuthorizationRequestQuery struct {
	AuthorizationRequestID uuid.UUID
}

type ClientResponse struct {
	ClientID           string  `json:"clientId"`
	Name               string  `json:"name"`
	LogoURI            *string `json:"logoUri,omitempty"`
	ClientURI          *string `json:"clientUri,omitempty"`
	IsMetadataDocument bool    `json:"isMetadataDocument"`
}

type ScopeResponse struct {
	Scope       string `json:"scope"`
	Description string `json:"description"`
}

type GetAuthorizationRequestResponse struct {
	ID          uuid.UUID        `json:"id"`
	Client      *ClientResponse  `json:"client"`
	Scopes      []*ScopeResponse `json:"scopes"`
	Resource    string           `json:"resource"`
	RedirectURI string           `json:"redirectUri"`
	Status      string           `json:"status"`
	IsPending   bool             `json:"isPending"`
	ExpiresAt   time.Time        `json:"expiresAt"`
}

func (s *Service) GetAuthorizationRequest(ctx context.Context, query *GetAuthorizationRequestQuery) (*GetAuthorizationRequestResponse, error) {
	request, err := s.authorizationRequestRepository.GetByID(ctx, query.AuthorizationRequestID)
	if err != nil {
		return nil, err
	}

	var client *mcpOAuthDomainClients.Client
	resolved, err := s.clientResolver.Resolve(ctx, request.ClientID)
	if err != nil {
		slog.WarnContext(ctx, "failed to resolve mcp oauth client",
			"component", "mcpoauth",
			"client_id", request.ClientID,
			"error", err.Error(),
		)
	} else {
		client = resolved
	}

	return MapAuthorizationRequestToResponse(request, client), nil
}

func MapAuthorizationRequestToResponse(
	request *mcpOAuthDomainAuthorizationRequests.AuthorizationRequest,
	client *mcpOAuthDomainClients.Client,
) *GetAuthorizationRequestResponse {
	return &GetAuthorizationRequestResponse{
		ID:          request.ID,
		Client:      MapClientToResponse(request.ClientID, client),
		Scopes:      MapScopesToResponse(request.Scopes),
		Resource:    request.Resource,
		RedirectURI: request.RedirectURI,
		Status:      request.Status.String(),
		IsPending:   request.IsPending(),
		ExpiresAt:   request.ExpiresAt,
	}
}

func MapClientToResponse(clientID string, client *mcpOAuthDomainClients.Client) *ClientResponse {
	response := &ClientResponse{
		ClientID:           clientID,
		Name:               clientID,
		IsMetadataDocument: mcpOAuthDomainClients.IsMetadataDocumentClientID(clientID),
	}

	if client != nil {
		response.Name = client.Name
		response.LogoURI = client.LogoURI
		response.ClientURI = client.ClientURI
	}

	return response
}

func MapScopesToResponse(scopes mcpOAuthDomainOAuth2.Scopes) []*ScopeResponse {
	responses := make([]*ScopeResponse, 0, len(scopes))
	for _, s := range scopes {
		responses = append(responses, &ScopeResponse{
			Scope:       s.String(),
			Description: s.Description(),
		})
	}
	return responses
}
