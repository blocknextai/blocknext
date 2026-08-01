package deleteworkflow

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	workflowsApplicationWorkflows "github.com/blocknextai/platform-api/internal/workflows/application/workflows"
	workflowsDomainWorkflows "github.com/blocknextai/platform-api/internal/workflows/domain/workflows"
)

type Handler struct {
	workflowRepository workflowsDomainWorkflows.WorkflowRepository
	transactionManager database.TransactionManager
}

func New(
	workflowRepository workflowsDomainWorkflows.WorkflowRepository,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		workflowRepository: workflowRepository,
		transactionManager: transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, request *DeleteWorkflowCommand) (*DeleteWorkflowResponse, error) {
	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		workflow, err := h.workflowRepository.GetByOrganizationIDAndID(txCtx, request.OrganizationID, request.WorkflowID)
		if err != nil {
			return workflowsDomainWorkflows.ErrWorkflowNotFound
		}

		workflow, err = workflow.Delete()
		if err != nil {
			return err
		}

		err = h.workflowRepository.Delete(txCtx, workflow)
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
