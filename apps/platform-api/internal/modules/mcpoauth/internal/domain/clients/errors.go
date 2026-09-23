package clients

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrClientNotFound                       = apperror.NotFound("oauth client not found")
	ErrClientAlreadyExists                  = apperror.Conflict("oauth client already exists")
	ErrInvalidClientName                    = apperror.Validation("invalid client name")
	ErrInvalidClientID                      = apperror.Validation("invalid client id")
	ErrInvalidRedirectURI                   = apperror.Validation("invalid redirect uri")
	ErrInvalidGrantTypes                    = apperror.Validation("invalid grant types")
	ErrInvalidResponseTypes                 = apperror.Validation("invalid response types")
	ErrInvalidScopes                        = apperror.Validation("invalid scopes")
	ErrInvalidTokenEndpointAuthMethod       = apperror.Validation("invalid token endpoint auth method")
	ErrUnauthorizedGrantType                = apperror.Validation("client is not allowed to use this grant type")
	ErrInvalidClientMetadataDocument        = apperror.Validation("invalid client id metadata document")
	ErrClientMetadataDocumentHostNotAllowed = apperror.Validation("client id metadata document host is not publicly routable")
	ErrInvalidClientCredentials             = apperror.Unauthorized("invalid client credentials")
)
