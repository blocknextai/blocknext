package task

import (
	"github.com/google/uuid"

	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
)

type TaskEnvelope struct {
	TaskID              uuid.UUID                          `json:"taskId"`
	TriggerContext      *nodeengineContract.TriggerContext `json:"triggerContext,omitempty"`
	PreviousNodeOutputs map[string][]map[string]any        `json:"previousNodeOutputs,omitempty"`
}
