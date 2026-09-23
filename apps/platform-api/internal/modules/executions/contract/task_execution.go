package contract

import (
	executionsApplicationTaskExecutions "github.com/blocknextai/platform-api/internal/modules/executions/internal/application/taskexecutions"
	executionsDomainTaskExecutions "github.com/blocknextai/platform-api/internal/modules/executions/internal/domain/taskexecutions"
)

type ExecutionType = executionsDomainTaskExecutions.ExecutionType

const ExecutionTypeAPI = executionsDomainTaskExecutions.ExecutionTypeAPI
const ExecutionTypeManual = executionsDomainTaskExecutions.ExecutionTypeManual
const ExecutionTypeSchedule = executionsDomainTaskExecutions.ExecutionTypeSchedule
const ExecutionTypeWebhook = executionsDomainTaskExecutions.ExecutionTypeWebhook

type TaskExecution = executionsDomainTaskExecutions.TaskExecution

type TaskExecutionService = executionsApplicationTaskExecutions.TaskExecutionService
