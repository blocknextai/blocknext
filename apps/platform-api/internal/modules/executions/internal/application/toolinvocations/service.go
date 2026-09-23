package toolinvocations

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/modules/executions/internal/domain/toolinvocations"
	"github.com/blocknextai/platform-api/internal/realtime"
	"github.com/blocknextai/platform-api/internal/realtime/events"
)

type ToolInvocationService interface {
	GetAllByOrganizationID(
		ctx context.Context,
		organizationID uuid.UUID,
		searchQuery string,
		offset int,
		limit int,
	) ([]*toolinvocations.ToolInvocation, int64, error)
	Record(
		ctx context.Context,
		organizationID uuid.UUID,
		apiKeyID *uuid.UUID,
		source toolinvocations.Source,
		toolID string,
		status toolinvocations.Status,
		parameters map[string]any,
		credentials map[string]any,
		outputs []map[string]any,
		errorMessage *string,
		startedAt time.Time,
		completedAt time.Time,
	) (uuid.UUID, error)
}

type toolInvocationService struct {
	toolInvocationRepository toolinvocations.ToolInvocationRepository
	transactionManager       database.TransactionManager
	broadcaster              realtime.Broadcaster
}

func NewToolInvocationService(
	toolInvocationRepository toolinvocations.ToolInvocationRepository,
	transactionManager database.TransactionManager,
	broadcaster realtime.Broadcaster,
) ToolInvocationService {
	return &toolInvocationService{
		toolInvocationRepository: toolInvocationRepository,
		transactionManager:       transactionManager,
		broadcaster:              broadcaster,
	}
}

func (s *toolInvocationService) GetAllByOrganizationID(
	ctx context.Context,
	organizationID uuid.UUID,
	searchQuery string,
	offset int,
	limit int,
) ([]*toolinvocations.ToolInvocation, int64, error) {
	return s.toolInvocationRepository.GetAllByOrganizationID(ctx, organizationID, searchQuery, offset, limit)
}

func (s *toolInvocationService) Record(
	ctx context.Context,
	organizationID uuid.UUID,
	apiKeyID *uuid.UUID,
	source toolinvocations.Source,
	toolID string,
	status toolinvocations.Status,
	parameters map[string]any,
	credentials map[string]any,
	outputs []map[string]any,
	errorMessage *string,
	startedAt time.Time,
	completedAt time.Time,
) (uuid.UUID, error) {
	var toolInvocation *toolinvocations.ToolInvocation

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		created, err := toolinvocations.New(
			organizationID,
			apiKeyID,
			source,
			toolID,
			status,
			parameters,
			credentials,
			outputs,
			errorMessage,
			startedAt,
			completedAt,
		)
		if err != nil {
			return err
		}

		if err := s.toolInvocationRepository.Create(txCtx, created); err != nil {
			return err
		}

		toolInvocation = created

		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	s.publish(ctx, toolInvocation)

	return toolInvocation.ID, nil
}

func (s *toolInvocationService) publish(ctx context.Context, toolInvocation *toolinvocations.ToolInvocation) {
	errorText := ""
	if toolInvocation.ErrorMessage != nil {
		errorText = *toolInvocation.ErrorMessage
	}

	event := &events.ToolInvocationEvent{
		ID:             toolInvocation.ID,
		Type:           events.ToolInvocationEventType,
		OrganizationID: toolInvocation.OrganizationID,
		Source:         string(toolInvocation.Source),
		ToolID:         toolInvocation.ToolID,
		Status:         string(toolInvocation.Status),
		Error:          errorText,
	}

	if err := s.broadcaster.PublishToolInvocationEvent(ctx, event); err != nil {
		slog.WarnContext(ctx, "failed to publish tool invocation event",
			"component", "executions",
			"organization_id", toolInvocation.OrganizationID,
			"tool_id", toolInvocation.ToolID,
			"error", err)
	}
}
