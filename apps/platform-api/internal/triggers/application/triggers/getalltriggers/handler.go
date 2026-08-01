package getalltriggers

import (
	"context"
	"log/slog"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/triggers/domain/triggers"
	workflowsApplicationWorkflows "github.com/blocknextai/platform-api/internal/workflows/application/workflows"
	"github.com/google/uuid"
)

type Handler struct {
	triggerRepository triggersDomainTriggers.TriggerRepository
	workflowService   workflowsApplicationWorkflows.WorkflowService
}

func New(
	triggerRepository triggersDomainTriggers.TriggerRepository,
	workflowService workflowsApplicationWorkflows.WorkflowService,
) *Handler {
	return &Handler{
		triggerRepository: triggerRepository,
		workflowService:   workflowService,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetAllTriggersQuery) (*GetAllTriggersResponse, error) {
	triggers, totalCount, err := h.triggerRepository.GetAllByOrganizationID(
		ctx,
		request.OrganizationID,
		request.Search.Query,
		request.Pagination.Offset,
		request.Pagination.Limit,
	)
	if err != nil {
		return nil, err
	}

	workflowsByID := make(map[uuid.UUID]Workflow)
	for _, trigger := range triggers {
		workflow := h.resolveWorkflow(ctx, trigger.ExecutionContext, trigger.ContextItemID, trigger.OrganizationID)
		workflowsByID[trigger.ID] = workflow
	}

	items := MapGetAllTriggersQueryToGetAllTriggersResponse(triggers, workflowsByID)

	return &GetAllTriggersResponse{
		Items:      items,
		TotalCount: totalCount,
	}, nil
}

func (h *Handler) resolveWorkflow(
	ctx context.Context,
	executionContext commonDomain.ExecutionContext,
	contextItemID uuid.UUID,
	organizationID uuid.UUID,
) Workflow {
	switch executionContext {
	case commonDomain.ExecutionContextWorkflow:
		workflow, err := h.workflowService.GetWorkflow(ctx, organizationID, contextItemID)
		if err != nil {
			slog.WarnContext(ctx, "Failed to resolve workflow info for trigger",
				"component", "getalltriggers.Handler",
				"organization_id", organizationID,
				"context_item_id", contextItemID,
				"error", err)
			return Workflow{Title: "[deleted]"}
		}
		return Workflow{
			ID:    workflow.ID,
			Title: workflow.Title,
		}
	default:
		return Workflow{Title: "[deleted]"}
	}
}
