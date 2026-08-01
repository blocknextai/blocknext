package exchangecode

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
	"strings"

	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/go-packages/json"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonDomainOAuth2 "github.com/blocknextai/platform-api/internal/common/domain/oauth2"
	credentialOAuthApplicationOAuth2 "github.com/blocknextai/platform-api/internal/credentialoauth/application/oauth2"
	credentialOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/credentialoauth/domain/oauth2"
	credentialsApplicationCredentials "github.com/blocknextai/platform-api/internal/credentials/application/credentials"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/credentials/domain/credentials"
	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/nodeengine/application/credentials"
	platformApplicationPlatformCredentials "github.com/blocknextai/platform-api/internal/platform/application/platformcredentials"
	"github.com/google/uuid"
)

type Handler struct {
	stateStore                  *credentialOAuthApplicationOAuth2.StateStore
	oauth2RedirectURL           string
	credentialService           credentialsApplicationCredentials.CredentialService
	nodeEngineCredentialService nodeEngineApplicationCredentials.CredentialService
	platformCredentialService   platformApplicationPlatformCredentials.PlatformCredentialService
}

func New(
	stateStore *credentialOAuthApplicationOAuth2.StateStore,
	oauth2RedirectURL string,
	credentialService credentialsApplicationCredentials.CredentialService,
	nodeEngineCredentialService nodeEngineApplicationCredentials.CredentialService,
	platformCredentialService platformApplicationPlatformCredentials.PlatformCredentialService,
) *Handler {
	return &Handler{
		stateStore:                  stateStore,
		oauth2RedirectURL:           oauth2RedirectURL,
		credentialService:           credentialService,
		nodeEngineCredentialService: nodeEngineCredentialService,
		platformCredentialService:   platformCredentialService,
	}
}

func (h *Handler) Handle(ctx context.Context, request *ExchangeCodeQuery) (*ExchangeCodeResponse, error) {
	state, err := h.stateStore.Consume(ctx, request.State)
	if err != nil {
		return nil, err
	}

	credentialInfo, err := h.credentialService.GetByIDForOwner(ctx, state.CredentialID, state.OwnerType, state.OwnerID)
	if err != nil {
		return nil, err
	}

	var credentialData commonDomainOAuth2.Credential
	if err := json.ArgsToStruct(credentialInfo.Data, &credentialData); err != nil {
		return nil, err
	}

	isPlatformCredential := credentialInfo.SourceType == credentialsDomainCredentials.SourceTypePlatform
	if isPlatformCredential {
		platformCred, ok := h.platformCredentialService.GetPlatformCredential(credentialInfo.Key)
		if !ok || platformCred == nil {
			return nil, credentialOAuthApplicationOAuth2.ErrPlatformCredentialMissing
		}
		if err := json.ArgsToStruct(platformCred.Data, &credentialData); err != nil {
			return nil, err
		}
		if !credentialData.HasClientCredentials() {
			return nil, credentialOAuthApplicationOAuth2.ErrPlatformCredentialMissingClientCreds
		}
	}

	schemaData := h.nodeEngineCredentialService.GetHiddenObjectsByCredentialID(credentialInfo.Key)
	if schemaData == nil {
		return nil, credentialOAuthApplicationOAuth2.ErrInvalidCredentialData
	}

	tokenURL, ok := schemaData["tokenUrl"].(string)
	if !ok || strings.TrimSpace(tokenURL) == "" {
		return nil, credentialOAuthApplicationOAuth2.ErrInvalidCredentialData
	}

	token, err := h.exchangeCodeForToken(ctx, credentialData, schemaData, request.Code, state.CodeVerifier, h.oauth2RedirectURL)
	if err != nil {
		return nil, err
	}

	if err := h.saveToken(ctx, state.OwnerType, state.OwnerID, state.CredentialID, credentialData, token, isPlatformCredential); err != nil {
		return nil, err
	}

	return MapToResponse("success"), nil
}

func (h *Handler) saveToken(ctx context.Context, ownerType commonDomain.OwnerType, ownerID uuid.UUID, credentialID uuid.UUID, credentialData commonDomainOAuth2.Credential, token *commonDomainOAuth2.Token, isPlatformCredential bool) error {
	var dataToSave any
	if isPlatformCredential {
		dataToSave = commonDomainOAuth2.Credential{OAuthToken: *token}
	} else {
		credentialData.OAuthToken = *token
		dataToSave = credentialData
	}
	return h.credentialService.SaveCredentialForOwner(ctx, credentialID, ownerType, ownerID, dataToSave)
}

func (h *Handler) exchangeCodeForToken(
	ctx context.Context,
	credential commonDomainOAuth2.Credential,
	schemaData map[string]any,
	code string,
	codeVerifier string,
	redirectURL string,
) (*commonDomainOAuth2.Token, error) {
	scope, _ := schemaData["scope"].(string)
	tokenURL, _ := schemaData["tokenUrl"].(string)
	authentication, _ := schemaData["authentication"].(string)

	form := url.Values{}
	form.Set("grant_type", credentialOAuthDomainOAuth2.AuthorizationCodeGrant.String())
	form.Set("code", code)
	form.Set("redirect_uri", redirectURL)
	form.Set("code_verifier", codeVerifier)

	if scope != "" {
		form.Set("scope", scope)
	}

	clientBuilder := httpclient.NewClientBuilder().
		Context(ctx).
		FormUrlencodedContentType()

	if authentication == credentialOAuthDomainOAuth2.HeaderAuthentication.String() && credential.HasClientCredentials() {
		clientBuilder.BasicAuth(credential.GetClientIdentifier(), credential.ClientSecret)
	} else {
		h.addClientCredentialsToForm(&form, credential)
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
		slog.WarnContext(ctx, "oauth token exchange failed",
			"providerError", errorResp.Error,
			"providerErrorDescription", errorResp.ErrorDescription,
		)
		if errorResp.Error != "" {
			return nil, credentialOAuthApplicationOAuth2.ErrTokenExchangeFailed.WithCause(errors.New(errorResp.GetErrorMessage()))
		}
		return nil, credentialOAuthApplicationOAuth2.ErrTokenExchangeFailed
	}

	return tokenResp.ToToken(), nil
}

func (h *Handler) addClientCredentialsToForm(form *url.Values, credential commonDomainOAuth2.Credential) {
	if credential.ClientID != "" {
		form.Set("client_id", credential.ClientID)
	}
	if credential.ClientKey != "" {
		form.Set("client_key", credential.ClientKey)
	}
	if credential.ClientSecret != "" {
		form.Set("client_secret", credential.ClientSecret)
	}
}
