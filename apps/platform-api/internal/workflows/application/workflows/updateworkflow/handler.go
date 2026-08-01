package updateworkflow

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

func (h *Handler) Handle(ctx context.Context, request *UpdateWorkflowCommand) (*UpdateWorkflowResponse, error) {
	var response *UpdateWorkflowResponse

	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		workflow, err := h.workflowRepository.GetByOrganizationIDAndID(txCtx, request.OrganizationID, request.WorkflowID)
		if err != nil {
			return err
		}

		if request.Title != nil && *request.Title != "" {
			workflow.Title = *request.Title
		}

		if request.Description != nil {
			workflow.Description = request.Description
		}

		if request.IsPinned != nil {
			workflow.IsPinned = *request.IsPinned
		}

		if request.IsActive != nil {
			workflow.IsActive = *request.IsActive
		}

		if len(request.Nodes) > 0 {
			workflow.Nodes = request.Nodes
		}

		if len(request.Edges) > 0 {
			workflow.Edges = request.Edges
		}

		err = h.workflowRepository.Update(txCtx, workflow)
		if err != nil {
			return workflowsApplicationWorkflows.ErrFailedToUpdateWorkflow.WithCause(err)
		}

		response = &UpdateWorkflowResponse{
			ID:          workflow.ID,
			Title:       workflow.Title,
			Description: workflow.Description,
			IsPinned:    workflow.IsPinned,
			IsActive:    workflow.IsActive,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
