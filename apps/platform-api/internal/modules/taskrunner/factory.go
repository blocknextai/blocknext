package taskrunner

import (
	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/redisclient"
	"github.com/blocknextai/platform-api/internal/config"
	executionsContract "github.com/blocknextai/platform-api/internal/modules/executions/contract"
	taskRunnerApplicationDispatchers "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/application/dispatchers"
	taskRunnerApplicationTaskRunner "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/application/taskrunner"
	"github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/domain/taskqueue"
	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/domain/taskrunner"
	taskRunnerLeaderMemory "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/leader/memory"
	taskRunnerLeaderRedis "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/leader/redis"
	taskRunnerTaskQueueRedis "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/taskqueue/redis"
)

var (
	ErrInvalidTaskQueueProvider = apperror.Internal("invalid task queue provider")
	ErrInvalidTaskLockProvider  = apperror.Internal("invalid task lock provider")
)

func NewTaskQueue(queueOptions config.TaskRunnerQueueOptions, consumerName string) (taskqueue.TaskQueue, error) {
	if queueOptions.Provider == config.TaskQueueProviderRedis {
		return taskRunnerTaskQueueRedis.NewRedisTaskQueue(
			queueOptions.Redis.Address,
			queueOptions.Redis.Password,
			queueOptions.Redis.DB,
			redisclient.PoolOptions{
				PoolSize:        queueOptions.Redis.PoolSize,
				MinIdleConns:    queueOptions.Redis.MinIdleConns,
				MaxIdleConns:    queueOptions.Redis.MaxIdleConns,
				PoolTimeout:     queueOptions.Redis.PoolTimeout,
				ConnMaxIdleTime: queueOptions.Redis.ConnMaxIdleTime,
				ConnMaxLifetime: queueOptions.Redis.ConnMaxLifetime,
			},
			taskRunnerTaskQueueRedis.Options{
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

func NewDispatcher(
	options config.TaskRunnerOptions,
	workerID string,
	workerPool taskRunnerDomainTaskRunner.WorkerPool,
	executor *taskRunnerApplicationTaskRunner.TaskExecutor,
	taskClaimService executionsContract.TaskClaimService,
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
	if leaderOptions.Provider == config.TaskLockProviderMemory {
		return taskRunnerLeaderMemory.New(), nil
	}

	if leaderOptions.Provider == config.TaskLockProviderRedis {
		return taskRunnerLeaderRedis.NewRedisLeader(
			leaderOptions.Redis.Address,
			leaderOptions.Redis.Password,
			leaderOptions.Redis.DB,
			redisclient.PoolOptions{
				PoolSize:        leaderOptions.Redis.PoolSize,
				MinIdleConns:    leaderOptions.Redis.MinIdleConns,
				MaxIdleConns:    leaderOptions.Redis.MaxIdleConns,
				PoolTimeout:     leaderOptions.Redis.PoolTimeout,
				ConnMaxIdleTime: leaderOptions.Redis.ConnMaxIdleTime,
				ConnMaxLifetime: leaderOptions.Redis.ConnMaxLifetime,
			},
			taskRunnerLeaderRedis.Options{
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
