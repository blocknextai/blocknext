package credentialresolver

import (
	"context"
	"errors"
	"log/slog"

	"github.com/blocknextai/go-packages/apperror"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonDomainCredential "github.com/blocknextai/platform-api/internal/common/domain/credential"
	credentialOAuthApplicationRegenerate "github.com/blocknextai/platform-api/internal/credentialoauth/application/regenerate"
	"github.com/google/uuid"
)

var (
	ErrCredentialReferenceInvalid = apperror.Validation("invalid credential reference format")
	ErrCredentialScopeMismatch    = apperror.Forbidden("credential scope does not match api key scope")
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
	tokenRegenerateService credentialOAuthApplicationRegenerate.CredentialOAuthTokenRegenerateService
}

func NewCredentialResolver(
	tokenRegenerateService credentialOAuthApplicationRegenerate.CredentialOAuthTokenRegenerateService,
) CredentialResolver {
	return &credentialResolver{
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

		data, err := r.tokenRegenerateService.RegenerateTokenIfNeeded(ctx, ownerType, ownerID, credentialID)
		if err != nil {
			return nil, err
		}

		resolved[credentialKey] = data
	}

	return resolved, nil
}

func scopeMatchesOwner(scope commonDomainCredential.Scope, ownerType commonDomain.OwnerType) bool {
	switch scope {
	case commonDomainCredential.OrganizationCredentialScope:
		return ownerType == commonDomain.OwnerTypeOrganization
	case commonDomainCredential.UserCredentialScope:
		return ownerType == commonDomain.OwnerTypeUser
	default:
		return false
	}
}
