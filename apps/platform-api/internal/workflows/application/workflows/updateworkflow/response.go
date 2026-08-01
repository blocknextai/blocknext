package updateworkflow

import (
	"github.com/google/uuid"
)

type UpdateWorkflowResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	IsPinned    bool      `json:"isPinned"`
	IsActive    bool      `json:"isActive"`
}
