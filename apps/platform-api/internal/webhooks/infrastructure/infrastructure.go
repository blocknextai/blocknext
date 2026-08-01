package infrastructure

import (
	taskRunnerApplicationWebhooks "github.com/blocknextai/platform-api/internal/taskrunner/application/webhooks"
)

type Handlers struct {
	TriggerProcessor taskRunnerApplicationWebhooks.WebhookProcessor
}

type RegisterInfrastructureDeps struct {
	TriggerWebhookProcessor taskRunnerApplicationWebhooks.WebhookProcessor
}

func RegisterInfrastructure(deps RegisterInfrastructureDeps) *Handlers {
	return &Handlers{
		TriggerProcessor: deps.TriggerWebhookProcessor,
	}
}
