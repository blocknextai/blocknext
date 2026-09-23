package taskrunner

import (
	"context"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/dag"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/domain/task"
	triggersContract "github.com/blocknextai/platform-api/internal/modules/triggers/contract"
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
		triggerContext *nodeengineContract.TriggerContext,
		cronPattern *string,
		runtimeConfig *triggersContract.RuntimeConfig,
		nodes []dag.Node,
		edges []dag.Edge,
	) (*taskRunnerDomainTask.Task, error)
}
