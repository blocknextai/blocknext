package gettoolinvocationbyid

import (
	"context"
	"log/slog"

	apiKeysApplicationAPIKeys "github.com/blocknextai/platform-api/internal/apikeys/application/apikeys"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/blocknextai/platform-api/internal/executions/domain/toolinvocations"
)

const deletedAPIKeyName = "[deleted]"

type Handler struct {
	toolInvocationRepository toolinvocations.ToolInvocationRepository
	apiKeyService            apiKeysApplicationAPIKeys.APIKeyService
}

func New(
	toolInvocationRepository toolinvocations.ToolInvocationRepository,
	apiKeyService apiKeysApplicationAPIKeys.APIKeyService,
) *Handler {
	return &Handler{
		toolInvocationRepository: toolInvocationRepository,
		apiKeyService:            apiKeyService,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetToolInvocationByIDQuery) (*GetToolInvocationByIDResponse, error) {
	toolInvocation, err := h.toolInvocationRepository.GetByIDAndOrganizationID(ctx, request.ID, request.OrganizationID)
	if err != nil {
		return nil, err
	}

	return MapToolInvocationToResponse(toolInvocation, h.resolveAPIKeyName(ctx, toolInvocation)), nil
}

func (h *Handler) resolveAPIKeyName(ctx context.Context, toolInvocation *toolinvocations.ToolInvocation) *string {
	if toolInvocation.APIKeyID == nil {
		return nil
	}

	apiKey, err := h.apiKeyService.GetByOwnerAndID(
		ctx,
		commonDomain.OwnerTypeOrganization,
		toolInvocation.OrganizationID,
		*toolInvocation.APIKeyID,
	)
	if err != nil {
		slog.WarnContext(ctx, "failed to resolve api key for tool invocation",
			"component", "gettoolinvocationbyid",
			"organization_id", toolInvocation.OrganizationID,
			"api_key_id", *toolInvocation.APIKeyID,
			"error", err)
		return new(deletedAPIKeyName)
	}

	return &apiKey.Name
}
