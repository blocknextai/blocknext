package credentialresolver

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/apperror"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonDomainCredential "github.com/blocknextai/platform-api/internal/common/domain/credential"
	credentialoauthContract "github.com/blocknextai/platform-api/internal/modules/credentialoauth/contract"
	credentialsContract "github.com/blocknextai/platform-api/internal/modules/credentials/contract"
)

var (
	ErrCredentialReferenceInvalid = apperror.Validation("invalid credential reference format")
	ErrCredentialScopeMismatch    = apperror.Forbidden("credential scope does not match the owner scope")
	ErrCredentialTypeMismatch     = apperror.Validation("credential type does not match the credential the tool requires")
)

type CredentialResolver interface {
	Resolve(
		ctx context.Context,
		ownerType commonDomain.OwnerType,
		ownerID uuid.UUID,
		references map[string]any,
	) (map[string]any, error)
}

type credentialResolver struct {
	credentialService      credentialsContract.CredentialService
	tokenRegenerateService credentialoauthContract.CredentialOAuthTokenRegenerateService
}

func NewCredentialResolver(
	credentialService credentialsContract.CredentialService,
	tokenRegenerateService credentialoauthContract.CredentialOAuthTokenRegenerateService,
) CredentialResolver {
	return &credentialResolver{
		credentialService:      credentialService,
		tokenRegenerateService: tokenRegenerateService,
	}
}

func (r *credentialResolver) Resolve(
	ctx context.Context,
	ownerType commonDomain.OwnerType,
	ownerID uuid.UUID,
	references map[string]any,
) (map[string]any, error) {
	resolved := make(map[string]any, len(references))

	for credentialKey, raw := range references {
		ref, ok := raw.(string)
		if !ok {
			return nil, ErrCredentialReferenceInvalid.WithCause(errors.New(credentialKey))
		}

		scope, credentialID, err := commonDomainCredential.ParseReference(ref)
		if err != nil {
			slog.WarnContext(ctx, "failed to parse credential reference",
				"component", "mcp_credential_resolver",
				"credential_key", credentialKey,
				"error", err,
			)
			return nil, ErrCredentialReferenceInvalid.WithCause(errors.New(credentialKey))
		}

		if !scopeMatchesOwner(scope, ownerType) {
			return nil, ErrCredentialScopeMismatch.WithCause(errors.New(credentialKey))
		}

		info, err := r.credentialService.GetByIDForOwner(ctx, credentialID, ownerType, ownerID)
		if err != nil {
			return nil, err
		}

		if info.Key != credentialKey {
			return nil, ErrCredentialTypeMismatch.WithCause(errors.New(credentialKey + " expected, got " + info.Key))
		}

		data, err := r.tokenRegenerateService.RegenerateTokenIfNeeded(ctx, ownerType, ownerID, credentialID)
		if err != nil {
			return nil, err
		}

		resolved[credentialKey] = data
	}

	return resolved, nil
}

func scopeMatchesOwner(scope commonDomainCredential.Scope, ownerType commonDomain.OwnerType) bool {
	return scope == commonDomainCredential.OrganizationCredentialScope &&
		ownerType == commonDomain.OwnerTypeOrganization
}
