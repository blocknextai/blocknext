package createworkflow

import (
	"github.com/google/uuid"
)

type CreateWorkflowResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	IsPinned    bool      `json:"isPinned"`
}
