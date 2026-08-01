package getallapikeys

import (
	"time"

	apiKeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/apikeys/domain/apikeys"
	"github.com/google/uuid"
)

type APIKeyResponse struct {
	ID         uuid.UUID                   `json:"id"`
	Name       string                      `json:"name"`
	Scopes     apiKeysDomainAPIKeys.Scopes `json:"scopes"`
	LastUsedAt *time.Time                  `json:"lastUsedAt,omitempty"`
	CreatedAt  time.Time                   `json:"createdAt"`
	UpdatedAt  time.Time                   `json:"updatedAt"`
}

type GetAllAPIKeysResponse struct {
	Items      []*APIKeyResponse
	TotalCount int64
}
