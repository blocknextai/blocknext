package taskrunner

import (
	"context"

	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
	"github.com/google/uuid"
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
