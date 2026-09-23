package contract

import (
	executionsApplicationToolInvocations "github.com/blocknextai/platform-api/internal/modules/executions/internal/application/toolinvocations"
	executionsDomainToolInvocations "github.com/blocknextai/platform-api/internal/modules/executions/internal/domain/toolinvocations"
)

const SourceMCP = executionsDomainToolInvocations.SourceMCP

type Status = executionsDomainToolInvocations.Status

const StatusFailed = executionsDomainToolInvocations.StatusFailed
const StatusSuccess = executionsDomainToolInvocations.StatusSuccess

type ToolInvocationService = executionsApplicationToolInvocations.ToolInvocationService

type ToolInvocation = executionsDomainToolInvocations.ToolInvocation
