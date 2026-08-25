package toolinvocations

import (
	"strings"
	"time"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	"github.com/google/uuid"
)

type ToolInvocation struct {
	database.BaseEntity

	OrganizationID uuid.UUID
	APIKeyID       *uuid.UUID
	Source         Source
	ToolID         string
	Status         Status
	Parameters     map[string]any
	Credentials    map[string]any
	Outputs        []map[string]any
	ErrorMessage   *string
	StartedAt      time.Time
	CompletedAt    time.Time
}

func New(
	organizationID uuid.UUID,
	apiKeyID *uuid.UUID,
	source Source,
	toolID string,
	status Status,
	parameters map[string]any,
	credentials map[string]any,
	outputs []map[string]any,
	errorMessage *string,
	startedAt time.Time,
	completedAt time.Time,
) (*ToolInvocation, error) {
	utcNow := time.Now().UTC()

	toolInvocation := &ToolInvocation{
		ID:             bnuuid.NewV7(),
		CreatedAt:      utcNow,
		UpdatedAt:      utcNow,
		DeletedAt:      nil,
		OrganizationID: organizationID,
		APIKeyID:       apiKeyID,
		Source:         source,
		ToolID:         toolID,
		Status:         status,
		Parameters:     parameters,
		Credentials:    credentials,
		Outputs:        outputs,
		ErrorMessage:   errorMessage,
		StartedAt:      startedAt,
		CompletedAt:    completedAt,
	}

	return toolInvocation.validateThenReturn()
}

func (t *ToolInvocation) validateThenReturn() (*ToolInvocation, error) {
	if t.OrganizationID == uuid.Nil {
		return nil, ErrOrganizationIDRequired
	}

	if !t.Source.IsValid() {
		return nil, ErrSourceInvalid
	}

	if strings.TrimSpace(t.ToolID) == "" {
		return nil, ErrToolIDIsRequired
	}

	if !t.Status.IsValid() {
		return nil, ErrStatusInvalid
	}

	return t, nil
}
