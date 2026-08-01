package duplicateworkflow

import (
	"github.com/google/uuid"
)

type DuplicateWorkflowResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	IsPinned    bool      `json:"isPinned"`
}
