package contract

import (
	nodeengineApplicationExecutors "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/executors"
	nodeengineDomainExecutors "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/executors"
)

type BranchingExecutor = nodeengineDomainExecutors.BranchingExecutor
type ExecutorManager = nodeengineDomainExecutors.ExecutorManager

var GetExecutor = nodeengineDomainExecutors.GetExecutor

type ExecutorService = nodeengineApplicationExecutors.ExecutorService
