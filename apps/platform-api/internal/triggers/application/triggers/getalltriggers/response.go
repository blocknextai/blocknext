package getalltriggers

import (
	"time"

	"github.com/blocknextai/platform-api/internal/triggers/domain/triggers"
	"github.com/google/uuid"
)

type Workflow struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type TriggerResponse struct {
	ID          uuid.UUID            `json:"id"`
	Type        triggers.TriggerType `json:"type"`
	CronPattern *string              `json:"cronPattern"`
	Timezone    *string              `json:"timezone"`
	IsActive    bool                 `json:"isActive"`
	Workflow    Workflow             `json:"workflow"`
	CreatedAt   time.Time            `json:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt"`
}

type GetAllTriggersResponse struct {
	Items      []*TriggerResponse
	TotalCount int64
}
