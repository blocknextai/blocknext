package contract

import (
	workflowsApplicationWorkflows "github.com/blocknextai/platform-api/internal/modules/workflows/internal/application/workflows"
	workflowsDomainWorkflows "github.com/blocknextai/platform-api/internal/modules/workflows/internal/domain/workflows"
)

type WorkflowService = workflowsApplicationWorkflows.WorkflowService

type Workflow = workflowsDomainWorkflows.Workflow
