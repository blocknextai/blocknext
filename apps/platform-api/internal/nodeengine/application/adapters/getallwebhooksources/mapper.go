package getallwebhooksources

import (
	"strings"

	adaptersDomain "github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
)

func MapGetAllWebhookSourcesQueryToGetAllWebhookSourcesResponse(
	adapters []adaptersDomain.TriggerAdapter,
	webhookURLTemplate string,
) *GetAllWebhookSourcesResponse {
	response := make(GetAllWebhookSourcesResponse, 0, len(adapters))
	for _, adapter := range adapters {
		url := strings.ReplaceAll(webhookURLTemplate, "{source}", adapter.GetID())
		_, supportsVerification := adaptersDomain.AsWebhookVerifier(adapter)
		response = append(response, WebhookSourceResponse{
			Source:               adapter.GetID(),
			Name:                 adapter.GetName(),
			SupportsVerification: supportsVerification,
			WebhookURL:           url,
		})
	}
	return &response
}
