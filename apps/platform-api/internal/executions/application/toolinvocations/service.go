package toolinvocations

import (
	"context"
	"time"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/executions/domain/toolinvocations"
	"github.com/google/uuid"
)

type ToolInvocationService interface {
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
}

func NewToolInvocationService(
	toolInvocationRepository toolinvocations.ToolInvocationRepository,
	transactionManager database.TransactionManager,
) ToolInvocationService {
	return &toolInvocationService{
		toolInvocationRepository: toolInvocationRepository,
		transactionManager:       transactionManager,
	}
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
	var id uuid.UUID

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		toolInvocation, err := toolinvocations.New(
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

		if err := s.toolInvocationRepository.Create(txCtx, toolInvocation); err != nil {
			return err
		}

		id = toolInvocation.ID

		return nil
	})

	return id, err
}
