package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/blocknextai/file-gateway-api/internal/auth"
	cacheInfrastructure "github.com/blocknextai/file-gateway-api/internal/cache/infrastructure"
	"github.com/blocknextai/file-gateway-api/internal/config"
	"github.com/blocknextai/file-gateway-api/internal/download"
	"github.com/blocknextai/file-gateway-api/internal/storage"
	"github.com/blocknextai/file-gateway-api/internal/upload"
	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/auth/jwt"
	"github.com/blocknextai/go-packages/fiber/errorhandler"
	"github.com/blocknextai/go-packages/fiber/middleware/recovery"
	cachestorage "github.com/blocknextai/go-packages/fiber/storage/cache"
	"github.com/blocknextai/go-packages/json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/static"
)

const (
	appName = "file-gateway-api"
)

func buildTrustProxyConfig(trustProxies string) fiber.TrustProxyConfig {
	trimmed := strings.TrimSpace(trustProxies)
	if trimmed == "" {
		return fiber.TrustProxyConfig{
			Loopback: true,
			Private:  true,
		}
	}
	parts := strings.Split(trimmed, ",")
	proxies := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			proxies = append(proxies, v)
		}
	}
	return fiber.TrustProxyConfig{
		Proxies: proxies,
	}
}

func main() {
	configuration, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	isProduction := configuration.AppEnv.IsProduction()

	cacheService, err := cacheInfrastructure.NewCacheService(configuration.Cache)
	if err != nil {
		slog.Error("failed to create cache service", "error", err)
		os.Exit(1)
	}

	storageModule, err := storage.NewModule(storage.Dependencies{
		Options: configuration.Storage,
	})
	if err != nil {
		slog.Error("failed to get storage provider", "error", err)
		os.Exit(1)
	}

	app := fiber.New(fiber.Config{
		AppName:           appName,
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
		ReadTimeout:       configuration.API.ReadTimeout,
		WriteTimeout:      configuration.API.WriteTimeout,
		IdleTimeout:       configuration.API.IdleTimeout,
		ReduceMemoryUsage: configuration.API.ReduceMemoryUsage,
		BodyLimit:         configuration.API.BodyLimit,
		ReadBufferSize:    configuration.API.ReadBufferSize,
		WriteBufferSize:   configuration.API.WriteBufferSize,
		TrustProxy:        configuration.API.TrustProxyEnabled,
		TrustProxyConfig:  buildTrustProxyConfig(configuration.API.TrustProxies),
		ProxyHeader:       configuration.API.ProxyHeader,
		ErrorHandler:      errorhandler.New(isProduction),
	})

	app.Use(requestid.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Split(configuration.API.AllowOrigins, ","),
		AllowHeaders: strings.Split(configuration.API.AllowHeaders, ","),
		AllowMethods: strings.Split(configuration.API.AllowMethods, ","),
		MaxAge:       configuration.API.CORSMaxAge,
	}))

	app.Use(helmet.New())

	/*var metricsMiddleware *metrics.Middleware
	if configuration.API.Metrics.Enabled {
		metricsMiddleware, err = metrics.New(
			metrics.WithNamespace(strings.ReplaceAll(appName, "-", "_")),
			metrics.WithGoCollectors(),
		)
		if err != nil {
			slog.Error("failed to create metrics middleware", "error", err)
			os.Exit(1)
		}
		app.Use(metricsMiddleware.Collect())
	}*/

	app.Use(recovery.New())

	if configuration.Storage.Driver == config.StorageDriverLocal {
		app.Get("/uploads/*", static.New(configuration.Storage.Local.Public.Path))
	}

	app.Get(healthcheck.LivenessEndpoint, healthcheck.New())

	app.Get(healthcheck.ReadinessEndpoint, healthcheck.New(healthcheck.Config{
		Probe: func(c fiber.Ctx) bool {
			ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
			defer cancel()

			if err := storageModule.Provider.HealthCheckPublic(ctx); err != nil {
				slog.Warn("readiness probe failed", "component", "storage.public", "error", err)
				return false
			}

			if err := storageModule.Provider.HealthCheckPrivate(ctx); err != nil {
				slog.Warn("readiness probe failed", "component", "storage.private", "error", err)
				return false
			}

			if err := cacheService.Ping(ctx); err != nil {
				slog.WarnContext(ctx, "readiness probe: cache ping failed", "error", err)
				return false
			}

			return true
		},
	}))

	jwtService, err := jwt.New(
		configuration.Auth.JWT.Issuer,
		configuration.Auth.JWT.Audience,
		configuration.Auth.JWT.Secret,
		configuration.Auth.JWT.TokenExpirationTime,
		configuration.Auth.JWT.TokenExpirationTime,
		configuration.Auth.JWT.TokenLeeway,
	)
	if err != nil {
		slog.Error("failed to create jwt service", "error", err)
		os.Exit(1)
	}

	authModule := auth.NewModule(auth.Dependencies{
		ServiceKey: configuration.Auth.ServiceKey,
		JWTService: jwtService,
	})

	app.Use(limiter.New(limiter.Config{
		Next:       authModule.MatchServiceKey,
		Max:        configuration.API.MaxRequests,
		Expiration: configuration.API.ExpirationTime,
		Storage:    cachestorage.New(cacheService, appName+":rate_limit:"),
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c fiber.Ctx) error {
			return apperror.RateLimited("rate limit exceeded")
		},
	}))

	uploadModule := upload.NewModule(upload.Dependencies{
		StorageProvider: storageModule.Provider,
	})

	downloadModule := download.NewModule(download.Dependencies{
		Options: configuration.Download,
	})

	authModule.Register(app)

	protectedRouter := app.Group("", authModule.Middleware)
	uploadModule.Register(protectedRouter)
	downloadModule.Register(protectedRouter)

	/*var metricsApp *fiber.App
	if metricsMiddleware != nil {
		metricsApp = fiber.New(fiber.Config{
			AppName:     appName + "-metrics",
			JSONEncoder: json.Marshal,
			JSONDecoder: json.Unmarshal,
		})
		metricsApp.Get("/metrics", metricsMiddleware.Expose())

		go func() {
			address := configuration.API.MetricsAddress()
			slog.Info("metrics listening", "address", address, "app", metricsApp.Config().AppName)
			if err := metricsApp.Listen(address, fiber.ListenConfig{DisableStartupMessage: true}); err != nil {
				slog.Error("metrics listen failed", "error", err)
			}
		}()
	}*/

	go func() {
		address := configuration.API.Address()
		if err := app.Listen(address, fiber.ListenConfig{EnablePrefork: configuration.API.IsPrefork}); err != nil {
			slog.Error("failed to start server", "address", address, "error", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	slog.Info("server ready, press ctrl+c to stop")
	<-sigChan
	slog.Info("shutting down server")

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		slog.Error("failed to shutdown app server", "error", err)
	}

	/*if metricsApp != nil {
		if err := metricsApp.ShutdownWithContext(shutdownCtx); err != nil {
			slog.Error("failed to shutdown metrics server", "error", err)
		}
	}*/

	if err := cacheService.Close(); err != nil {
		slog.Error("failed to close cache service", "error", err)
	}

	slog.Info("shutdown complete")
}
