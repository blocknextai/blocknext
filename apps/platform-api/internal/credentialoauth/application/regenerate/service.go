package regenerate

import (
	"context"
	"errors"
	"log/slog"
	"maps"
	"net/url"
	"time"

	"github.com/blocknextai/go-packages/cache"
	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/go-packages/json"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonDomainOAuth2 "github.com/blocknextai/platform-api/internal/common/domain/oauth2"
	credentialOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/credentialoauth/domain/oauth2"
	credentialsApplicationCredentials "github.com/blocknextai/platform-api/internal/credentials/application/credentials"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/credentials/domain/credentials"
	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/nodeengine/application/credentials"
	platformApplicationPlatformCredentials "github.com/blocknextai/platform-api/internal/platform/application/platformcredentials"
	"github.com/google/uuid"
)

const (
	refreshLockKeyPrefix    = "credentialoauth:oauth2:refresh:"
	refreshLockTTL          = 30 * time.Second
	refreshLockWaitTimeout  = 10 * time.Second
	refreshLockRetryBackoff = 200 * time.Millisecond
)

type CredentialOAuthTokenRegenerateService interface {
	RegenerateTokenIfNeeded(ctx context.Context, ownerType commonDomain.OwnerType, ownerID uuid.UUID, credentialID uuid.UUID) (map[string]any, error)
}

type credentialOAuthTokenRegenerateService struct {
	credentialService           credentialsApplicationCredentials.CredentialService
	nodeEngineCredentialService nodeEngineApplicationCredentials.CredentialService
	platformCredentialService   platformApplicationPlatformCredentials.PlatformCredentialService
	cache                       cache.Service
}

func NewCredentialOAuthTokenRegenerateService(
	credentialService credentialsApplicationCredentials.CredentialService,
	nodeEngineCredentialService nodeEngineApplicationCredentials.CredentialService,
	platformCredentialService platformApplicationPlatformCredentials.PlatformCredentialService,
	cacheService cache.Service,
) CredentialOAuthTokenRegenerateService {
	return &credentialOAuthTokenRegenerateService{
		credentialService:           credentialService,
		nodeEngineCredentialService: nodeEngineCredentialService,
		platformCredentialService:   platformCredentialService,
		cache:                       cacheService,
	}
}

func (s *credentialOAuthTokenRegenerateService) RegenerateTokenIfNeeded(ctx context.Context, ownerType commonDomain.OwnerType, ownerID uuid.UUID, credentialID uuid.UUID) (map[string]any, error) {
	credentialInfo, err := s.credentialService.GetByIDForOwner(ctx, credentialID, ownerType, ownerID)
	if err != nil {
		return nil, err
	}

	if !s.nodeEngineCredentialService.IsOAuth2Credential(credentialInfo.Key) {
		return s.resolvePlatformCredentialData(credentialInfo)
	}

	return s.regenerateOAuth2Token(ctx, ownerType, ownerID, credentialID, credentialInfo)
}

func (s *credentialOAuthTokenRegenerateService) resolvePlatformCredentialData(credentialInfo *credentialsDomainCredentials.CredentialInfo) (map[string]any, error) {
	if credentialInfo.SourceType != credentialsDomainCredentials.SourceTypePlatform {
		return credentialInfo.Data, nil
	}

	platformCred, ok := s.platformCredentialService.GetPlatformCredential(credentialInfo.Key)
	if !ok || platformCred == nil {
		return credentialInfo.Data, ErrPlatformCredentialMissing
	}

	result := make(map[string]any)
	maps.Copy(result, platformCred.Data)
	maps.Copy(result, credentialInfo.Data)

	return result, nil
}

func (s *credentialOAuthTokenRegenerateService) regenerateOAuth2Token(
	ctx context.Context,
	ownerType commonDomain.OwnerType,
	ownerID uuid.UUID,
	credentialID uuid.UUID,
	credentialInfo *credentialsDomainCredentials.CredentialInfo,
) (map[string]any, error) {
	credentialData, _, err := s.parseCredentialData(credentialInfo)
	if err != nil {
		return credentialInfo.Data, err
	}

	if !credentialData.OAuthToken.HasRefreshToken() || !credentialData.OAuthToken.NeedsRefresh() {
		return credentialInfo.Data, nil
	}

	lockKey := refreshLockKeyPrefix + credentialID.String()
	acquired, err := s.acquireRefreshLock(ctx, lockKey)
	if err != nil {
		return credentialInfo.Data, err
	}

	if !acquired {
		return s.readLatestCredentialData(ctx, ownerType, ownerID, credentialID, credentialInfo.Data)
	}

	defer func() {
		if _, releaseErr := s.cache.ReleaseSemaphoreAtomic(ctx, lockKey); releaseErr != nil {
			slog.WarnContext(ctx, "failed to release refresh lock", "credentialID", credentialID, "error", releaseErr)
		}
	}()

	freshInfo, err := s.credentialService.GetByIDForOwner(ctx, credentialID, ownerType, ownerID)
	if err != nil {
		return credentialInfo.Data, err
	}

	freshData, isPlatformCredential, err := s.parseCredentialData(freshInfo)
	if err != nil {
		return freshInfo.Data, err
	}

	if !freshData.OAuthToken.HasRefreshToken() || !freshData.OAuthToken.NeedsRefresh() {
		return freshInfo.Data, nil
	}

	oldRefreshToken := freshData.OAuthToken.RefreshToken

	newToken, err := s.getNewToken(ctx, freshInfo.Key, freshData)
	if err != nil {
		return freshInfo.Data, err
	}

	if newToken.RefreshToken != "" && newToken.RefreshToken != oldRefreshToken {
		slog.InfoContext(ctx, "oauth refresh token rotated",
			"credentialID", credentialID,
			"credentialKey", freshInfo.Key,
		)
	}

	if err := s.saveToken(ctx, ownerType, ownerID, credentialID, freshData, newToken, isPlatformCredential); err != nil {
		return freshInfo.Data, err
	}

	freshData.OAuthToken = *newToken
	var result map[string]any
	if err := json.ArgsToStruct(freshData, &result); err != nil {
		return freshInfo.Data, err
	}

	return result, nil
}

func (s *credentialOAuthTokenRegenerateService) parseCredentialData(credentialInfo *credentialsDomainCredentials.CredentialInfo) (commonDomainOAuth2.Credential, bool, error) {
	var credentialData commonDomainOAuth2.Credential
	if err := json.ArgsToStruct(credentialInfo.Data, &credentialData); err != nil {
		return credentialData, false, err
	}

	isPlatformCredential := credentialInfo.SourceType == credentialsDomainCredentials.SourceTypePlatform
	if isPlatformCredential {
		platformCred, ok := s.platformCredentialService.GetPlatformCredential(credentialInfo.Key)
		if !ok || platformCred == nil {
			return credentialData, isPlatformCredential, ErrPlatformCredentialMissing
		}
		if err := json.ArgsToStruct(platformCred.Data, &credentialData); err != nil {
			return credentialData, isPlatformCredential, err
		}
	}

	return credentialData, isPlatformCredential, nil
}

func (s *credentialOAuthTokenRegenerateService) acquireRefreshLock(ctx context.Context, key string) (bool, error) {
	deadline := time.Now().Add(refreshLockWaitTimeout)

	for {
		acquired, err := s.cache.AcquireSemaphoreAtomic(ctx, key, 1, refreshLockTTL)
		if err != nil {
			return false, err
		}
		if acquired {
			return true, nil
		}

		if time.Now().After(deadline) {
			return false, nil
		}

		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(refreshLockRetryBackoff):
		}
	}
}

func (s *credentialOAuthTokenRegenerateService) readLatestCredentialData(
	ctx context.Context,
	ownerType commonDomain.OwnerType,
	ownerID uuid.UUID,
	credentialID uuid.UUID,
	fallback map[string]any,
) (map[string]any, error) {
	freshInfo, err := s.credentialService.GetByIDForOwner(ctx, credentialID, ownerType, ownerID)
	if err != nil {
		return fallback, err
	}
	return freshInfo.Data, nil
}

func (s *credentialOAuthTokenRegenerateService) getNewToken(
	ctx context.Context,
	credentialKey string,
	credentialData commonDomainOAuth2.Credential,
) (*commonDomainOAuth2.Token, error) {
	if refreshable, ok := s.nodeEngineCredentialService.GetRefreshableCredential(credentialKey); ok {
		return refreshable.RefreshToken(ctx, credentialData)
	}

	schemaData := s.nodeEngineCredentialService.GetHiddenObjectsByCredentialID(credentialKey)
	if schemaData == nil {
		return nil, nil
	}

	tokenURL, ok := schemaData["tokenUrl"].(string)
	if !ok {
		return nil, ErrInvalidTokenURL
	}

	authentication, _ := schemaData["authentication"].(string)
	scope, _ := schemaData["scope"].(string)

	return s.refreshOAuthToken(ctx, credentialData, authentication, tokenURL, scope)
}

func (s *credentialOAuthTokenRegenerateService) saveToken(
	ctx context.Context,
	ownerType commonDomain.OwnerType,
	ownerID uuid.UUID,
	credentialID uuid.UUID,
	credentialData commonDomainOAuth2.Credential,
	newToken *commonDomainOAuth2.Token,
	isPlatformCredential bool,
) error {
	var dataToSave any
	if isPlatformCredential {
		dataToSave = commonDomainOAuth2.Credential{OAuthToken: *newToken}
	} else {
		credentialData.OAuthToken = *newToken
		dataToSave = credentialData
	}
	return s.credentialService.SaveCredentialForOwner(ctx, credentialID, ownerType, ownerID, dataToSave)
}

func (s *credentialOAuthTokenRegenerateService) refreshOAuthToken(
	ctx context.Context,
	credential commonDomainOAuth2.Credential,
	authentication string,
	tokenURL string,
	scope string,
) (*commonDomainOAuth2.Token, error) {
	form := url.Values{}
	form.Set("grant_type", credentialOAuthDomainOAuth2.RefreshTokenGrant.String())
	form.Set("refresh_token", credential.OAuthToken.RefreshToken)

	if credential.ClientKey != "" {
		form.Set("client_key", credential.ClientKey)
	}
	if credential.ClientID != "" {
		form.Set("client_id", credential.ClientID)
	}
	if credential.ClientSecret != "" {
		form.Set("client_secret", credential.ClientSecret)
	}
	if scope != "" {
		form.Set("scope", scope)
	}

	clientBuilder := httpclient.NewClientBuilder().
		Context(ctx).
		FormUrlencodedContentType()

	if authentication == credentialOAuthDomainOAuth2.HeaderAuthentication.String() {
		clientBuilder.BasicAuth(credential.GetClientIdentifier(), credential.ClientSecret)
	}

	var tokenResp credentialOAuthDomainOAuth2.TokenResponse
	var errorResp credentialOAuthDomainOAuth2.TokenErrorResponse
	resp, err := clientBuilder.Build().Post(tokenURL).
		FormUrlencodedContentType().
		Body(form).
		Do(&tokenResp, &errorResp)

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		slog.WarnContext(ctx, "oauth token refresh failed",
			"providerError", errorResp.Error,
			"providerErrorDescription", errorResp.ErrorDescription,
		)
		if errorResp.Error == credentialOAuthDomainOAuth2.InvalidGrantError {
			return nil, ErrRefreshTokenInvalid.WithCause(errors.New(errorResp.GetErrorMessage()))
		}
		if errorResp.Error != "" {
			return nil, ErrTokenRefreshFailed.WithCause(errors.New(errorResp.GetErrorMessage()))
		}
		return nil, ErrTokenRefreshFailed
	}

	token := &commonDomainOAuth2.Token{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: credential.OAuthToken.RefreshToken,
	}
	if tokenResp.RefreshToken != "" {
		token.RefreshToken = tokenResp.RefreshToken
	}
	token.SetExpiration(tokenResp.ExpiresIn)

	return token, nil
}
