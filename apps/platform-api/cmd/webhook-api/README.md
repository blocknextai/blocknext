# webhook-api

> The inbound webhook edge: a Fiber HTTP server that receives external workflow trigger webhooks and dispatches them to the taskrunner's webhook processor.

## What it does
Loads the webhook API config and bootstraps a near-full module graph via `bootstrap.NewWebhookAPI`, then serves a Fiber app that registers only the `WebhooksModule` routes. The webhooks module is wired with the taskrunner `WebhookProcessor`, which resolves and fires workflow trigger webhooks. It includes a realtime broadcaster, liveness/readiness probes, and a Redis-backed rate limiter.

## Bootstrap & config
- **Assembler:** `bootstrap.NewCore` + `bootstrap.NewWebhookAPI` — wires the workflow/trigger/taskrunner stack needed to process trigger webhooks.
- **Config:** `config.LoadWebhookAPI()` → `config.WebhookAPIConfig` (embeds `SharedConfig`, `HTTPServer`, `TaskRunner`).
- **Runs as:** HTTP server on `HTTP_SERVER_*` address.

## Bounded contexts activated
- common, account, organizations, credentialoauth
- nodeengine, platform, credentials, llm
- workflows, executions, triggers
- taskrunner, webhooks

## Notes
- Only `WebhooksModule.Register` mounts routes; no auth/cache/api-key middleware is applied to webhook endpoints.
- Builds a `realtime` broadcaster (for task-runner execution events); `Health` pings DB, cache, broadcaster, and the task runner.
- Shutdown closes the task runner module and broadcaster, then core.
