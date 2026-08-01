package taskrunner

import (
	"context"

	"github.com/blocknextai/go-packages/dag"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	nodeEngineDomainAdapters "github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/triggers/domain/triggers"
	"github.com/google/uuid"
)

type TaskService interface {
	CancelTask(
		ctx context.Context,
		triggeredByUserID *uuid.UUID,
		organizationID uuid.UUID,
		taskID uuid.UUID,
	) error

	RerunAll(
		ctx context.Context,
		triggeredByUserID *uuid.UUID,
		organizationID uuid.UUID,
		taskID uuid.UUID,
		triggerType taskRunnerDomainTask.TaskTriggerType,
	) (*taskRunnerDomainTask.Task, error)

	RerunFailed(
		ctx context.Context,
		triggeredByUserID *uuid.UUID,
		organizationID uuid.UUID,
		taskID uuid.UUID,
		triggerType taskRunnerDomainTask.TaskTriggerType,
	) (*taskRunnerDomainTask.Task, error)

	ExecuteTask(
		ctx context.Context,
		triggeredByUserID *uuid.UUID,
		organizationID uuid.UUID,
		executionContext commonDomain.ExecutionContext,
		contextItemID uuid.UUID,
		triggerType taskRunnerDomainTask.TaskTriggerType,
		triggerContext *nodeEngineDomainAdapters.TriggerContext,
		cronPattern *string,
		runtimeConfig *triggersDomainTriggers.RuntimeConfig,
		nodes []dag.Node,
		edges []dag.Edge,
	) (*taskRunnerDomainTask.Task, error)
}
