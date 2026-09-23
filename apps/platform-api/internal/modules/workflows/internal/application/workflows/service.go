package workflows

import (
	"context"

	"github.com/google/uuid"

	"github.com/blocknextai/platform-api/internal/modules/workflows/internal/domain/workflows"
)

type WorkflowService interface {
	GetWorkflow(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID) (*workflows.Workflow, error)
	GetAllByOrganizationID(ctx context.Context, organizationID uuid.UUID, searchQuery string, offset int, limit int) ([]*workflows.Workflow, int64, error)
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

func (s *workflowService) GetAllByOrganizationID(
	ctx context.Context,
	organizationID uuid.UUID,
	searchQuery string,
	offset int,
	limit int,
) ([]*workflows.Workflow, int64, error) {
	return s.workflowRepository.GetAllByOrganizationID(ctx, organizationID, searchQuery, offset, limit)
}

func (s *workflowService) GetCountByOrganizationID(ctx context.Context, organizationID uuid.UUID) (int64, error) {
	return s.workflowRepository.GetCountByOrganizationID(ctx, organizationID)
}
