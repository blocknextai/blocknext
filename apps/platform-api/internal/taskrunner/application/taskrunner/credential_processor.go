package taskrunner

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonDomainCredential "github.com/blocknextai/platform-api/internal/common/domain/credential"
	credentialOAuthApplicationRegenerate "github.com/blocknextai/platform-api/internal/credentialoauth/application/regenerate"
	"github.com/google/uuid"
)

type CredentialProcessor struct {
	tokenRegenerateService credentialOAuthApplicationRegenerate.CredentialOAuthTokenRegenerateService
}

func NewCredentialProcessor(
	tokenRegenerateService credentialOAuthApplicationRegenerate.CredentialOAuthTokenRegenerateService,
) *CredentialProcessor {
	return &CredentialProcessor{
		tokenRegenerateService: tokenRegenerateService,
	}
}

func (p *CredentialProcessor) ProcessCredentials(
	ctx context.Context,
	organizationID uuid.UUID,
	triggeredByUserID *uuid.UUID,
	nodeType string,
	credentials map[string]any,
) (map[string]any, error) {
	processed := make(map[string]any, len(credentials))

	for k, v := range credentials {
		if value, ok := v.(string); ok && strings.HasPrefix(value, commonDomainCredential.CredentialPrefix) {
			resolved, err := p.processCredentialValue(ctx, organizationID, triggeredByUserID, nodeType, value)
			if err != nil {
				return nil, err
			}
			processed[k] = resolved
		} else {
			processed[k] = v
		}
	}

	return processed, nil
}

func (p *CredentialProcessor) processCredentialValue(
	ctx context.Context,
	organizationID uuid.UUID,
	triggeredByUserID *uuid.UUID,
	nodeType string,
	value string,
) (any, error) {
	scope, credentialID, err := commonDomainCredential.ParseReference(value)
	if err != nil {
		slog.WarnContext(ctx, "failed to parse credential reference",
			"component", "taskrunner_credential_processor",
			"nodeType", nodeType,
			"error", err,
		)
		return value, nil
	}

	var ownerType commonDomain.OwnerType
	var ownerID uuid.UUID
	switch scope {
	case commonDomainCredential.OrganizationCredentialScope:
		ownerType = commonDomain.OwnerTypeOrganization
		ownerID = organizationID
	case commonDomainCredential.UserCredentialScope:
		if triggeredByUserID == nil {
			return value, nil
		}
		ownerType = commonDomain.OwnerTypeUser
		ownerID = *triggeredByUserID
	default:
		return value, nil
	}

	updatedCredentialData, err := p.tokenRegenerateService.RegenerateTokenIfNeeded(ctx, ownerType, ownerID, credentialID)
	if err != nil {
		if errors.Is(err, credentialOAuthApplicationRegenerate.ErrRefreshTokenInvalid) {
			slog.WarnContext(ctx, "credential needs re-authentication",
				"credentialID", credentialID,
				"ownerType", ownerType,
				"nodeType", nodeType,
			)
		} else {
			slog.ErrorContext(ctx, "failed to regenerate oauth token",
				"credentialID", credentialID,
				"ownerType", ownerType,
				"nodeType", nodeType,
				"error", err,
			)
		}
		return nil, err
	}

	return updatedCredentialData, nil
}
