package getallorganizations

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	IsVerified  bool      `json:"isVerified"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type GetAllOrganizationsResponse = []OrganizationResponse
