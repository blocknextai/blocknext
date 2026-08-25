package getalltoolinvocations

import (
	"time"

	"github.com/google/uuid"
)

type ToolInvocationResponse struct {
	ID           uuid.UUID  `json:"id,omitempty"`
	Source       string     `json:"source,omitempty"`
	ToolID       string     `json:"toolId,omitempty"`
	Status       string     `json:"status,omitempty"`
	ErrorMessage *string    `json:"errorMessage,omitempty"`
	StartedAt    *time.Time `json:"startedAt,omitempty"`
	CompletedAt  time.Time  `json:"completedAt"`
}

type GetAllToolInvocationsResponse struct {
	Items      []*ToolInvocationResponse `json:"items"`
	TotalCount int64                     `json:"totalCount"`
}
