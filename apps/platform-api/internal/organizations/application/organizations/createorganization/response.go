package createorganization

import (
	"github.com/google/uuid"
)

type CreateOrganizationResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
}
