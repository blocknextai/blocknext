package updateorganization

import (
	"github.com/google/uuid"
)

type UpdateOrganizationResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
}
