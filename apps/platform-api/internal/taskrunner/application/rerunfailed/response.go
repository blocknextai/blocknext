package rerunfailed

import (
	"github.com/google/uuid"
)

type RerunFailedResponse struct {
	ID uuid.UUID `json:"id"`
}
