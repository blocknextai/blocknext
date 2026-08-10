package task

import (
	"sync"
	"time"

	"github.com/blocknextai/go-packages/dag"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	nodeEngineDomainAdapters "github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
	"github.com/blocknextai/platform-api/internal/taskrunner/domain"
	"github.com/google/uuid"
)

type Task struct {
	ID                uuid.UUID
	OrganizationID    uuid.UUID
	TriggeredByUserID *uuid.UUID
	ExecutionContext  commonDomain.ExecutionContext
	ContextItemID     uuid.UUID
	DAG               *dag.DAG

	Status             domain.Status
	WebhookToken       *string
	StartTime          *time.Time
	EndTime            *time.Time
	NodeExecutionIDMap map[string]uuid.UUID
	StartedNodes       sync.Map
	NodeStatuses       sync.Map

	PreviousNodeOutputs map[string][]map[string]any
	TriggerContext      *nodeEngineDomainAdapters.TriggerContext
}
