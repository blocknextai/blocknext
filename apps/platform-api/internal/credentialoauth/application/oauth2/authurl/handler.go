package authurl

import (
	"context"
	"net/url"

	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/go-packages/pkce"
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

type Credential struct {
	ID         uuid.UUID
	Key        string
	SourceType credentialsDomainCredentials.SourceType
	Data       commonDomainOAuth2.Credential
}

func (h *Handler) Handle(ctx context.Context, command *AuthURLCommand) (*AuthURLResponse, error) {
	credential, err := h.getCredentialData(ctx, command.OwnerType, command.OwnerID, command.CredentialID)
	if err != nil {
		return nil, err
	}

	schemaData := h.nodeEngineCredentialService.GetHiddenObjectsByCredentialID(credential.Key)
	if schemaData == nil {
		return nil, credentialOAuthApplicationOAuth2.ErrInvalidCredential
	}

	scope, _ := schemaData["scope"].(string)
	authURL, _ := schemaData["authUrl"].(string)
	authQueryParameters, _ := schemaData["authQueryParameters"].(string)

	challenge, err := pkce.GenerateChallenge()
	if err != nil {
		return nil, err
	}

	stateID, err := h.stateStore.Issue(ctx, &credentialOAuthDomainOAuth2.State{
		CredentialID: credential.ID,
		OwnerType:    command.OwnerType,
		OwnerID:      command.OwnerID,
		CodeVerifier: challenge.CodeVerifier,
	})
	if err != nil {
		return nil, err
	}

	authURL = h.buildAuthURL(credential.Data, scope, stateID, authURL, challenge.CodeChallenge, challenge.CodeChallengeMethod, authQueryParameters, h.oauth2RedirectURL)

	return &AuthURLResponse{URL: authURL}, nil
}

func (h *Handler) getCredentialData(ctx context.Context, ownerType commonDomain.OwnerType, ownerID uuid.UUID, credentialID uuid.UUID) (*Credential, error) {
	credentialInfo, err := h.credentialService.GetByIDForOwner(ctx, credentialID, ownerType, ownerID)
	if err != nil {
		return nil, err
	}

	var credentialData commonDomainOAuth2.Credential
	if err := json.ArgsToStruct(credentialInfo.Data, &credentialData); err != nil {
		return nil, err
	}

	if credentialInfo.SourceType == credentialsDomainCredentials.SourceTypePlatform {
		platformCred, ok := h.platformCredentialService.GetPlatformCredential(credentialInfo.Key)
		if !ok || platformCred == nil {
			return nil, credentialOAuthApplicationOAuth2.ErrInvalidCredential
		}
		if err := json.ArgsToStruct(platformCred.Data, &credentialData); err != nil {
			return nil, err
		}
		if !credentialData.HasClientCredentials() {
			return nil, credentialOAuthApplicationOAuth2.ErrPlatformCredentialMissingClientCreds
		}
	}

	return &Credential{
		ID:         credentialID,
		Key:        credentialInfo.Key,
		SourceType: credentialInfo.SourceType,
		Data:       credentialData,
	}, nil
}

func (h *Handler) buildAuthURL(cred commonDomainOAuth2.Credential, scope, state, authURL, codeChallenge, codeChallengeMethod, authQueryParameters, redirectURL string) string {
	v := url.Values{}

	if authQueryParameters != "" {
		if extra, err := url.ParseQuery(authQueryParameters); err == nil {
			for key, values := range extra {
				if len(values) > 0 {
					v.Set(key, values[0])
				}
			}
		}
	}

	if cred.ClientKey != "" {
		v.Set("client_key", cred.ClientKey)
	}
	if cred.ClientID != "" {
		v.Set("client_id", cred.ClientID)
	}
	if scope != "" {
		v.Set("scope", scope)
	}

	v.Set("code_challenge", codeChallenge)
	v.Set("code_challenge_method", codeChallengeMethod)
	v.Set("redirect_uri", redirectURL)
	v.Set("state", state)
	v.Set("response_type", "code")

	u, _ := url.Parse(authURL)
	u.RawQuery = v.Encode()
	return u.String()
}
