package bootstrap

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/fiber/errorhandler"
	loggingmiddleware "github.com/blocknextai/go-packages/fiber/middleware/logging"
	"github.com/blocknextai/go-packages/fiber/middleware/recovery"
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/config"
)

const (
	fiberShutdownTimeout = 10 * time.Second
)

func NewFiber(appName string, opts config.HTTPServerOptions, isProduction bool) (*fiber.App /* *metrics.Middleware,*/, error) {
	app := fiber.New(fiber.Config{
		AppName:           appName,
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
		ReadTimeout:       opts.ReadTimeout,
		WriteTimeout:      opts.WriteTimeout,
		IdleTimeout:       opts.IdleTimeout,
		ReduceMemoryUsage: opts.ReduceMemoryUsage,
		BodyLimit:         opts.BodyLimit,
		ReadBufferSize:    opts.ReadBufferSize,
		WriteBufferSize:   opts.WriteBufferSize,
		TrustProxy:        opts.TrustProxyEnabled,
		TrustProxyConfig:  buildTrustProxyConfig(opts.TrustProxies),
		ProxyHeader:       opts.ProxyHeader,
		ErrorHandler:      newErrorHandler(isProduction),
	})

	/*var metricsMiddleware *metrics.Middleware
	if opts.Metrics.Enabled {
		m, err := metrics.New(
			metrics.WithNamespace(strings.ReplaceAll(appName, "-", "_")),
			metrics.WithGoCollectors(),
		)
		if err != nil {
			return nil, nil, err
		}
		metricsMiddleware = m
	}*/

	app.Use(requestid.New())
	app.Use(loggingmiddleware.New())

	if strings.TrimSpace(opts.AllowOrigins) != "" {
		app.Use(cors.New(cors.Config{
			AllowOrigins:  strings.Split(opts.AllowOrigins, ","),
			AllowHeaders:  strings.Split(opts.AllowHeaders, ","),
			AllowMethods:  strings.Split(opts.AllowMethods, ","),
			ExposeHeaders: strings.Split(opts.ExposeHeaders, ","),
			MaxAge:        opts.CORSMaxAge,
		}))
	}

	app.Use(helmet.New())
	/*if metricsMiddleware != nil {
		app.Use(metricsMiddleware.Collect())
	}*/
	app.Use(recovery.New())

	return app /*metricsMiddleware,*/, nil
}

func newErrorHandler(isProduction bool) fiber.ErrorHandler {
	handler := errorhandler.New(isProduction)

	return func(c fiber.Ctx, err error) error {
		var fiberErr *fiber.Error
		if _, isAppError := apperror.As(err); isAppError || !errors.As(err, &fiberErr) {
			return handler(c, err)
		}

		return c.Status(fiberErr.Code).JSON(result.Fail[any](fiberErr))
	}
}

func ListenAndWait(app *fiber.App, opts config.HTTPServerOptions, shutdownFns ...func() error) {
	go func() {
		address := opts.Address()
		slog.Info("listening", "address", address, "app", app.Config().AppName)
		if err := app.Listen(address, fiber.ListenConfig{EnablePrefork: opts.IsPrefork}); err != nil {
			slog.Error("listen failed", "error", err)
			os.Exit(1)
		}
	}()

	fiberShutdown := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), fiberShutdownTimeout)
		defer cancel()
		return app.ShutdownWithContext(ctx)
	}

	allFns := append([]func() error{fiberShutdown}, shutdownFns...)
	WaitForShutdown(fiberShutdownTimeout+5*time.Second, allFns...)
}

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
	return fiber.TrustProxyConfig{Proxies: proxies}
}
