package toolinvocations

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	executionsDomainToolInvocations "github.com/blocknextai/platform-api/internal/modules/executions/internal/domain/toolinvocations"
)

type GetToolInvocationByIDQuery struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
}

type GetToolInvocationByIDResponse struct {
	ID           uuid.UUID        `json:"id,omitempty"`
	Source       string           `json:"source,omitempty"`
	APIKeyID     *uuid.UUID       `json:"apiKeyId,omitempty"`
	APIKeyName   *string          `json:"apiKeyName,omitempty"`
	ToolID       string           `json:"toolId,omitempty"`
	Status       string           `json:"status,omitempty"`
	Parameters   map[string]any   `json:"parameters,omitempty"`
	Outputs      []map[string]any `json:"outputs,omitempty"`
	ErrorMessage *string          `json:"errorMessage,omitempty"`
	StartedAt    *time.Time       `json:"startedAt,omitempty"`
	CompletedAt  time.Time        `json:"completedAt"`
}

const deletedAPIKeyName = "[deleted]"

func (s *Service) GetToolInvocationByID(ctx context.Context, request *GetToolInvocationByIDQuery) (*GetToolInvocationByIDResponse, error) {
	toolInvocation, err := s.toolInvocationRepository.GetByIDAndOrganizationID(ctx, request.ID, request.OrganizationID)
	if err != nil {
		return nil, err
	}

	return MapToolInvocationToResponse(toolInvocation, s.resolveAPIKeyName(ctx, toolInvocation)), nil
}

func (s *Service) resolveAPIKeyName(ctx context.Context, toolInvocation *executionsDomainToolInvocations.ToolInvocation) *string {
	if toolInvocation.APIKeyID == nil {
		return nil
	}

	apiKey, err := s.apiKeyService.GetByOwnerAndID(
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

func MapToolInvocationToResponse(
	toolInvocation *executionsDomainToolInvocations.ToolInvocation,
	apiKeyName *string,
) *GetToolInvocationByIDResponse {
	return &GetToolInvocationByIDResponse{
		ID:           toolInvocation.ID,
		Source:       toolInvocation.Source.String(),
		APIKeyID:     toolInvocation.APIKeyID,
		APIKeyName:   apiKeyName,
		ToolID:       toolInvocation.ToolID,
		Status:       toolInvocation.Status.String(),
		Parameters:   toolInvocation.Parameters,
		Outputs:      toolInvocation.Outputs,
		ErrorMessage: toolInvocation.ErrorMessage,
		StartedAt:    new(toolInvocation.StartedAt),
		CompletedAt:  toolInvocation.CompletedAt,
	}
}
