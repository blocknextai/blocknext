package events

import (
	"github.com/blocknextai/go-packages/json"
	executionsDomainToolInvocations "github.com/blocknextai/platform-api/internal/executions/domain/toolinvocations"
	taskRunnerDomainNode "github.com/blocknextai/platform-api/internal/taskrunner/domain/node"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
)

const SubscriberBuffer = 64

func MarshalTask(event *taskRunnerDomainTask.TaskEvent) (string, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return "", ErrFailedToMarshalTaskEvent
	}

	return string(payload), nil
}

func MarshalNode(event *taskRunnerDomainNode.NodeEvent) (string, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return "", ErrFailedToMarshalNodeEvent
	}

	return string(payload), nil
}

func MarshalToolInvocation(event *executionsDomainToolInvocations.ToolInvocationEvent) (string, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return "", ErrFailedToMarshalToolInvocationEvent
	}

	return string(payload), nil
}
