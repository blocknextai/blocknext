package workflows

import (
	"context"

	"github.com/google/uuid"

	workflowsApplicationWorkflows "github.com/blocknextai/platform-api/internal/modules/workflows/internal/application/workflows"
	workflowsDomainWorkflows "github.com/blocknextai/platform-api/internal/modules/workflows/internal/domain/workflows"
)

type DeleteWorkflowCommand struct {
	OrganizationID uuid.UUID
	WorkflowID     uuid.UUID
}

type DeleteWorkflowResponse struct{}

func (s *Service) DeleteWorkflow(ctx context.Context, request *DeleteWorkflowCommand) (*DeleteWorkflowResponse, error) {
	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		workflow, err := s.workflowRepository.GetByOrganizationIDAndID(txCtx, request.OrganizationID, request.WorkflowID)
		if err != nil {
			return workflowsDomainWorkflows.ErrWorkflowNotFound
		}

		workflow, err = workflow.Delete()
		if err != nil {
			return err
		}

		err = s.workflowRepository.Delete(txCtx, workflow)
		if err != nil {
			return workflowsApplicationWorkflows.ErrFailedToDeleteWorkflow.WithCause(err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &DeleteWorkflowResponse{}, nil
}
