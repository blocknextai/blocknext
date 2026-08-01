package nodeexecutions

import (
	"context"
	"time"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/executions/domain/nodeexecutions"
	"github.com/google/uuid"
)

type NodeExecutionService interface {
	Create(
		ctx context.Context,
		taskID uuid.UUID,
		nodeType string,
		nodeID string,
		status string,
		inputs []map[string]any,
		errorMessage *string,
		startedAt *time.Time,
	) (uuid.UUID, error)

	Update(
		ctx context.Context,
		taskID uuid.UUID,
		nodeExecutionID uuid.UUID,
		status string,
		inputs []map[string]any,
		outputs []map[string]any,
		functionCallingOutputs []map[string]any,
		errorMessage *string,
		startedAt *time.Time,
		completedAt *time.Time,
	) error

	GetAllByTaskID(ctx context.Context, taskID uuid.UUID) ([]*nodeexecutions.NodeExecution, error)
}

type nodeExecutionService struct {
	nodeExecutionRepository nodeexecutions.NodeExecutionRepository
	transactionManager      database.TransactionManager
}

func NewNodeExecutionService(
	nodeExecutionRepository nodeexecutions.NodeExecutionRepository,
	transactionManager database.TransactionManager,
) NodeExecutionService {
	return &nodeExecutionService{
		nodeExecutionRepository: nodeExecutionRepository,
		transactionManager:      transactionManager,
	}
}

func (s *nodeExecutionService) Create(
	ctx context.Context,
	taskID uuid.UUID,
	nodeType string,
	nodeID string,
	status string,
	inputs []map[string]any,
	errorMessage *string,
	startedAt *time.Time,
) (uuid.UUID, error) {
	var nodeExecutionID uuid.UUID

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		nodeExecution, err := nodeexecutions.New(
			taskID,
			nodeType,
			nodeID,
			status,
			inputs,
			errorMessage,
			startedAt,
		)
		if err != nil {
			return err
		}

		if err := s.nodeExecutionRepository.Create(txCtx, nodeExecution); err != nil {
			return err
		}

		nodeExecutionID = nodeExecution.ID
		return nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	return nodeExecutionID, nil
}

func (s *nodeExecutionService) Update(
	ctx context.Context,
	taskID uuid.UUID,
	nodeExecutionID uuid.UUID,
	status string,
	inputs []map[string]any,
	outputs []map[string]any,
	functionCallingOutputs []map[string]any,
	errorMessage *string,
	startedAt *time.Time,
	completedAt *time.Time,
) error {
	return s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		nodeExecution, err := s.nodeExecutionRepository.GetByIDAndTaskID(txCtx, nodeExecutionID, taskID)
		if err != nil {
			return err
		}

		if inputs == nil {
			inputs = nodeExecution.Inputs
		}
		if outputs == nil {
			outputs = nodeExecution.Outputs
		}
		if functionCallingOutputs == nil {
			functionCallingOutputs = nodeExecution.FunctionCallingOutputs
		}
		if startedAt == nil {
			startedAt = nodeExecution.StartedAt
		}
		if completedAt == nil {
			completedAt = nodeExecution.CompletedAt
		}

		nodeExecution, err = nodeExecution.Update(status, inputs, outputs, functionCallingOutputs, errorMessage, startedAt, completedAt)
		if err != nil {
			return err
		}

		if err := s.nodeExecutionRepository.Update(txCtx, nodeExecution); err != nil {
			return err
		}

		return nil
	})
}

func (s *nodeExecutionService) GetAllByTaskID(ctx context.Context, taskID uuid.UUID) ([]*nodeexecutions.NodeExecution, error) {
	return s.nodeExecutionRepository.GetAllByTaskID(ctx, taskID)
}
