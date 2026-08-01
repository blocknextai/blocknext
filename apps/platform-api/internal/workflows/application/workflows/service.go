package workflows

import (
	"context"

	"github.com/blocknextai/platform-api/internal/workflows/domain/workflows"
	"github.com/google/uuid"
)

type WorkflowService interface {
	GetWorkflow(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID) (*workflows.Workflow, error)
	GetCountByOrganizationID(ctx context.Context, organizationID uuid.UUID) (int64, error)
}

type workflowService struct {
	workflowRepository workflows.WorkflowRepository
}

func NewWorkflowService(workflowRepository workflows.WorkflowRepository) WorkflowService {
	return &workflowService{
		workflowRepository: workflowRepository,
	}
}

func (s *workflowService) GetWorkflow(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID) (*workflows.Workflow, error) {
	workflow, err := s.workflowRepository.GetByOrganizationIDAndID(ctx, organizationID, workflowID)
	if err != nil {
		return nil, err
	}

	return workflow, nil
}

func (s *workflowService) GetCountByOrganizationID(ctx context.Context, organizationID uuid.UUID) (int64, error) {
	return s.workflowRepository.GetCountByOrganizationID(ctx, organizationID)
}
