package triggertask

import (
	"github.com/blocknextai/go-packages/dag"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	"github.com/google/uuid"
)

type TriggerTaskCommand struct {
	TriggeredByUserID *uuid.UUID
	OrganizationID    uuid.UUID
	ExecutionContext  commonDomain.ExecutionContext
	ContextItemID     uuid.UUID
	TriggerType       taskRunnerDomainTask.TaskTriggerType
	CronPattern       *string

	RuntimePrompt string
	Nodes         []dag.Node
}
