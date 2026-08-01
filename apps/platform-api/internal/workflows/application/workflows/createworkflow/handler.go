package createworkflow

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	organizationsApplicationOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/application/organizationusers"
	workflowsDomainWorkflows "github.com/blocknextai/platform-api/internal/workflows/domain/workflows"
)

type Handler struct {
	workflowRepository      workflowsDomainWorkflows.WorkflowRepository
	transactionManager      database.TransactionManager
	organizationUserService organizationsApplicationOrganizationUsers.OrganizationUserService
}

func New(
	workflowRepository workflowsDomainWorkflows.WorkflowRepository,
	transactionManager database.TransactionManager,
	organizationUserService organizationsApplicationOrganizationUsers.OrganizationUserService,
) *Handler {
	return &Handler{
		workflowRepository:      workflowRepository,
		transactionManager:      transactionManager,
		organizationUserService: organizationUserService,
	}
}

func (h *Handler) Handle(ctx context.Context, request *CreateWorkflowCommand) (*CreateWorkflowResponse, error) {
	var response *CreateWorkflowResponse

	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		user, err := h.organizationUserService.GetByOrganizationIDAndUserID(txCtx, request.OrganizationID, request.UserID)
		if err != nil {
			return err
		}

		workflow, err := workflowsDomainWorkflows.New(
			request.OrganizationID,
			user.ID,
			request.Title,
			&request.Description,
			false,
			true,
			request.Nodes,
			request.Edges,
		)
		if err != nil {
			return err
		}

		err = h.workflowRepository.Create(txCtx, workflow)
		if err != nil {
			return err
		}

		response = &CreateWorkflowResponse{
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
