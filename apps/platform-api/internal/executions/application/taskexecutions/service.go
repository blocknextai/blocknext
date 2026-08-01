package taskexecutions

import (
	"context"
	"time"

	"github.com/blocknextai/go-packages/dag"
	"github.com/blocknextai/go-packages/database"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	taskclaimsApplication "github.com/blocknextai/platform-api/internal/executions/application/taskclaims"
	"github.com/blocknextai/platform-api/internal/executions/domain/taskexecutions"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizationusers"
	"github.com/google/uuid"
)

type TaskExecutionService interface {
	Create(
		ctx context.Context,
		id uuid.UUID,
		organizationID uuid.UUID,
		triggeredByUserID *uuid.UUID,
		flowTriggerID *uuid.UUID,
		executionContext commonDomain.ExecutionContext,
		contextItemID uuid.UUID,
		status string,
		executionType taskexecutions.ExecutionType,
		errorMessage *string,
		nodes []dag.Node,
		edges []dag.Edge,
		startedAt *time.Time,
	) error

	Update(
		ctx context.Context,
		id uuid.UUID,
		organizationID uuid.UUID,
		status string,
		errorMessage *string,
		completedAt *time.Time,
	) error

	GetByIDAndOrganizationID(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) (*taskexecutions.TaskExecution, error)
	GetByID(ctx context.Context, id uuid.UUID) (*taskexecutions.TaskExecution, error)
	GetAllByStatuses(ctx context.Context, statuses []string) ([]*taskexecutions.TaskExecution, error)
}

type taskExecutionService struct {
	taskExecutionRepository taskexecutions.TaskExecutionRepository
	taskClaimService        taskclaimsApplication.TaskClaimService
	organizationUserService organizationusers.OrganizationUserService
	transactionManager      database.TransactionManager
}

func NewTaskExecutionService(
	taskExecutionRepository taskexecutions.TaskExecutionRepository,
	taskClaimService taskclaimsApplication.TaskClaimService,
	organizationUserService organizationusers.OrganizationUserService,
	transactionManager database.TransactionManager,
) TaskExecutionService {
	return &taskExecutionService{
		taskExecutionRepository: taskExecutionRepository,
		taskClaimService:        taskClaimService,
		organizationUserService: organizationUserService,
		transactionManager:      transactionManager,
	}
}

func (s *taskExecutionService) Create(
	ctx context.Context,
	id uuid.UUID,
	organizationID uuid.UUID,
	triggeredByUserID *uuid.UUID,
	flowTriggerID *uuid.UUID,
	executionContext commonDomain.ExecutionContext,
	contextItemID uuid.UUID,
	status string,
	executionType taskexecutions.ExecutionType,
	errorMessage *string,
	nodes []dag.Node,
	edges []dag.Edge,
	startedAt *time.Time,
) error {
	return s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		var organizationUserID *uuid.UUID
		if triggeredByUserID != nil {
			organizationUser, err := s.organizationUserService.GetByOrganizationIDAndUserID(txCtx, organizationID, *triggeredByUserID)
			if err != nil {
				return err
			}
			organizationUserID = &organizationUser.ID
		}

		taskExecution, err := taskexecutions.New(
			id,
			organizationID,
			organizationUserID,
			flowTriggerID,
			executionContext,
			contextItemID,
			status,
			executionType,
			errorMessage,
			nodes,
			edges,
			startedAt,
		)
		if err != nil {
			return err
		}

		if err := s.taskExecutionRepository.Create(txCtx, taskExecution); err != nil {
			return err
		}

		if err := s.taskClaimService.Create(txCtx, taskExecution.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *taskExecutionService) Update(
	ctx context.Context,
	id uuid.UUID,
	organizationID uuid.UUID,
	status string,
	errorMessage *string,
	completedAt *time.Time,
) error {
	return s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		taskExecution, err := s.taskExecutionRepository.GetByIDAndOrganizationID(txCtx, id, organizationID)
		if err != nil {
			return err
		}

		taskExecution, err = taskExecution.Update(status, errorMessage, completedAt)
		if err != nil {
			return err
		}

		if err := s.taskExecutionRepository.Update(txCtx, taskExecution); err != nil {
			return err
		}

		return nil
	})
}

func (s *taskExecutionService) GetByIDAndOrganizationID(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) (*taskexecutions.TaskExecution, error) {
	return s.taskExecutionRepository.GetByIDAndOrganizationID(ctx, id, organizationID)
}

func (s *taskExecutionService) GetByID(ctx context.Context, id uuid.UUID) (*taskexecutions.TaskExecution, error) {
	return s.taskExecutionRepository.GetByID(ctx, id)
}

func (s *taskExecutionService) GetAllByStatuses(ctx context.Context, statuses []string) ([]*taskexecutions.TaskExecution, error) {
	return s.taskExecutionRepository.GetAllByStatuses(ctx, statuses)
}
