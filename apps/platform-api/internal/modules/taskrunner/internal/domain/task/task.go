package task

import (
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/dag"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
	"github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/domain"
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
	TriggerContext      *nodeengineContract.TriggerContext
}
