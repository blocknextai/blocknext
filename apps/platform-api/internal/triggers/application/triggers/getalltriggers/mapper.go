package getalltriggers

import (
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/triggers/domain/triggers"
	"github.com/google/uuid"
)

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
