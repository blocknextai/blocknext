package gettoolinvocationbyid

import (
	"time"

	"github.com/google/uuid"
)

type GetToolInvocationByIDResponse struct {
	ID           uuid.UUID        `json:"id,omitempty"`
	Source       string           `json:"source,omitempty"`
	APIKeyID     *uuid.UUID       `json:"apiKeyId,omitempty"`
	APIKeyName   *string          `json:"apiKeyName,omitempty"`
	ToolID       string           `json:"toolId,omitempty"`
	Status       string           `json:"status,omitempty"`
	Parameters   map[string]any   `json:"parameters,omitempty"`
	Outputs      []map[string]any `json:"outputs,omitempty"`
	ErrorMessage *string          `json:"errorMessage,omitempty"`
	StartedAt    *time.Time       `json:"startedAt,omitempty"`
	CompletedAt  time.Time        `json:"completedAt"`
}
