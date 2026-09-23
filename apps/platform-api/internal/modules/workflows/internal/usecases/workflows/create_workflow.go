package workflows

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/dag"
	workflowsDomainWorkflows "github.com/blocknextai/platform-api/internal/modules/workflows/internal/domain/workflows"
)

type CreateWorkflowCommand struct {
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Title          string
	Description    string
	Nodes          []dag.Node
	Edges          []dag.Edge
}

const (
	MaxTitleLength = 255
)

func (c *CreateWorkflowCommand) Validate() error {
	if strings.TrimSpace(c.Title) == "" {
		return workflowsDomainWorkflows.ErrTitleIsRequired
	}

	if len(strings.TrimSpace(c.Title)) > MaxTitleLength {
		return workflowsDomainWorkflows.ErrTitleTooLong
	}

	return nil
}

type CreateWorkflowResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	IsPinned    bool      `json:"isPinned"`
}

func (s *Service) CreateWorkflow(ctx context.Context, request *CreateWorkflowCommand) (*CreateWorkflowResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	var response *CreateWorkflowResponse

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		user, err := s.organizationUserService.GetByOrganizationIDAndUserID(txCtx, request.OrganizationID, request.UserID)
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

		err = s.workflowRepository.Create(txCtx, workflow)
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
