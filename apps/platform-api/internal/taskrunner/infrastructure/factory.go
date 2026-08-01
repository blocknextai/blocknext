package infrastructure

import (
	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/platform-api/internal/config"
	executionsApplicationTaskClaims "github.com/blocknextai/platform-api/internal/executions/application/taskclaims"
	taskRunnerApplicationDispatchers "github.com/blocknextai/platform-api/internal/taskrunner/application/dispatchers"
	taskRunnerApplicationTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/application/taskrunner"
	"github.com/blocknextai/platform-api/internal/taskrunner/domain/taskqueue"
	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
	taskRunnerInfrastructureLeaderRedis "github.com/blocknextai/platform-api/internal/taskrunner/infrastructure/leader/redis"
	taskRunnerInfrastructureSemaphoreRedis "github.com/blocknextai/platform-api/internal/taskrunner/infrastructure/semaphore/redis"
	taskRunnerInfrastructureTaskQueueRedis "github.com/blocknextai/platform-api/internal/taskrunner/infrastructure/taskqueue/redis"
)

var (
	ErrInvalidTaskQueueProvider     = apperror.Internal("invalid task queue provider")
	ErrInvalidTaskLockProvider      = apperror.Internal("invalid task lock provider")
	ErrInvalidTaskSemaphoreProvider = apperror.Internal("invalid task semaphore provider")
	ErrInvalidTaskRunnerMode        = apperror.Internal("invalid task runner mode")
)

func NewTaskQueue(queueOptions config.TaskRunnerQueueOptions, consumerName string) (taskqueue.TaskQueue, error) {
	if queueOptions.Provider == config.TaskQueueProviderRedis {
		return taskRunnerInfrastructureTaskQueueRedis.NewRedisTaskQueue(
			queueOptions.Redis.Address,
			queueOptions.Redis.Password,
			queueOptions.Redis.DB,
			taskRunnerInfrastructureTaskQueueRedis.PoolOptions{
				PoolSize:        queueOptions.Redis.PoolSize,
				MinIdleConns:    queueOptions.Redis.MinIdleConns,
				MaxIdleConns:    queueOptions.Redis.MaxIdleConns,
				PoolTimeout:     queueOptions.Redis.PoolTimeout,
				ConnMaxIdleTime: queueOptions.Redis.ConnMaxIdleTime,
				ConnMaxLifetime: queueOptions.Redis.ConnMaxLifetime,
			},
			taskRunnerInfrastructureTaskQueueRedis.Options{
				StreamName:    queueOptions.Redis.StreamName,
				ConsumerGroup: queueOptions.Redis.ConsumerGroup,
				ConsumerName:  consumerName,
				BlockTimeout:  queueOptions.Redis.BlockTimeout,
				PrefetchCount: queueOptions.Redis.PrefetchCount,
			},
		)
	}

	return nil, ErrInvalidTaskQueueProvider
}

func NewSemaphore(semaphoreOptions config.TaskRunnerSemaphoreOptions) (taskRunnerDomainTaskRunner.SemaphoreManager, error) {
	if semaphoreOptions.Provider == config.TaskSemaphoreProviderRedis {
		return taskRunnerInfrastructureSemaphoreRedis.NewRedisSemaphore(
			semaphoreOptions.Redis.Address,
			semaphoreOptions.Redis.Password,
			semaphoreOptions.Redis.DB,
			taskRunnerInfrastructureSemaphoreRedis.PoolOptions{
				PoolSize:        semaphoreOptions.Redis.PoolSize,
				MinIdleConns:    semaphoreOptions.Redis.MinIdleConns,
				MaxIdleConns:    semaphoreOptions.Redis.MaxIdleConns,
				PoolTimeout:     semaphoreOptions.Redis.PoolTimeout,
				ConnMaxIdleTime: semaphoreOptions.Redis.ConnMaxIdleTime,
				ConnMaxLifetime: semaphoreOptions.Redis.ConnMaxLifetime,
			},
			semaphoreOptions.TTL,
		)
	}

	return nil, ErrInvalidTaskSemaphoreProvider
}

func NewDispatcher(
	options config.TaskRunnerOptions,
	workerID string,
	workerPool taskRunnerDomainTaskRunner.WorkerPool,
	executor *taskRunnerApplicationTaskRunner.TaskExecutor,
	taskClaimService executionsApplicationTaskClaims.TaskClaimService,
) (taskRunnerDomainTaskRunner.TaskDispatcher, error) {
	if options.Mode == config.TaskRunnerModeEmbedded {
		return taskRunnerApplicationDispatchers.NewEmbedded(workerPool, executor), nil
	}

	if options.Mode == config.TaskRunnerModeQueue {
		queue, err := NewTaskQueue(options.Queue, workerID)
		if err != nil {
			return nil, err
		}
		return taskRunnerApplicationDispatchers.NewQueue(
			queue,
			workerPool,
			executor,
			taskClaimService,
			options.Queue.IdleTimeout,
			options.Queue.MaxRetries,
		), nil
	}

	return nil, ErrInvalidTaskRunnerMode
}

func NewLeaderRunner(leaderOptions config.TaskRunnerLeaderOptions, instanceID string) (taskRunnerDomainTaskRunner.LeaderRunner, error) {
	if leaderOptions.Provider == config.TaskLockProviderRedis {
		return taskRunnerInfrastructureLeaderRedis.NewRedisLeader(
			leaderOptions.Redis.Address,
			leaderOptions.Redis.Password,
			leaderOptions.Redis.DB,
			taskRunnerInfrastructureLeaderRedis.PoolOptions{
				PoolSize:        leaderOptions.Redis.PoolSize,
				MinIdleConns:    leaderOptions.Redis.MinIdleConns,
				MaxIdleConns:    leaderOptions.Redis.MaxIdleConns,
				PoolTimeout:     leaderOptions.Redis.PoolTimeout,
				ConnMaxIdleTime: leaderOptions.Redis.ConnMaxIdleTime,
				ConnMaxLifetime: leaderOptions.Redis.ConnMaxLifetime,
			},
			taskRunnerInfrastructureLeaderRedis.Options{
				Key:           leaderOptions.Key,
				InstanceID:    instanceID,
				TTL:           leaderOptions.TTL,
				PollInterval:  leaderOptions.PollInterval,
				RenewInterval: leaderOptions.RenewInterval,
			},
		)
	}

	return nil, ErrInvalidTaskLockProvider
}
