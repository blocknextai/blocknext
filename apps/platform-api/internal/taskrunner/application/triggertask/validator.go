package triggertask

import (
	"github.com/blocknextai/go-packages/apperror"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	"github.com/google/uuid"
)

var (
	ErrInvalidExecutionContext = apperror.Validation("invalid execution context")
	ErrInvalidContextItemID    = apperror.Validation("invalid context item id")
	ErrInvalidTriggerType      = apperror.Validation("invalid trigger type")
	ErrInvalidNodes            = apperror.Validation("invalid nodes")
	ErrInvalidCronPattern      = apperror.Validation("invalid cron pattern")
)

func (c *TriggerTaskCommand) Validate() error {
	if !c.ExecutionContext.IsValid() {
		return ErrInvalidExecutionContext
	}

	if c.ContextItemID == uuid.Nil {
		return ErrInvalidContextItemID
	}

	if !c.TriggerType.IsValid() {
		return ErrInvalidTriggerType
	}

	if c.TriggerType == taskRunnerDomainTask.TaskTriggerTypeSchedule {
		if c.CronPattern == nil {
			return ErrInvalidCronPattern
		}
	}

	if len(c.Nodes) == 0 {
		return ErrInvalidNodes
	}

	return nil
}
