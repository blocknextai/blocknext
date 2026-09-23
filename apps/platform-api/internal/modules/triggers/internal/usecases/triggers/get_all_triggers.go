package triggers

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/modules/triggers/internal/domain/triggers"
)

type GetAllTriggersQuery struct {
	OrganizationID uuid.UUID
	Search         resultPkg.SearchRequest
	Pagination     resultPkg.PaginationRequest
}

type Workflow struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type TriggerResponse struct {
	ID               uuid.UUID                          `json:"id"`
	Type             triggersDomainTriggers.TriggerType `json:"type"`
	CronPattern      *string                            `json:"cronPattern"`
	Timezone         *string                            `json:"timezone"`
	HasWebhookSecret bool                               `json:"hasWebhookSecret"`
	IsActive         bool                               `json:"isActive"`
	Workflow         Workflow                           `json:"workflow"`
	CreatedAt        time.Time                          `json:"createdAt"`
	UpdatedAt        time.Time                          `json:"updatedAt"`
}

type GetAllTriggersResponse struct {
	Items      []*TriggerResponse
	TotalCount int64
}

func (s *Service) GetAllTriggers(ctx context.Context, request *GetAllTriggersQuery) (*GetAllTriggersResponse, error) {
	triggers, totalCount, err := s.triggerRepository.GetAllByOrganizationID(
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
		workflow := s.resolveWorkflow(ctx, trigger.ExecutionContext, trigger.ContextItemID, trigger.OrganizationID)
		workflowsByID[trigger.ID] = workflow
	}

	items := MapGetAllTriggersQueryToGetAllTriggersResponse(triggers, workflowsByID)

	return &GetAllTriggersResponse{
		Items:      items,
		TotalCount: totalCount,
	}, nil
}

func (s *Service) resolveWorkflow(
	ctx context.Context,
	executionContext commonDomain.ExecutionContext,
	contextItemID uuid.UUID,
	organizationID uuid.UUID,
) Workflow {
	switch executionContext {
	case commonDomain.ExecutionContextWorkflow:
		workflow, err := s.workflowService.GetWorkflow(ctx, organizationID, contextItemID)
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

func MapGetAllTriggersQueryToGetAllTriggersResponse(
	triggers []*triggersDomainTriggers.Trigger,
	workflowsByID map[uuid.UUID]Workflow,
) []*TriggerResponse {
	response := make([]*TriggerResponse, 0, len(triggers))
	for _, trigger := range triggers {
		workflow := workflowsByID[trigger.ID]

		response = append(response, &TriggerResponse{
			ID:               trigger.ID,
			Type:             trigger.Type,
			CronPattern:      trigger.CronPattern,
			Timezone:         trigger.Timezone,
			HasWebhookSecret: trigger.WebhookSecret != nil,
			IsActive:         trigger.IsActive,
			Workflow:         workflow,
			CreatedAt:        trigger.CreatedAt,
			UpdatedAt:        trigger.UpdatedAt,
		})
	}
	return response
}
