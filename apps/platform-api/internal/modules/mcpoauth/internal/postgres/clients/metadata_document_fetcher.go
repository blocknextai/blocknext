package clients

import (
	"context"
	"log/slog"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/cache"
	"github.com/blocknextai/go-packages/digest"
	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/go-packages/json"
	mcpOAuthDomainClients "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/clients"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

const (
	metadataDocumentCacheKeyPrefix = "mcpoauth:clients:metadata_document:"
	metadataDocumentCacheTTL       = 15 * time.Minute
	metadataDocumentFetchTimeout   = 5 * time.Second
	carrierGradeNATFirstOctet      = 100
	carrierGradeNATSecondOctetMin  = 64
	carrierGradeNATSecondOctetMax  = 127
)

type metadataDocument struct {
	ClientID      string   `json:"client_id"`
	ClientName    string   `json:"client_name"`
	RedirectURIs  []string `json:"redirect_uris"`
	GrantTypes    []string `json:"grant_types"`
	ResponseTypes []string `json:"response_types"`
	Scope         string   `json:"scope"`
	LogoURI       string   `json:"logo_uri"`
	ClientURI     string   `json:"client_uri"`
}

type MetadataDocumentFetcher struct {
	cacheService cache.Service
}

func NewMetadataDocumentFetcher(cacheService cache.Service) mcpOAuthDomainClients.MetadataDocumentFetcher {
	return &MetadataDocumentFetcher{
		cacheService: cacheService,
	}
}

func (f *MetadataDocumentFetcher) Fetch(ctx context.Context, clientID string) (*mcpOAuthDomainClients.Client, error) {
	if !mcpOAuthDomainClients.IsMetadataDocumentClientID(clientID) || strings.Contains(clientID, "#") {
		return nil, mcpOAuthDomainClients.ErrInvalidClientID
	}

	document, err := f.load(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if document.ClientID != clientID {
		return nil, mcpOAuthDomainClients.ErrInvalidClientMetadataDocument
	}

	name := strings.TrimSpace(document.ClientName)
	if name == "" {
		name = clientID
	}

	return mcpOAuthDomainClients.New(
		clientID,
		nil,
		name,
		document.RedirectURIs,
		toGrantTypes(document.GrantTypes),
		[]mcpOAuthDomainOAuth2.ResponseType{mcpOAuthDomainOAuth2.CodeResponseType},
		mcpOAuthDomainOAuth2.ParseScopes(document.Scope),
		mcpOAuthDomainOAuth2.NoneAuthentication,
		optionalString(document.LogoURI),
		optionalString(document.ClientURI),
	)
}

func (f *MetadataDocumentFetcher) load(ctx context.Context, clientID string) (*metadataDocument, error) {
	cacheKey := metadataDocumentCacheKeyPrefix + digest.SHA256Hex(clientID)

	if cached, err := f.cacheService.Get(ctx, cacheKey); err == nil && strings.TrimSpace(cached) != "" {
		document := new(metadataDocument)
		if err := json.Unmarshal([]byte(cached), document); err == nil {
			return document, nil
		}
	}

	client := httpclient.NewClientBuilder().
		Context(ctx).
		Timeout(metadataDocumentFetchTimeout).
		JSONContentType().
		AllowDestination(allowPublicHTTPS).
		Build()

	document := new(metadataDocument)
	response, err := client.Get(clientID).Do(document, nil)
	if err != nil {
		return nil, mcpOAuthDomainClients.ErrInvalidClientMetadataDocument.WithCause(err)
	}
	if !response.IsSuccess() {
		return nil, mcpOAuthDomainClients.ErrInvalidClientMetadataDocument
	}

	if encoded, err := json.Marshal(document); err == nil {
		if err := f.cacheService.Set(ctx, cacheKey, string(encoded), metadataDocumentCacheTTL); err != nil {
			slog.WarnContext(ctx, "failed to cache client id metadata document",
				"component", "mcpoauth",
				"client_id", clientID,
				"error", err.Error(),
			)
		}
	}

	return document, nil
}

func allowPublicHTTPS(target *url.URL) error {
	if target.Scheme != "https" || strings.TrimSpace(target.Hostname()) == "" {
		return mcpOAuthDomainClients.ErrClientMetadataDocumentHostNotAllowed
	}

	addresses, err := net.LookupIP(target.Hostname())
	if err != nil {
		return mcpOAuthDomainClients.ErrClientMetadataDocumentHostNotAllowed.WithCause(err)
	}

	for _, address := range addresses {
		if !isPublicIP(address) {
			return mcpOAuthDomainClients.ErrClientMetadataDocumentHostNotAllowed
		}
	}

	return nil
}

func isPublicIP(address net.IP) bool {
	if address.IsLoopback() || address.IsPrivate() || address.IsUnspecified() ||
		address.IsLinkLocalUnicast() || address.IsLinkLocalMulticast() ||
		address.IsInterfaceLocalMulticast() || address.IsMulticast() {
		return false
	}

	if v4 := address.To4(); v4 != nil &&
		v4[0] == carrierGradeNATFirstOctet &&
		v4[1] >= carrierGradeNATSecondOctetMin &&
		v4[1] <= carrierGradeNATSecondOctetMax {
		return false
	}

	return true
}

func toGrantTypes(values []string) []mcpOAuthDomainOAuth2.GrantType {
	if len(values) == 0 {
		return []mcpOAuthDomainOAuth2.GrantType{
			mcpOAuthDomainOAuth2.AuthorizationCodeGrant,
			mcpOAuthDomainOAuth2.RefreshTokenGrant,
		}
	}

	return mcpOAuthDomainOAuth2.SupportedGrantTypes(values)
}

func optionalString(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return new(trimmed)
}
