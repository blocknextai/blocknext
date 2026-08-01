package bootstrap

/*func StartMetrics(middleware *metrics.Middleware, appName string, opts config.HTTPServerOptions) func() error {
	if middleware == nil {
		return nil
	}

	metricsApp := fiber.New(fiber.Config{
		AppName:     appName + "-metrics",
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	metricsApp.Get("/metrics", middleware.Expose())

	go func() {
		address := opts.MetricsAddress()
		slog.Info("metrics listening", "address", address, "app", metricsApp.Config().AppName)
		if err := metricsApp.Listen(address, fiber.ListenConfig{DisableStartupMessage: true}); err != nil {
			slog.Error("metrics listen failed", "error", err)
		}
	}()

	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), fiberShutdownTimeout)
		defer cancel()
		return metricsApp.ShutdownWithContext(ctx)
	}
}*/
