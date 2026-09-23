package taskrunner

import (
	"context"

	"github.com/google/uuid"

	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/domain/taskrunner"
)

const (
	defaultMaxConcurrentTasks = 1
)

type concurrencyLimitResolver struct {
	maxConcurrentTasks int64
}

func NewConcurrencyLimitResolver(
	maxConcurrentTasks int64,
) taskRunnerDomainTaskRunner.ConcurrencyLimitResolver {
	return &concurrencyLimitResolver{
		maxConcurrentTasks: maxConcurrentTasks,
	}
}

func (r *concurrencyLimitResolver) GetMaxConcurrentTasks(_ context.Context, _ uuid.UUID) int64 {
	if r.maxConcurrentTasks <= 0 {
		return defaultMaxConcurrentTasks
	}
	return r.maxConcurrentTasks
}
