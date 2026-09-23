package taskrunner

import (
	"context"
	"log/slog"
	"os"
	"strings"

	apikeysContract "github.com/blocknextai/platform-api/internal/modules/apikeys/contract"

	commonDomainSemaphore "github.com/blocknextai/platform-api/internal/common/domain/semaphore"
	commonSemaphore "github.com/blocknextai/platform-api/internal/common/semaphore"

	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/cache"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	"github.com/blocknextai/platform-api/internal/common/auth"
	"github.com/blocknextai/platform-api/internal/config"
	credentialoauthContract "github.com/blocknextai/platform-api/internal/modules/credentialoauth/contract"
	executionsContract "github.com/blocknextai/platform-api/internal/modules/executions/contract"
	llmContract "github.com/blocknextai/platform-api/internal/modules/llm/contract"
	taskRunnerApplicationContextResolver "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/application/contextresolver"
	taskRunnerApplicationTaskRunner "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/application/taskrunner"
	taskRunnerApplicationWebhooks "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/application/webhooks"
	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/domain/taskrunner"
	taskRunnerHTTP "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/http"
	taskRunnerUseCases "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/usecases"
	triggersContract "github.com/blocknextai/platform-api/internal/modules/triggers/contract"
	workflowsContract "github.com/blocknextai/platform-api/internal/modules/workflows/contract"
	"github.com/blocknextai/platform-api/internal/realtime"
)

var (
	ErrInvalidTaskRunnerMode = apperror.Internal("invalid task runner mode")
)

type Dependencies struct {
	CacheService cache.Service
	Broadcaster  realtime.Broadcaster

	TaskRunnerOptions config.TaskRunnerOptions
	SemaphoreOptions  config.SemaphoreOptions

	FunctionCallingService                llmContract.FunctionCallingService
	CredentialOAuthTokenRegenerateService credentialoauthContract.CredentialOAuthTokenRegenerateService
	WorkflowService                       workflowsContract.WorkflowService
	TaskExecutionService                  executionsContract.TaskExecutionService
	TaskClaimService                      executionsContract.TaskClaimService
	NodeExecutionService                  executionsContract.NodeExecutionService
	TriggerService                        triggersContract.TriggerService
	WebhookResolver                       triggersContract.WebhookResolver
}

type Module struct {
	TaskService      taskRunnerDomainTaskRunner.TaskService
	ContextResolver  taskRunnerApplicationContextResolver.ContextResolver
	WebhookProcessor taskRunnerApplicationWebhooks.WebhookProcessor

	mode                config.TaskRunnerMode
	dispatcher          taskRunnerDomainTaskRunner.TaskDispatcher
	workerPool          taskRunnerDomainTaskRunner.WorkerPool
	cronService         taskRunnerDomainTaskRunner.CronService
	taskRecoveryService *taskRunnerApplicationTaskRunner.TaskRecoveryService
	semaphoreManager    commonDomainSemaphore.SemaphoreManager
	leaderRunner        taskRunnerDomainTaskRunner.LeaderRunner

	useCases *taskRunnerUseCases.Services
}

func NewModule(deps Dependencies) (*Module, error) {
	outputStore := taskRunnerApplicationTaskRunner.NewOutputStore()
	dataProcessor := taskRunnerApplicationTaskRunner.NewDataProcessor(outputStore)
	credentialProcessor := taskRunnerApplicationTaskRunner.NewCredentialProcessor(deps.CredentialOAuthTokenRegenerateService)
	eventPublisher := taskRunnerApplicationTaskRunner.NewEventPublisher(deps.Broadcaster)
	nodeExecutor := taskRunnerApplicationTaskRunner.NewNodeExecutor(
		deps.FunctionCallingService,
		dataProcessor,
		credentialProcessor,
		eventPublisher,
		deps.NodeExecutionService,
		outputStore,
	)
	executionCoordinator := taskRunnerApplicationTaskRunner.NewTaskExecutionCoordinator(
		deps.TaskExecutionService,
		deps.NodeExecutionService,
		deps.CredentialOAuthTokenRegenerateService,
		nodeExecutor,
		eventPublisher,
	)
	lifecycleManager := taskRunnerApplicationTaskRunner.NewTaskLifecycleManager(
		deps.TaskExecutionService,
		deps.NodeExecutionService,
		eventPublisher,
	)
	workerPool := taskRunnerApplicationTaskRunner.NewWorkerPool(
		deps.TaskRunnerOptions.Worker.PoolSize,
		deps.TaskRunnerOptions.Worker.PoolQueueSize,
	)
	contextManager := taskRunnerApplicationTaskRunner.NewContextManager()
	cronService := taskRunnerApplicationTaskRunner.NewCronService()
	contextResolver := taskRunnerApplicationContextResolver.NewContextResolver(
		deps.WorkflowService,
	)
	concurrencyLimitResolver := taskRunnerApplicationTaskRunner.NewConcurrencyLimitResolver(
		deps.SemaphoreOptions.MaxConcurrentExecutions,
	)

	workerID := resolveWorkerID(deps.TaskRunnerOptions.Queue.Redis.ConsumerName)

	semaphoreManager, err := commonSemaphore.New(deps.SemaphoreOptions)
	if err != nil {
		return nil, err
	}

	executor := taskRunnerApplicationTaskRunner.NewTaskExecutor(
		deps.TaskExecutionService,
		deps.TaskClaimService,
		deps.NodeExecutionService,
		executionCoordinator,
		lifecycleManager,
		semaphoreManager,
		contextManager,
		concurrencyLimitResolver,
		workerID,
		deps.TaskRunnerOptions.HeartbeatInterval,
		deps.TaskRunnerOptions.MaxExecutionTime,
	)

	dispatcher, err := NewDispatcher(
		deps.TaskRunnerOptions,
		workerID,
		workerPool,
		executor,
		deps.TaskClaimService,
	)
	if err != nil {
		return nil, err
	}

	taskService := taskRunnerApplicationTaskRunner.NewTaskService(
		contextResolver,
		executionCoordinator,
		dispatcher,
		contextManager,
		deps.TaskExecutionService,
		deps.NodeExecutionService,
		deps.TriggerService,
	)
	webhookProcessor := taskRunnerApplicationWebhooks.NewWebhookProcessor(deps.WebhookResolver, taskService, contextResolver)
	taskRecoveryService := taskRunnerApplicationTaskRunner.NewTaskRecoveryService(
		deps.TaskExecutionService,
		deps.TaskClaimService,
		deps.NodeExecutionService,
		dispatcher,
		deps.TriggerService,
		cronService,
		taskService,
		contextResolver,
		deps.TaskRunnerOptions.Queue.StaleClaimTimeout,
		deps.TaskRunnerOptions.RecoveryInterval,
	)

	useCases := taskRunnerUseCases.NewServices(taskRunnerUseCases.ServiceDependencies{
		TaskService:     taskService,
		ContextResolver: contextResolver,
	})

	leaderRunner, err := NewLeaderRunner(
		deps.TaskRunnerOptions.Leader,
		workerID,
	)
	if err != nil {
		return nil, err
	}

	return &Module{
		TaskService:         taskService,
		ContextResolver:     contextResolver,
		WebhookProcessor:    webhookProcessor,
		mode:                deps.TaskRunnerOptions.Mode,
		dispatcher:          dispatcher,
		workerPool:          workerPool,
		cronService:         cronService,
		taskRecoveryService: taskRecoveryService,
		semaphoreManager:    semaphoreManager,
		leaderRunner:        leaderRunner,
		useCases:            useCases,
	}, nil
}

func (m *Module) Health(ctx context.Context) error {
	if err := m.dispatcher.Ping(ctx); err != nil {
		return err
	}
	if m.mode != config.TaskRunnerModeQueue {
		return nil
	}
	if err := m.semaphoreManager.Ping(ctx); err != nil {
		return err
	}
	return m.leaderRunner.Ping(ctx)
}

func resolveWorkerID(configured string) string {
	if trimmed := strings.TrimSpace(configured); trimmed != "" {
		return trimmed
	}

	hostname, err := os.Hostname()
	if err != nil || strings.TrimSpace(hostname) == "" {
		slog.Warn("hostname not available, falling back to ephemeral worker id",
			"component", "taskrunner_module")
		return "task-runner-" + bnuuid.NewV7().String()[:8]
	}
	return hostname
}

func (m *Module) StartAsMain(ctx context.Context) error {
	m.workerPool.Start(ctx)

	if err := m.dispatcher.Start(ctx); err != nil {
		return err
	}

	go m.leaderRunner.Run(ctx, func(leaderCtx context.Context) {
		m.cronService.Start()
		defer func() {
			if err := m.cronService.Stop(); err != nil {
				slog.ErrorContext(leaderCtx, "failed to stop cron after losing leadership", "error", err)
			}
		}()

		go m.taskRecoveryService.Start(leaderCtx)

		<-leaderCtx.Done()
	})

	return nil
}

func (m *Module) StartAsWorker(ctx context.Context) error {
	if m.mode != config.TaskRunnerModeQueue {
		slog.ErrorContext(ctx, "task worker started but TASK_RUNNER_MODE is not 'queue'",
			"component", "taskrunner_module",
			"mode", string(m.mode))
		return ErrInvalidTaskRunnerMode
	}

	m.workerPool.Start(ctx)
	return m.dispatcher.Start(ctx)
}

func (m *Module) Shutdown() error {
	if err := m.dispatcher.Stop(); err != nil {
		slog.Error("failed to stop task dispatcher", "error", err)
	}
	m.workerPool.Stop()
	return nil
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware, apiKeyMiddleware *auth.APIKeyMiddleware[apikeysContract.Scope]) {
	taskRunnerHTTP.RegisterRoutes(router, authMiddleware, apiKeyMiddleware, m.useCases)
}
