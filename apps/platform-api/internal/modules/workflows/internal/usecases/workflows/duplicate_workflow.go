package workflows

import (
	"context"

	"github.com/google/uuid"

	"github.com/blocknextai/platform-api/internal/modules/workflows/internal/domain/workflows"
)

type DuplicateWorkflowCommand struct {
	WorkflowID     uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Title          string
	Description    *string
}

type DuplicateWorkflowResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	IsPinned    bool      `json:"isPinned"`
}

func (s *Service) DuplicateWorkflow(ctx context.Context, request *DuplicateWorkflowCommand) (*DuplicateWorkflowResponse, error) {
	var response *DuplicateWorkflowResponse

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		originalWorkflow, err := s.workflowRepository.GetByOrganizationIDAndID(txCtx, request.OrganizationID, request.WorkflowID)
		if err != nil {
			return err
		}

		user, err := s.organizationUserService.GetByOrganizationIDAndUserID(txCtx, request.OrganizationID, request.UserID)
		if err != nil {
			return err
		}

		workflow, err := workflows.New(
			request.OrganizationID,
			user.ID,
			request.Title,
			request.Description,
			false,
			originalWorkflow.IsActive,
			originalWorkflow.Nodes,
			originalWorkflow.Edges,
		)
		if err != nil {
			return err
		}

		err = s.workflowRepository.Create(txCtx, workflow)
		if err != nil {
			return err
		}

		response = &DuplicateWorkflowResponse{
			ID:          workflow.ID,
			Title:       workflow.Title,
			Description: workflow.Description,
			IsPinned:    workflow.IsPinned,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
