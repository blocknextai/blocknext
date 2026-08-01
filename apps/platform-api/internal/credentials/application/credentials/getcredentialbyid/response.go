package getcredentialbyid

import (
	"time"

	"github.com/blocknextai/platform-api/internal/credentials/domain/credentials"
	"github.com/google/uuid"
)

type GetCredentialByIDResponse struct {
	ID         uuid.UUID              `json:"id"`
	Title      string                 `json:"title"`
	Key        string                 `json:"key"`
	UIKey      string                 `json:"uiKey"`
	Data       map[string]any         `json:"data,omitempty"`
	SourceType credentials.SourceType `json:"sourceType"`
	CreatedAt  time.Time              `json:"createdAt"`
	UpdatedAt  time.Time              `json:"updatedAt"`
}
