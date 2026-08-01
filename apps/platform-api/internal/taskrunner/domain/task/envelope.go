package task

import (
	nodeEngineDomainAdapters "github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
	"github.com/google/uuid"
)

type TaskEnvelope struct {
	TaskID              uuid.UUID                                `json:"taskId"`
	TriggerContext      *nodeEngineDomainAdapters.TriggerContext `json:"triggerContext,omitempty"`
	PreviousNodeOutputs map[string][]map[string]any              `json:"previousNodeOutputs,omitempty"`
}
