package workflows

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/dag"
	workflowsApplicationWorkflows "github.com/blocknextai/platform-api/internal/modules/workflows/internal/application/workflows"
	workflowsDomainWorkflows "github.com/blocknextai/platform-api/internal/modules/workflows/internal/domain/workflows"
)

type UpdateWorkflowCommand struct {
	OrganizationID uuid.UUID
	WorkflowID     uuid.UUID
	Title          *string
	Description    *string
	IsPinned       *bool
	IsActive       *bool
	Nodes          []dag.Node
	Edges          []dag.Edge
}

const (
	UpdateWorkflowMaxTitleLength = 255
)

func (c *UpdateWorkflowCommand) Validate() error {
	if c.Title != nil {
		if strings.TrimSpace(*c.Title) == "" {
			return workflowsDomainWorkflows.ErrTitleIsRequired
		}

		if len(strings.TrimSpace(*c.Title)) > UpdateWorkflowMaxTitleLength {
			return workflowsDomainWorkflows.ErrTitleTooLong
		}
	}

	return nil
}

type UpdateWorkflowResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	IsPinned    bool      `json:"isPinned"`
	IsActive    bool      `json:"isActive"`
}

func (s *Service) UpdateWorkflow(ctx context.Context, request *UpdateWorkflowCommand) (*UpdateWorkflowResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	var response *UpdateWorkflowResponse

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		workflow, err := s.workflowRepository.GetByOrganizationIDAndID(txCtx, request.OrganizationID, request.WorkflowID)
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

		err = s.workflowRepository.Update(txCtx, workflow)
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
