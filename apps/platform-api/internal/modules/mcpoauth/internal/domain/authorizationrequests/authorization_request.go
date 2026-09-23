package authorizationrequests

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

type AuthorizationRequest struct {
	database.BaseEntity

	ClientID            string
	RedirectURI         string
	Scopes              mcpOAuthDomainOAuth2.Scopes
	State               *string
	CodeChallenge       string
	CodeChallengeMethod mcpOAuthDomainOAuth2.CodeChallengeMethod
	Resource            string
	Status              Status
	UserID              *uuid.UUID
	OrganizationID      *uuid.UUID
	ExpiresAt           time.Time
	ResolvedAt          *time.Time
}

func New(
	clientID string,
	redirectURI string,
	scopes mcpOAuthDomainOAuth2.Scopes,
	state *string,
	codeChallenge string,
	codeChallengeMethod mcpOAuthDomainOAuth2.CodeChallengeMethod,
	resource string,
	ttl time.Duration,
) (*AuthorizationRequest, error) {
	now := time.Now().UTC()

	request := &AuthorizationRequest{
		ID:                  bnuuid.NewV7(),
		CreatedAt:           now,
		UpdatedAt:           now,
		ClientID:            clientID,
		RedirectURI:         redirectURI,
		Scopes:              scopes,
		State:               state,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
		Resource:            resource,
		Status:              PendingStatus,
		ExpiresAt:           now.Add(ttl),
	}

	return request.validateThenReturn()
}

func (r *AuthorizationRequest) IsExpired() bool {
	return time.Now().UTC().After(r.ExpiresAt)
}

func (r *AuthorizationRequest) IsPending() bool {
	return r.Status == PendingStatus && !r.IsExpired()
}

func (r *AuthorizationRequest) Approve(
	userID uuid.UUID,
	organizationID uuid.UUID,
	scopes mcpOAuthDomainOAuth2.Scopes,
) (*AuthorizationRequest, error) {
	if err := r.ensureResolvable(); err != nil {
		return nil, err
	}

	if len(scopes) == 0 || !r.Scopes.HasAll(scopes) {
		return nil, ErrInvalidScopes
	}

	now := time.Now().UTC()

	r.Scopes = scopes
	r.Status = ApprovedStatus
	r.UserID = new(userID)
	r.OrganizationID = new(organizationID)
	r.ResolvedAt = new(now)
	r.UpdatedAt = now

	return r.validateThenReturn()
}

func (r *AuthorizationRequest) Deny(userID uuid.UUID) (*AuthorizationRequest, error) {
	if err := r.ensureResolvable(); err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	r.Status = DeniedStatus
	r.UserID = new(userID)
	r.ResolvedAt = new(now)
	r.UpdatedAt = now

	return r.validateThenReturn()
}

func (r *AuthorizationRequest) ensureResolvable() error {
	if r.Status != PendingStatus {
		return ErrAuthorizationRequestResolved
	}

	if r.IsExpired() {
		return ErrAuthorizationRequestExpired
	}

	return nil
}

func (r *AuthorizationRequest) validateThenReturn() (*AuthorizationRequest, error) {
	if strings.TrimSpace(r.ClientID) == "" {
		return nil, ErrInvalidClientID
	}

	if strings.TrimSpace(r.RedirectURI) == "" {
		return nil, ErrInvalidRedirectURI
	}

	if strings.TrimSpace(r.CodeChallenge) == "" {
		return nil, ErrInvalidCodeChallenge
	}

	if !r.CodeChallengeMethod.IsValid() {
		return nil, ErrInvalidCodeChallengeMethod
	}

	if strings.TrimSpace(r.Resource) == "" {
		return nil, ErrInvalidResource
	}

	if len(r.Scopes) == 0 || !r.Scopes.IsValid() {
		return nil, ErrInvalidScopes
	}

	if !r.Status.IsValid() {
		return nil, ErrInvalidStatus
	}

	if r.UserID != nil && *r.UserID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	if r.OrganizationID != nil && *r.OrganizationID == uuid.Nil {
		return nil, ErrInvalidOrganizationID
	}

	return r, nil
}
