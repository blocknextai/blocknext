package oauth2

import (
	"context"
	"net/url"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/go-packages/pkce"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonDomainOAuth2 "github.com/blocknextai/platform-api/internal/common/domain/oauth2"
	credentialOAuthApplicationOAuth2 "github.com/blocknextai/platform-api/internal/modules/credentialoauth/internal/application/oauth2"
	credentialOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/credentialoauth/internal/domain/oauth2"
	credentialsContract "github.com/blocknextai/platform-api/internal/modules/credentials/contract"
)

type AuthURLCommand struct {
	OwnerType    commonDomain.OwnerType
	OwnerID      uuid.UUID
	CredentialID uuid.UUID
}

var ErrCredentialIDIsRequired = apperror.Validation("credential id is required")

func (c *AuthURLCommand) Validate() error {
	if c.CredentialID == uuid.Nil {
		return ErrCredentialIDIsRequired
	}
	return nil
}

type AuthURLResponse struct {
	URL string `json:"url"`
}

type Credential struct {
	ID         uuid.UUID
	Key        string
	SourceType credentialsContract.SourceType
	Data       commonDomainOAuth2.Credential
}

func (s *Service) AuthURL(ctx context.Context, command *AuthURLCommand) (*AuthURLResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	credential, err := s.getCredentialData(ctx, command.OwnerType, command.OwnerID, command.CredentialID)
	if err != nil {
		return nil, err
	}

	schemaData := s.nodeEngineCredentialService.GetHiddenObjectsByCredentialID(credential.Key)
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

	stateID, err := s.stateStore.Issue(ctx, &credentialOAuthDomainOAuth2.State{
		CredentialID: credential.ID,
		OwnerType:    command.OwnerType,
		OwnerID:      command.OwnerID,
		CodeVerifier: challenge.CodeVerifier,
	})
	if err != nil {
		return nil, err
	}

	authURL = s.buildAuthURL(credential.Data, scope, stateID, authURL, challenge.CodeChallenge, challenge.CodeChallengeMethod, authQueryParameters, s.oauth2RedirectURL)

	return &AuthURLResponse{URL: authURL}, nil
}

func (s *Service) getCredentialData(ctx context.Context, ownerType commonDomain.OwnerType, ownerID uuid.UUID, credentialID uuid.UUID) (*Credential, error) {
	credentialInfo, err := s.credentialService.GetByIDForOwner(ctx, credentialID, ownerType, ownerID)
	if err != nil {
		return nil, err
	}

	var credentialData commonDomainOAuth2.Credential
	if err := json.ArgsToStruct(credentialInfo.Data, &credentialData); err != nil {
		return nil, err
	}

	if credentialInfo.SourceType == credentialsContract.SourceTypePlatform {
		platformCred, ok := s.platformCredentialService.GetPlatformCredential(credentialInfo.Key)
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

func (s *Service) buildAuthURL(cred commonDomainOAuth2.Credential, scope, state, authURL, codeChallenge, codeChallengeMethod, authQueryParameters, redirectURL string) string {
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
