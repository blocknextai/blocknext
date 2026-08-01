package duplicateworkflow

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	organizationsApplicationOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/application/organizationusers"
	"github.com/blocknextai/platform-api/internal/workflows/domain/workflows"
)

type Handler struct {
	workflowRepository      workflows.WorkflowRepository
	transactionManager      database.TransactionManager
	organizationUserService organizationsApplicationOrganizationUsers.OrganizationUserService
}

func New(
	workflowRepository workflows.WorkflowRepository,
	transactionManager database.TransactionManager,
	organizationUserService organizationsApplicationOrganizationUsers.OrganizationUserService,
) *Handler {
	return &Handler{
		workflowRepository:      workflowRepository,
		transactionManager:      transactionManager,
		organizationUserService: organizationUserService,
	}
}

func (h *Handler) Handle(ctx context.Context, request *DuplicateWorkflowCommand) (*DuplicateWorkflowResponse, error) {
	var response *DuplicateWorkflowResponse

	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		originalWorkflow, err := h.workflowRepository.GetByOrganizationIDAndID(txCtx, request.OrganizationID, request.WorkflowID)
		if err != nil {
			return err
		}

		user, err := h.organizationUserService.GetByOrganizationIDAndUserID(txCtx, request.OrganizationID, request.UserID)
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

		err = h.workflowRepository.Create(txCtx, workflow)
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
