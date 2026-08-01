package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	cachestorage "github.com/blocknextai/go-packages/fiber/storage/cache"
	"github.com/blocknextai/platform-api/internal/bootstrap"
	commonPresentationAuth "github.com/blocknextai/platform-api/internal/common/presentation/auth"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/config"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

const (
	appName = "platform-api"
)

func main() {
	cfg, err := config.LoadPlatformAPI()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	core, err := bootstrap.NewCore(cfg.Shared)
	if err != nil {
		slog.Error("failed to bootstrap core", "error", err)
		os.Exit(1)
	}

	app, err := bootstrap.NewPlatformAPI(core, cfg)
	if err != nil {
		slog.Error("failed to bootstrap platform-api", "error", err)
		os.Exit(1)
	}

	fiberApp /*metricsMiddleware,*/, err := bootstrap.NewFiber(appName, cfg.HTTPServer, cfg.Shared.AppEnv.IsProduction())
	if err != nil {
		slog.Error("failed to bootstrap fiber", "error", err)
		os.Exit(1)
	}

	//metricsShutdown := bootstrap.StartMetrics(metricsMiddleware, appName, cfg.HTTPServer)

	fiberApp.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	fiberApp.Get(healthcheck.ReadinessEndpoint, healthcheck.New(healthcheck.Config{
		Probe: func(c fiber.Ctx) bool {
			ctx, cancel := context.WithTimeout(c.RequestCtx(), 2*time.Second)
			defer cancel()

			if err := app.Health(ctx); err != nil {
				slog.WarnContext(ctx, "readiness probe failed", "error", err)
				return false
			}
			return true
		},
	}))

	if cfg.HTTPServer.MaxRequests > 0 {
		fiberApp.Use(limiter.New(limiter.Config{
			Max:        cfg.HTTPServer.MaxRequests,
			Expiration: cfg.HTTPServer.ExpirationTime,
			Storage:    cachestorage.New(core.CacheService, appName+":rate_limit:"),
			KeyGenerator: func(c fiber.Ctx) string {
				return c.IP()
			},
			LimitReached: func(c fiber.Ctx) error {
				return commonHTTP.ErrRateLimitReached
			},
		}))
	}

	authMiddleware := commonPresentationAuth.NewAuthMiddleware(
		app.JWTService,
		app.AccountModule.UserPermissionChecker,
		app.OrganizationsModule.OrganizationPermissionChecker,
		app.AccountModule.SessionService,
	)
	apiKeyMiddleware := commonPresentationAuth.NewAPIKeyMiddleware(
		app.APIKeysModule.APIKeyValidator,
	)
	cacheMiddleware := cachemiddleware.New(core.CacheService, appName+":cache:")

	app.AccountModule.Register(fiberApp, authMiddleware, cacheMiddleware)
	app.OrganizationsModule.Register(fiberApp, authMiddleware, cacheMiddleware)
	app.ExecutionsModule.Register(fiberApp, authMiddleware)
	app.TriggersModule.Register(fiberApp, authMiddleware)
	app.CredentialsModule.Register(fiberApp, authMiddleware)
	app.APIKeysModule.Register(fiberApp, authMiddleware, cacheMiddleware)
	app.NotificationsModule.Register(fiberApp, authMiddleware)
	app.WorkflowsModule.Register(fiberApp, authMiddleware)
	app.CredentialOAuthModule.Register(fiberApp, authMiddleware)
	app.NodeEngineModule.Register(fiberApp, cacheMiddleware)
	app.TaskRunnerModule.Register(fiberApp, authMiddleware, apiKeyMiddleware)
	app.PlatformModule.Register(fiberApp, authMiddleware, cacheMiddleware)
	app.WSModule.Register(fiberApp, authMiddleware)

	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := app.TaskRunnerModule.StartAsMain(appCtx); err != nil {
		slog.Error("failed to start task runner module", "error", err)
		os.Exit(1)
	}

	if cfg.TaskRunner.Mode == config.TaskRunnerModeEmbedded {
		core.EventBus.StartRelay(appCtx)
	}

	bootstrap.ListenAndWait(fiberApp, cfg.HTTPServer,
		func() error { cancel(); return nil },
		//metricsShutdown,
		app.Shutdown,
		core.Shutdown,
	)
}
