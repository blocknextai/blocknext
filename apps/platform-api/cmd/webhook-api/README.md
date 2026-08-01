# webhook-api

> The inbound webhook edge: a Fiber HTTP server that receives external webhooks (payment provider callbacks and workflow trigger webhooks) and dispatches them to the right processor.

## What it does
Loads the webhook API config and bootstraps a near-full module graph via `bootstrap.NewWebhookAPI`, then serves a Fiber app that registers only the `WebhooksModule` routes. The webhooks module is wired with two processors: the payments `WebhookProcessor` (provider callbacks) and the taskrunner `WebhookProcessor` (resolves and fires workflow trigger webhooks). It includes a realtime broadcaster, liveness/readiness probes, a Redis-backed rate limiter, and metrics.

## Bootstrap & config
- **Assembler:** `bootstrap.NewCore` + `bootstrap.NewWebhookAPI` — wires the workflow/trigger/taskrunner stack plus payments so both webhook kinds can be processed.
- **Config:** `config.LoadWebhookAPI()` → `config.WebhookAPIConfig` (embeds `SharedConfig`, `HTTPServer`, `TaskRunner`).
- **Runs as:** HTTP server on `HTTP_SERVER_*` address; metrics on a separate metrics server.

## Bounded contexts activated
- common, account, organizations, web3, oauth, plans, payments, subscriptions, quota
- nodeengine, platform, credentials, llm
- workflows, marketplace, library, executions, triggers
- taskrunner, webhooks

## Notes
- Only `WebhooksModule.Register` mounts routes; no auth/cache/api-key middleware is applied to webhook endpoints.
- Builds a `realtime` broadcaster (for task-runner execution events); `Health` pings DB, cache, broadcaster, and the task runner.
- Shutdown closes the task runner module and broadcaster, then core.
