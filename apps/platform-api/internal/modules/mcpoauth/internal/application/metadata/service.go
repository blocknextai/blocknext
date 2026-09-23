package metadata

import (
	"strings"

	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

type MetadataService interface {
	GetResourceURL(resourcePath string) string
	GetProtectedResourceMetadataURL(resourcePath string) string
	GetProtectedResourceMetadata(resourcePath string, scopes []string) *mcpOAuthDomainOAuth2.ProtectedResourceMetadata
}

type metadataService struct {
	issuerURL       string
	resourceBaseURL string
}

func NewMetadataService(issuerURL string, resourceBaseURL string) MetadataService {
	return &metadataService{
		issuerURL:       issuerURL,
		resourceBaseURL: strings.TrimSuffix(resourceBaseURL, "/"),
	}
}

func (s *metadataService) GetResourceURL(resourcePath string) string {
	return s.resourceBaseURL + resourcePath
}

func (s *metadataService) GetProtectedResourceMetadataURL(resourcePath string) string {
	return s.resourceBaseURL + mcpOAuthDomainOAuth2.ProtectedResourceMetadataPath + resourcePath
}

func (s *metadataService) GetProtectedResourceMetadata(
	resourcePath string,
	scopes []string,
) *mcpOAuthDomainOAuth2.ProtectedResourceMetadata {
	return &mcpOAuthDomainOAuth2.ProtectedResourceMetadata{
		Resource:               s.GetResourceURL(resourcePath),
		AuthorizationServers:   []string{s.issuerURL},
		ScopesSupported:        scopes,
		BearerMethodsSupported: []string{mcpOAuthDomainOAuth2.HeaderBearerMethod},
	}
}
