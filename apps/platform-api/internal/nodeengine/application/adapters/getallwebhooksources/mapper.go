package getallwebhooksources

import (
	"strings"

	"github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
)

func MapGetAllWebhookSourcesQueryToGetAllWebhookSourcesResponse(
	adapters []adapters.TriggerAdapter,
	webhookURLTemplate string,
) *GetAllWebhookSourcesResponse {
	response := make(GetAllWebhookSourcesResponse, 0, len(adapters))
	for _, adapter := range adapters {
		url := strings.ReplaceAll(webhookURLTemplate, "{source}", adapter.GetID())
		response = append(response, WebhookSourceResponse{
			Source:     adapter.GetID(),
			WebhookURL: url,
		})
	}
	return &response
}
