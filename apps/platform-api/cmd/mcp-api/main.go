package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/limiter"

	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	cachestorage "github.com/blocknextai/go-packages/fiber/storage/cache"
	"github.com/blocknextai/platform-api/internal/bootstrap"
	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	"github.com/blocknextai/platform-api/internal/config"
)

const (
	appName = "mcp-api"
)

func main() {
	cfg, err := config.LoadMCPAPI()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	core, err := bootstrap.NewCore(cfg.Shared)
	if err != nil {
		slog.Error("failed to bootstrap core", "error", err)
		os.Exit(1)
	}

	app, err := bootstrap.NewMCPAPI(core, cfg)
	if err != nil {
		slog.Error("failed to bootstrap mcp-api", "error", err)
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

	authMiddleware := commonAuth.NewAuthMiddleware(
		app.JWTService,
		app.AccountModule.UserPermissionChecker,
		app.OrganizationsModule.OrganizationPermissionChecker,
		app.AccountModule.SessionService,
	)
	apiKeyMiddleware := commonAuth.NewAPIKeyMiddleware(
		app.APIKeysModule.APIKeyValidator,
	)
	cacheMiddleware := cachemiddleware.New(core.CacheService, appName+":cache:")

	accessTokenMiddleware := commonAuth.NewAccessTokenMiddleware(
		app.MCPOAuthModule.AccessTokenValidator,
		app.MCPOAuthModule.MetadataService.GetProtectedResourceMetadataURL(""),
	)

	app.MCPOAuthModule.Register(fiberApp, authMiddleware, cacheMiddleware)
	app.MCPModule.Register(fiberApp, cacheMiddleware, apiKeyMiddleware, accessTokenMiddleware)

	bootstrap.ListenAndWait(fiberApp, cfg.HTTPServer,
		/*metricsShutdown,*/
		core.Shutdown,
	)
}
